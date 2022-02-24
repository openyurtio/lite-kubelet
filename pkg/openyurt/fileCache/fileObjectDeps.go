package fileCache

import (
	"fmt"
	"os"

	"path/filepath"
	"sort"

	coordinationv1 "k8s.io/api/coordination/v1"
	v1 "k8s.io/api/core/v1"
	"k8s.io/client-go/tools/cache"
	"k8s.io/klog/v2"
	"k8s.io/kubernetes/pkg/openyurt/manifest"
)

type FileObjectDeps interface {
	GetDir() string
	GetDefaultObject() interface{}
	GetObjectKeyFunc() cache.KeyFunc
	GetFileKeyFunc() FileIndexFunc
	GetObjectString() string
	GetAllFiles() []string
	GetFullFileName(obj interface{}) (string, error)
}

type FileConfigMapDeps struct {
	PathDir     string
	FileKeyFunc FileIndexFunc
}

func (c *FileConfigMapDeps) GetFullFileName(obj interface{}) (string, error) {
	name, err := CreateFileNameByObject(obj)
	if err != nil {
		return "", err
	}
	return filepath.Join(c.GetDir(), name), nil
}

func (c *FileConfigMapDeps) GetAllFiles() []string {
	return getAllFiles(c.GetDir())
}

func NewDefaultFileConfigMapDeps() *FileConfigMapDeps {
	return &FileConfigMapDeps{
		PathDir: manifest.GetConfigmapsManifestPath(),
	}
}

func (c *FileConfigMapDeps) GetObjectString() string {
	return "ConfigMap"
}
func (c *FileConfigMapDeps) GetDir() string {
	return c.PathDir
}

func (c *FileConfigMapDeps) GetDefaultObject() interface{} {
	return &v1.ConfigMap{}
}

func (c *FileConfigMapDeps) GetObjectKeyFunc() cache.KeyFunc {
	//	cache.MetaNamespaceKeyFunc, FileNamespaceNameKeyFunc
	return cache.MetaNamespaceKeyFunc
}

func (c *FileConfigMapDeps) GetFileKeyFunc() FileIndexFunc {
	return FileNamespaceNameKeyFunc
}

type FileServiceDeps struct {
	PathDir     string
	FileKeyFunc FileIndexFunc
}

func (f *FileServiceDeps) GetFullFileName(obj interface{}) (string, error) {
	name, err := CreateFileNameByObject(obj)
	if err != nil {
		return "", err
	}
	return filepath.Join(f.GetDir(), name), nil
}

func (f *FileServiceDeps) GetAllFiles() []string {
	return getAllFiles(f.GetDir())
}

func NewDefaultFileServiceDeps() *FileServiceDeps {
	return &FileServiceDeps{
		PathDir: manifest.GetServicesManifestPath(),
	}
}
func (f *FileServiceDeps) GetObjectString() string {
	return "Service"
}
func (f *FileServiceDeps) GetDir() string {
	return f.PathDir
}

func (f *FileServiceDeps) GetDefaultObject() interface{} {
	return &v1.Service{}
}

func (f *FileServiceDeps) GetObjectKeyFunc() cache.KeyFunc {
	return cache.MetaNamespaceKeyFunc
}

func (f *FileServiceDeps) GetFileKeyFunc() FileIndexFunc {
	return FileNamespaceNameKeyFunc
}

type FileNodeDeps struct {
	PathDir     string
	FileKeyFunc FileIndexFunc
}

func (f *FileNodeDeps) GetFullFileName(obj interface{}) (string, error) {
	name, err := CreateFileNameByObject(obj)
	if err != nil {
		return "", err
	}
	return filepath.Join(f.GetDir(), name), nil
}

func (f *FileNodeDeps) GetAllFiles() []string {
	return getAllFiles(f.GetDir())
}

func NewDefaultFileNodeDeps() *FileNodeDeps {
	return &FileNodeDeps{
		PathDir: manifest.GetNodesManifestPath(),
	}
}
func (f *FileNodeDeps) GetObjectString() string {
	return "Node"
}

func (f *FileNodeDeps) GetDir() string {
	return f.PathDir
}

func (f *FileNodeDeps) GetDefaultObject() interface{} {
	return &v1.Node{}
}

func (f *FileNodeDeps) GetObjectKeyFunc() cache.KeyFunc {
	return cache.MetaNamespaceKeyFunc
}

func (f *FileNodeDeps) GetFileKeyFunc() FileIndexFunc {
	return FileNoNamespacedKeyFunc
}

type FileSecretDeps struct {
	PathDir     string
	FileKeyFunc FileIndexFunc
}

func (f *FileSecretDeps) GetFullFileName(obj interface{}) (string, error) {
	name, err := CreateFileNameByObject(obj)
	if err != nil {
		return "", err
	}
	return filepath.Join(f.GetDir(), name), nil
}

func (f *FileSecretDeps) GetAllFiles() []string {
	return getAllFiles(f.GetDir())
}

func (f *FileSecretDeps) GetDir() string {
	return f.PathDir
}

func (f *FileSecretDeps) GetDefaultObject() interface{} {
	return &v1.Secret{}
}

func (f *FileSecretDeps) GetObjectKeyFunc() cache.KeyFunc {
	return cache.MetaNamespaceKeyFunc
}

