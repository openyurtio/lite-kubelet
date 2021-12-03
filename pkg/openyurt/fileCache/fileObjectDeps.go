package fileCache

import (
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

func NewDefaultFileConfigMapDeps() *FileConfigMapDeps {
	return &FileConfigMapDeps{
		PathDir: "/etc/kubernetes/fileCache/configmaps",
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

func NewDefaultFileServiceDeps() *FileServiceDeps {
	return &FileServiceDeps{
		PathDir: "/etc/kubernetes/fileCache/services",
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

func NewDefaultFileNodeDeps() *FileNodeDeps {
	return &FileNodeDeps{
		PathDir: "/etc/kubernetes/fileCache/nodes",
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

func NewDefaultFileSecretDeps() *FileSecretDeps {
	return &FileSecretDeps{
		PathDir: "/etc/kubernetes/fileCache/secrets",
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

func NewDefaultFilePodDeps() *FilePodDeps {
	return &FilePodDeps{
		PathDir: "/etc/kubernetes/fileCache/pods",
	}
}

var _ FileObjectDeps = &FileConfigMapDeps{}
var _ FileObjectDeps = &FileServiceDeps{}
var _ FileObjectDeps = &FileNodeDeps{}
var _ FileObjectDeps = &FileSecretDeps{}
var _ FileObjectDeps = &FilePodDeps{}
