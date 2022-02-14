package fileCache

import (
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strings"

	"k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/client-go/tools/cache"
	"k8s.io/klog/v2"
	utilio "k8s.io/utils/io"
	"sigs.k8s.io/yaml"
)

const NameSpaceSpliceFileKey = "#"

var FileNameSpaceNameSpecificError = fmt.Errorf("fileName does not meet the specifications of [NameSpace#Name.yaml]")
var FileNoNameSpacedSpecificError = fmt.Errorf("fileName does not meet the specifications of [Name.yaml]")

// FileIndexFunc knows how to compute the set of indexed values for an file.
type FileIndexFunc func(path string) (string, error)

func CreateFileNameByObject(obj interface{}) (string, error) {
	meta, err := meta.Accessor(obj)
	if err != nil {
		return "", fmt.Errorf("object has no meta: %v", err)
	}
	if len(meta.GetNamespace()) > 0 {
		return meta.GetNamespace() + NameSpaceSpliceFileKey + meta.GetName() + FileSuffix, nil
	}
	return meta.GetName() + FileSuffix, nil
}

func FileNamespaceNameKeyFunc(path string) (string, error) {
	base := filepath.Base(path)
	if !strings.HasSuffix(base, FileSuffix) {
		return "", fmt.Errorf("%s need has suffix %s", path, FileSuffix)
	}
	nameSpaceName := strings.TrimSuffix(base, FileSuffix)
	if strings.HasSuffix(nameSpaceName, NameSpaceSpliceFileKey) || strings.HasPrefix(nameSpaceName, NameSpaceSpliceFileKey) {
		return "", FileNameSpaceNameSpecificError
	}

	split := strings.Split(nameSpaceName, NameSpaceSpliceFileKey)
	switch len(split) {
	case 1:
		return "", FileNameSpaceNameSpecificError
	case 2:
		return split[0] + "/" + split[1], nil
	default:
		return "", FileNameSpaceNameSpecificError
	}
}

func FileNoNamespacedKeyFunc(path string) (string, error) {
	base := filepath.Base(path)
	if !strings.HasSuffix(base, FileSuffix) {
		return "", fmt.Errorf("%s need has suffix %s", path, FileSuffix)
	}
	name := strings.TrimSuffix(base, FileSuffix)
	if strings.Contains(name, NameSpaceSpliceFileKey) {
		return "", FileNoNameSpacedSpecificError
	}
	return name, nil
}

func getNewDefault(v interface{}) interface{} {
	t := reflect.TypeOf(v)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	return reflect.New(t).Interface()
}

func NewFileObiectIndexer(ofi FileObjectDeps, regularList bool, send func(indexer cache.Indexer)) cache.Indexer {

	path := ofi.GetDir()

	watchEvents := make(chan *watchFileEvent, 10)
	NewObjectSourceFile(path, regularList, watchEvents)
	keyFunc := ofi.GetObjectKeyFunc()

	indexer := cache.NewIndexer(keyFunc, cache.Indexers{cache.NamespaceIndex: cache.MetaNamespaceIndexFunc})

	go func() {
		for {
			select {
			case e := <-watchEvents:
				klog.V(5).Info(e)
				switch e.eventType {
				case FileAdd, FileListAdd:
					filekey, err := ofi.GetFileKeyFunc()(e.fileName)
					if err != nil {
						klog.Errorf("Get file %s NameSpaceName Key error %v", e.fileName, err)
						break
					}
					if e.eventType == FileListAdd {
						_, exist, _ := indexer.GetByKey(filekey)
						if exist {
							klog.V(5).Infof("File %s is listAdd found, and has exist in cache, do nothing", e.fileName)
							break
						}
					}
					c := getNewDefault(ofi.GetDefaultObject())
					err = decodeFileObject(e.fileName, c)
					if err != nil {
						klog.Errorf("Decode file %s object[%s] error %v", e.fileName, ofi.GetObjectString(), err)
						break
					}

					objKey, err := keyFunc(c)
					if err != nil {
						klog.Errorf("Get object[%s][%s] key error %v", ofi.GetObjectString(), filekey, err)
						break
					}
					if objKey == filekey {
						indexer.Add(c)
						if send != nil {
							send(indexer)
						}
					} else {
						klog.Warningf("Object[%s] key %s dose not equal filekey %s, so do not add cache", ofi.GetObjectString(), objKey, filekey)
					}
				case FileModify:
					c := getNewDefault(ofi.GetDefaultObject())
					err := decodeFileObject(e.fileName, c)
					if err != nil {
						klog.Errorf("Decode file %s to object[%s] error %v", e.fileName, ofi.GetObjectString(), err)
						break
					}
					indexer.Update(c)
					if send != nil {
						send(indexer)
					}
				case FileDelete:
					key, err := ofi.GetFileKeyFunc()(e.fileName)
					if err != nil {
						klog.Errorf("Get file %s filekey error %v", e.fileName, err)
						break
					}
					old, exist, _ := indexer.GetByKey(key)
					if exist {
						indexer.Delete(old)
						if send != nil {
							send(indexer)
						}
					}
				}
			}
		}
	}()

	return indexer
}

func decodeFileObject(filename string, obj interface{}) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	data, err := utilio.ReadAtMost(file, maxConfigLength)
	if err != nil {
		return err
	}
	if err := yaml.Unmarshal(data, obj); err != nil {
		return err
	}
	return nil
}