func (f *FileSecretDeps) GetFileKeyFunc() FileIndexFunc {
	return FileNamespaceNameKeyFunc
}

func (f *FileSecretDeps) GetObjectString() string {
	return "Secret"
}

func NewDefaultFileSecretDeps() *FileSecretDeps {
	return &FileSecretDeps{
		PathDir: manifest.GetSecretManifestPath(),
	}
}

type FilePodDeps struct {
	PathDir     string
	FileKeyFunc FileIndexFunc
}

func (f *FilePodDeps) GetFullFileName(obj interface{}) (string, error) {
	name, err := CreateFileNameByObject(obj)
	if err != nil {
		return "", err
	}
	if len(name) == 0 {
		return "", fmt.Errorf("can not get filename by object")
	}
	return filepath.Join(f.GetDir(), name), nil
}

func (f *FilePodDeps) GetAllFiles() []string {
	return getAllFiles(f.GetDir())
}

func (f *FilePodDeps) GetDir() string {
	return f.PathDir
}

func (f *FilePodDeps) GetDefaultObject() interface{} {
	return &v1.Pod{}
}

func (f *FilePodDeps) GetObjectKeyFunc() cache.KeyFunc {
	return cache.MetaNamespaceKeyFunc
}

func (f *FilePodDeps) GetFileKeyFunc() FileIndexFunc {
	return FileNamespaceNameKeyFunc
}

func (f *FilePodDeps) GetObjectString() string {
	return "Pod"
}

func NewDefaultFilePodDeps() *FilePodDeps {
	return &FilePodDeps{
		PathDir: manifest.GetPodsManifestPath(),
	}
}

type FileLeaseDeps struct {
	PathDir     string
	FileKeyFunc FileIndexFunc
}

func (f *FileLeaseDeps) GetFullFileName(obj interface{}) (string, error) {
	name, err := CreateFileNameByObject(obj)
	if err != nil {
		return "", err
	}
	return filepath.Join(f.GetDir(), name), nil
}

func (f *FileLeaseDeps) GetAllFiles() []string {
	return getAllFiles(f.GetDir())
}

func (f *FileLeaseDeps) GetDir() string {
	return f.PathDir
}

func (f *FileLeaseDeps) GetDefaultObject() interface{} {
	return &coordinationv1.Lease{}
}

func (f *FileLeaseDeps) GetObjectKeyFunc() cache.KeyFunc {
	return cache.MetaNamespaceKeyFunc
}

func (f *FileLeaseDeps) GetFileKeyFunc() FileIndexFunc {
	return FileNamespaceNameKeyFunc
}

func (f *FileLeaseDeps) GetObjectString() string {
	return "Lease"
}

func NewDefaultFileLeaseDeps() *FileLeaseDeps {
	return &FileLeaseDeps{
		PathDir: manifest.GetLeasesManifestPath(),
	}
}

type FileEventDeps struct {
	PathDir     string
	FileKeyFunc FileIndexFunc
}

func (f *FileEventDeps) GetFullFileName(obj interface{}) (string, error) {
	name, err := CreateFileNameByObject(obj)
	if err != nil {
		return "", err
	}
	return filepath.Join(f.GetDir(), name), nil
}

func (f *FileEventDeps) GetAllFiles() []string {
	return getAllFiles(f.GetDir())
}

func (f *FileEventDeps) GetDir() string {
	return f.PathDir
}

func (f *FileEventDeps) GetDefaultObject() interface{} {
	return &v1.Event{}
}

func (f *FileEventDeps) GetObjectKeyFunc() cache.KeyFunc {
	return cache.MetaNamespaceKeyFunc
}

func (f *FileEventDeps) GetFileKeyFunc() FileIndexFunc {
	return FileNamespaceNameKeyFunc
}

func (f *FileEventDeps) GetObjectString() string {
	return "Event"
}

func NewDefaultFileEventDeps() *FileEventDeps {
	return &FileEventDeps{
		PathDir: manifest.GetEventsManifestPath(),
	}
}

// get all full path files
func getAllFiles(dir string) []string {
	names := make([]string, 0, 10)
	dirents, err := filepath.Glob(filepath.Join(dir, fmt.Sprintf("[^.]*%s", FileSuffix)))
	if err != nil {
		klog.Errorf("glob dir %s failed: %v", dir, err)
		return names
	}

	if len(dirents) == 0 {
		return nil
	}

	sort.Strings(dirents)
	for _, path := range dirents {
		statInfo, err := os.Stat(path)
		if err != nil {
			klog.Errorf("Can't get stat for %q: %v", path, err)
			continue
		}

		switch {
		case statInfo.Mode().IsDir():
			klog.Errorf("%s is dir, do nothing", path)
		case statInfo.Mode().IsRegular():
			//names = append(names, filepath.Base(path))
			names = append(names, path)
		default:
			klog.Errorf("Path %q is not a directory or file: %v, do nothing", path, statInfo.Mode())
		}
	}
	return names
}

var _ FileObjectDeps = &FileConfigMapDeps{}
var _ FileObjectDeps = &FileServiceDeps{}
var _ FileObjectDeps = &FileNodeDeps{}
var _ FileObjectDeps = &FileSecretDeps{}
var _ FileObjectDeps = &FilePodDeps{}
var _ FileObjectDeps = &FileLeaseDeps{}
var _ FileObjectDeps = &FileEventDeps{}
