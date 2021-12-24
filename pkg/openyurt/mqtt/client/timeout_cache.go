package client

import (
	"sync"
	"time"
)

type TimeoutCache struct {
	cache sync.Map
	cond  *sync.Cond
}

type cacheData struct {
	value *PublishAckData
	t     time.Time
}

var timecache *TimeoutCache

func init() {
	timecache = NewTimeoutCache(time.Second*5, time.Hour*6)
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

func (t *TimeoutCache) run(period, timeout time.Duration) {
	ticker := time.NewTicker(period)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			n := time.Now()
			t.cache.Range(func(key, value interface{}) bool {
				if n.Sub(value.(*cacheData).t) > timeout {
					t.cache.Delete(key)
				}
				return true
			})
			t.cond.Broadcast()
		}
	}
}

// Pop retrun false , when timeout
func (t *TimeoutCache) Pop(key string, timeout time.Duration) (*PublishAckData, bool) {
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

func (t *TimeoutCache) Set(value *PublishAckData) {
	d := &cacheData{
		value: value,
		t:     time.Now(),
	}
	t.cache.Store(value.Identity, d)
	t.cond.Broadcast()
}
