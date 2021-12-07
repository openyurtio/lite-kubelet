package fileCache

import (
	coordinationv1 "k8s.io/api/coordination/v1"
	v1 "k8s.io/api/core/v1"
	"k8s.io/client-go/tools/cache"
)

type FileObjectDeps interface {
	GetDir() string
	GetDefaultObject() interface{}
	GetObjectKeyFunc() cache.KeyFunc
	GetFileKeyFunc() FileIndexFunc
	GetObjectString() string
}

type FileConfigMapDeps struct {
	PathDir     string
	FileKeyFunc FileIndexFunc
}

func NewDefaultFileConfigMapDeps(dir string) *FileConfigMapDeps {
	return &FileConfigMapDeps{
		PathDir: dir,
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

func NewDefaultFileServiceDeps(dir string) *FileServiceDeps {
	return &FileServiceDeps{
		PathDir: dir,
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

func NewDefaultFileNodeDeps(dir string) *FileNodeDeps {
	return &FileNodeDeps{
		PathDir: dir,
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

func NewDefaultFileSecretDeps(dir string) *FileSecretDeps {
	return &FileSecretDeps{
		PathDir: dir,
	}
}

type FilePodDeps struct {
	PathDir     string
	FileKeyFunc FileIndexFunc
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

func NewDefaultFilePodDeps(dir string) *FilePodDeps {
	return &FilePodDeps{
		PathDir: dir,
	}
}

type FileLeaseDeps struct {
	PathDir     string
	FileKeyFunc FileIndexFunc
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

func NewDefaultFileLeaseDeps(dir string) *FileLeaseDeps {
	return &FileLeaseDeps{
		PathDir: dir,
	}
}

type FileEventDeps struct {
	PathDir     string
	FileKeyFunc FileIndexFunc
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

func NewDefaultFileEventDeps(dir string) *FileEventDeps {
	return &FileEventDeps{
		PathDir: dir,
	}
}

var _ FileObjectDeps = &FileConfigMapDeps{}
var _ FileObjectDeps = &FileServiceDeps{}
var _ FileObjectDeps = &FileNodeDeps{}
var _ FileObjectDeps = &FileSecretDeps{}
var _ FileObjectDeps = &FilePodDeps{}
var _ FileObjectDeps = &FileLeaseDeps{}
var _ FileObjectDeps = &FileEventDeps{}
