package fileCache

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/fsnotify/fsnotify"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/util/flowcontrol"
	"k8s.io/klog/v2"
)

type fileEventType int

const (
	FileAdd fileEventType = iota
	FileModify
	FileDelete
	FileListAdd
)

const (
	retryPeriod        = 1 * time.Second
	maxRetryPeriod     = 20 * time.Second
	FileCheckFrequency = 15 * time.Second
)

const (
	maxConfigLength = 10 * 1 << 20 // 10MB
)

const FileSuffix = ".yaml"

type watchFileEvent struct {
	fileName  string
	eventType fileEventType
}

type ObjectSourceFile struct {
	path        string
	regularList bool
	watchEvents chan<- *watchFileEvent
}

func (w *watchFileEvent) String() string {
	ope := "UnKnown"
	switch w.eventType {
	case FileAdd:
		ope = "Add"
	case FileModify:
		ope = "Modify"
	case FileDelete:
		ope = "Delete"
	case FileListAdd:
		ope = "ListAdd"
	}
	return fmt.Sprintf("%s File %s", ope, w.fileName)
}

// NewObjectSourceFile watches a object file for changes.
func NewObjectSourceFile(path string, regularList bool, events chan<- *watchFileEvent) {
	// "github.com/sigma/go-inotify" requires a path without trailing "/"
	path = strings.TrimRight(path, string(os.PathSeparator))

	config := newObjectSourceFile(path, regularList, events)
	config.run()
}

func newObjectSourceFile(path string, regularList bool, events chan<- *watchFileEvent) *ObjectSourceFile {
	return &ObjectSourceFile{
		path:        path,
		regularList: regularList,
		watchEvents: events,
	}
}

func (s *ObjectSourceFile) run() {

	listTicker := time.NewTicker(FileCheckFrequency)

	go func() {
		// Read path immediately to speed up startup.
		if err := s.listConfig(); err != nil {
			klog.Errorf("Unable to read config path %q: %v", s.path, err)
		}
		if !s.regularList {
			klog.Warningf("ObjectSourceFile set regularList %v, so do not regular list dir %v", s.regularList, s.path)
			return
		}
		for {
			select {
			case <-listTicker.C:
				if err := s.listConfig(); err != nil {
					klog.Errorf("Unable to read config path %q: %v", s.path, err)
				}
			}
		}
	}()

	s.startWatch()
}

func (s *ObjectSourceFile) startWatch() {
	backOff := flowcontrol.NewBackOff(retryPeriod, maxRetryPeriod)
	backOffID := "watch"

	go wait.Forever(func() {
		if backOff.IsInBackOffSinceUpdate(backOffID, time.Now()) {
			return
		}

		if err := s.doWatch(); err != nil {
			klog.Errorf("Unable to read config path %q: %v", s.path, err)
		}
	}, retryPeriod)
}

func (s *ObjectSourceFile) listConfig() error {
	path := s.path
	statInfo, err := os.Stat(path)
	if err != nil {
		if !os.IsNotExist(err) {
			return err
		}
		// Emit an update with an empty PodList to allow FileSource to be marked as seen
		return fmt.Errorf("path does not exist, ignoring")
	}

	switch {
	case statInfo.Mode().IsDir():
		return s.extractFromDir(path)

	case statInfo.Mode().IsRegular():
		klog.V(4).Infof("listConfig: find path %s is regular, create watchfile event", path)
		s.watchEvents <- &watchFileEvent{path, FileListAdd}
		return nil
	default:
		return fmt.Errorf("path is not a directory or file")
	}
}

// Get as many object manifests as we can from a directory. Return an error if and only if something
// prevented us from reading anything at all. Do not return an error if only some files
// were problematic.
func (s *ObjectSourceFile) extractFromDir(name string) error {
	dirents, err := filepath.Glob(filepath.Join(name, fmt.Sprintf("[^.]*%s", FileSuffix)))
	if err != nil {
		return fmt.Errorf("glob failed: %v", err)
	}

	if len(dirents) == 0 {
		return nil
	}

	sort.Strings(dirents)
	for _, path := range dirents {
		statInfo, err := os.Stat(path)
		if err != nil {
			klog.Errorf("Can't get metadata for %q: %v", path, err)
			continue
		}

		switch {
		case statInfo.Mode().IsDir():
			klog.Errorf("Not recursing into manifest path %q", path)
		case statInfo.Mode().IsRegular():
			klog.V(4).Infof("extractFromDir: find path %s is regular, create watchfile event", path)
			s.watchEvents <- &watchFileEvent{path, FileListAdd}
		default:
			klog.Errorf("Manifest path %q is not a directory or file: %v", path, statInfo.Mode())
		}
	}
	return nil
}

func (s *ObjectSourceFile) doWatch() error {
	w, err := fsnotify.NewWatcher()
	if err != nil {
		return fmt.Errorf("unable to create inotify: %v", err)
	}
	defer w.Close()
	path := s.path

	err = w.Add(path)
	if err != nil {
		return fmt.Errorf("unable to create inotify for path %q: %v", path, err)
	}

	for {
		select {
		case event := <-w.Events:
			if err = s.produceWatchEvent(&event); err != nil {
				return fmt.Errorf("error while processing inotify event (%+v): %v", event, err)
			}
		case err = <-w.Errors:
			return fmt.Errorf("error while watching %q: %v", path, err)
		}
	}
}

func (s *ObjectSourceFile) produceWatchEvent(e *fsnotify.Event) error {
	// Ignore file start with dots
	if strings.HasPrefix(filepath.Base(e.Name), ".") {
		klog.V(4).Infof("Ignored object manifest: %s, because it starts with dots", e.Name)
		return nil
	}

	if !strings.HasSuffix(filepath.Base(e.Name), FileSuffix) {
		klog.V(4).Infof("Ignored object manifest: %s, because it not end with %s", e.Name, FileSuffix)
		return nil
	}
	var eventType fileEventType
	switch {
	case (e.Op & fsnotify.Create) > 0:
		eventType = FileAdd
	case (e.Op & fsnotify.Write) > 0:
		eventType = FileModify
		/*
			case (e.Op & fsnotify.Chmod) > 0:
				fmt.Println("Chmod ...")
				eventType = FileModify
		*/
	case (e.Op & fsnotify.Remove) > 0:
		eventType = FileDelete
	case (e.Op & fsnotify.Rename) > 0:
		eventType = FileDelete
	default:
		// Ignore rest events
		return nil
	}

	klog.V(4).Infof("produceWatchEvent: name %s, type %v", e.Name, eventType)
	s.watchEvents <- &watchFileEvent{e.Name, eventType}
	return nil
}
