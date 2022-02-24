/*
Copyright 2015 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package oyImages

import (
	"k8s.io/kubernetes/pkg/kubelet/images"
	"sort"
	"sync"
	"time"

	v1 "k8s.io/api/core/v1"
	"k8s.io/klog/v2"

	"k8s.io/apimachinery/pkg/util/sets"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/tools/record"
	"k8s.io/kubernetes/pkg/kubelet/container"
	"k8s.io/kubernetes/pkg/kubelet/util/sliceutils"
)

// ImageGCPolicy is a policy for garbage collecting images. Policy defines an allowed band in
// which garbage collection will be run.
type ImageGCPolicy struct {
	// Minimum age at which an image can be garbage collected.
	MinAge time.Duration
}

type realImageGCManager struct {
	// Container runtime
	runtime container.Runtime

	// Records of images and their use.
	imageRecords     map[string]*imageRecord
	imageRecordsLock sync.Mutex

	// The image garbage collection policy in use.
	policy ImageGCPolicy

	// statsProvider provides stats used during image garbage collection.
	statsProvider images.StatsProvider

	// Recorder for Kubernetes events.
	recorder record.EventRecorder

	// Reference to this node.
	nodeRef *v1.ObjectReference

	// Track initialization
	initialized bool

	// imageCache is the cache of latest image list.
	imageCache imageCache

	// sandbox image exempted from GC
	sandboxImage string
}

// imageCache caches latest result of ListImages.
type imageCache struct {
	// sync.Mutex is the mutex protects the image cache.
	sync.Mutex
	// images is the image cache.
	images []container.Image
}

// set sorts the input list and updates image cache.
// 'i' takes ownership of the list, you should not reference the list again
// after calling this function.
func (i *imageCache) set(images []container.Image) {
	i.Lock()
	defer i.Unlock()
	// The image list needs to be sorted when it gets read and used in
	// setNodeStatusImages. We sort the list on write instead of on read,
	// because the image cache is more often read than written
	sort.Sort(sliceutils.ByImageSize(images))
	i.images = images
}

// get gets image list from image cache.
// NOTE: The caller of get() should not do mutating operations on the
// returned list that could cause data race against other readers (e.g.
// in-place sorting the returned list)
func (i *imageCache) get() []container.Image {
	i.Lock()
	defer i.Unlock()
	return i.images
}

// Information about the images we track.
type imageRecord struct {
	// Time when this image was first detected.
	firstDetected time.Time

	// Time when we last saw this image being used.
	lastUsed time.Time

	// Size of the image in bytes.
	size int64
}

// NewImageGCManager instantiates a new ImageGCManager object.
func NewImageGCManager(runtime container.Runtime, statsProvider images.StatsProvider,
	recorder record.EventRecorder, nodeRef *v1.ObjectReference, policy ImageGCPolicy, sandboxImage string) (images.ImageGCManager, error) {
	im := &realImageGCManager{
		runtime:       runtime,
		policy:        policy,
		imageRecords:  make(map[string]*imageRecord),
		statsProvider: statsProvider,
		recorder:      recorder,
		nodeRef:       nodeRef,
		initialized:   false,
		sandboxImage:  sandboxImage,
	}

	return im, nil
}

func (im *realImageGCManager) Start() {
	go wait.Until(func() {
		// Initial detection make detected time "unknown" in the past.
		var ts time.Time
		if im.initialized {
			ts = time.Now()
		}
		_, err := im.detectImages(ts)
		if err != nil {
			klog.Warningf("[imageGCManager] Failed to monitor images: %v", err)
		} else {
			im.initialized = true
		}
	}, 5*time.Minute, wait.NeverStop)

	// Start a goroutine periodically updates image cache.
	go wait.Until(func() {
		images, err := im.runtime.ListImages()
		if err != nil {
			klog.Warningf("[imageGCManager] Failed to update image list: %v", err)
		} else {
			im.imageCache.set(images)
		}
	}, 30*time.Second, wait.NeverStop)

}

// Get a list of images on this node
func (im *realImageGCManager) GetImageList() ([]container.Image, error) {
	return im.imageCache.get(), nil
}

func (im *realImageGCManager) detectImages(detectTime time.Time) (sets.String, error) {
	imagesInUse := sets.NewString()

	// Always consider the container runtime pod sandbox image in use
	imageRef, err := im.runtime.GetImageRef(container.ImageSpec{Image: im.sandboxImage})
	if err == nil && imageRef != "" {
		imagesInUse.Insert(imageRef)
	}

	images, err := im.runtime.ListImages()
	if err != nil {
		return imagesInUse, err
	}
	pods, err := im.runtime.GetPods(true)
	if err != nil {
		return imagesInUse, err
	}

	// Make a set of images in use by containers.
	for _, pod := range pods {
		for _, container := range pod.Containers {
			klog.V(5).Infof("Pod %s/%s, container %s uses image %s(%s)", pod.Namespace, pod.Name, container.Name, container.Image, container.ImageID)
			imagesInUse.Insert(container.ImageID)
		}
	}

	// Add new images and record those being used.
	now := time.Now()
	currentImages := sets.NewString()
	im.imageRecordsLock.Lock()
	defer im.imageRecordsLock.Unlock()
	for _, image := range images {
		klog.V(5).Infof("Adding image ID %s to currentImages", image.ID)
		currentImages.Insert(image.ID)

		// New image, set it as detected now.
		if _, ok := im.imageRecords[image.ID]; !ok {
			klog.V(5).Infof("Image ID %s is new", image.ID)
			im.imageRecords[image.ID] = &imageRecord{
				firstDetected: detectTime,
			}
		}

		// Set last used time to now if the image is being used.
		if isImageUsed(image.ID, imagesInUse) {
			klog.V(5).Infof("Setting Image ID %s lastUsed to %v", image.ID, now)
			im.imageRecords[image.ID].lastUsed = now
		}

		klog.V(5).Infof("Image ID %s has size %d", image.ID, image.Size)
		im.imageRecords[image.ID].size = image.Size
	}

	// Remove old images from our records.
	for image := range im.imageRecords {
		if !currentImages.Has(image) {
			klog.V(5).Infof("Image ID %s is no longer present; removing from imageRecords", image)
			delete(im.imageRecords, image)
		}
	}

	return imagesInUse, nil
}

func (im *realImageGCManager) GarbageCollect() error {

	// If over the max threshold, free enough to place us at the lower threshold.
	klog.Infof("[imageGCManager]: Trying to free useless image.")
	return im.freeSpace(time.Now())
}

func (im *realImageGCManager) DeleteUnusedImages() error {
	klog.Infof("attempting to delete unused images")
	return im.freeSpace(time.Now())
}

// Tries to free useless images on the disk.
//
func (im *realImageGCManager) freeSpace(freeTime time.Time) (error) {
	imagesInUse, err := im.detectImages(freeTime)
	if err != nil {
		return err
	}

	im.imageRecordsLock.Lock()
	defer im.imageRecordsLock.Unlock()

	// Get all images in eviction order.
	images := make([]evictionInfo, 0, len(im.imageRecords))
	for image, record := range im.imageRecords {
		if isImageUsed(image, imagesInUse) {
			klog.V(5).Infof("Image ID %s is being used", image)
			continue
		}
		images = append(images, evictionInfo{
			id:          image,
			imageRecord: *record,
		})
	}
	sort.Sort(byLastUsedAndDetected(images))

	// Delete unused images until we've freed up enough space.
	var deletionErrors []error
	for _, image := range images {
		klog.V(5).Infof("Evaluating image ID %s for possible garbage collection", image.id)
		// Images that are currently in used were given a newer lastUsed.
		if image.lastUsed.Equal(freeTime) || image.lastUsed.After(freeTime) {
			klog.V(5).Infof("Image ID %s has lastUsed=%v which is >= freeTime=%v, not eligible for garbage collection", image.id, image.lastUsed, freeTime)
			continue
		}

		// Avoid garbage collect the image if the image is not old enough.
		// In such a case, the image may have just been pulled down, and will be used by a container right away.

		if freeTime.Sub(image.firstDetected) < im.policy.MinAge {
			klog.V(5).Infof("Image ID %s has age %v which is less than the policy's minAge of %v, not eligible for garbage collection", image.id, freeTime.Sub(image.firstDetected), im.policy.MinAge)
			continue
		}

		// Remove image. Continue despite errors.
		klog.Infof("[imageGCManager]: Removing image %q to free %d bytes", image.id, image.size)
		err := im.runtime.RemoveImage(container.ImageSpec{Image: image.id})
		if err != nil {
			deletionErrors = append(deletionErrors, err)
			continue
		}
		delete(im.imageRecords, image.id)
	}

	return nil
}

type evictionInfo struct {
	id string
	imageRecord
}

type byLastUsedAndDetected []evictionInfo

func (ev byLastUsedAndDetected) Len() int      { return len(ev) }
func (ev byLastUsedAndDetected) Swap(i, j int) { ev[i], ev[j] = ev[j], ev[i] }
func (ev byLastUsedAndDetected) Less(i, j int) bool {
	// Sort by last used, break ties by detected.
	if ev[i].lastUsed.Equal(ev[j].lastUsed) {
		return ev[i].firstDetected.Before(ev[j].firstDetected)
	}
	return ev[i].lastUsed.Before(ev[j].lastUsed)
}

func isImageUsed(imageID string, imagesInUse sets.String) bool {
	// Check the image ID.
	if _, ok := imagesInUse[imageID]; ok {
		return true
	}
	return false
}
