/*
Copyright 2022 The OpenYurt Authors.

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
package message

import (
	"sync"
	"time"

	"k8s.io/klog/v2"
)

type TimeoutCache struct {
	cache sync.Map
	cond  *sync.Cond
}

type cacheData struct {
	value *AckData
	t     time.Time
}

var timecache *TimeoutCache

func init() {
	timecache = NewTimeoutCache(time.Second*5, time.Hour*2)
}

func GetDefaultTimeoutCache() *TimeoutCache {
	return timecache
}

func NewTimeoutCache(period, timeout time.Duration) *TimeoutCache {
	t := &TimeoutCache{
		cache: sync.Map{},
		cond:  sync.NewCond(&sync.Mutex{}),
	}
	go t.run(period, timeout)
	return t
}

// Expired keys in the cache are cleared every period and broadcast every period to prevent excessive data in the cache and reduce memory usage
func (t *TimeoutCache) run(period, timeout time.Duration) {
	ticker := time.NewTicker(period)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			n := time.Now()
			t.cache.Range(func(key, value interface{}) bool {
				if n.Sub(value.(*cacheData).t) > timeout {
					klog.Infof("key %s data %d now %d exceed %d s", key, value.(*cacheData).t.Unix(), n.Unix(), timeout.Milliseconds()/1000)
					t.cache.Delete(key)
				}
				return true
			})
			t.cond.Broadcast()
		}
	}
}

// PopWait retrun false , when timeout
// Get the key value, return true if get value, and delete it from the cache
// If not, wait until the corresponding key and value are obtained
// If no data is retrieved until timeout, return false
func (t *TimeoutCache) PopWait(key string, timeout time.Duration) (*AckData, bool) {
	ticker := time.NewTicker(timeout)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			return nil, false
		default:
		}
		v, ok := t.cache.Load(key)
		if ok {
			t.cache.Delete(key)
			return v.(*cacheData).value, true
		}
		t.cond.L.Lock()
		t.cond.Wait()
		t.cond.L.Unlock()
	}
}

func (t *TimeoutCache) Set(value *AckData) {
	d := &cacheData{
		value: value,
		t:     time.Now(),
	}
	t.cache.Store(value.Identity, d)
	t.cond.Broadcast()
}
