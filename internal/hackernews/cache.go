package hackernews

import (
	"errors"
	"github.com/patrickmn/go-cache"
	"sync"
	"time"
)

type ChanLocker struct {
	sync.Map
}


type Cache map[uint]Item

type readCache map[uint]time.Time

var (
	acquireLocker ChanLocker
	ErrCacheItemNotFound = errors.New("cache item not found")
	readState map[uint]time.Time
	_ = readCache{}
)

// NewCache creates a new cache with the specified duration for each datum stores within
func NewCache(expiryDuration time.Duration) *cache.Cache {
	return cache.New(expiryDuration, expiryDuration)
}


// Lock uses a sync.Map to ensure that the first goroutine to ask for a key will add a channel that all other goroutines will wait on. 
// See: https://lakefs.io/in-process-caching-in-go-scaling-lakefs-to-100k-requests-second/
func (c *ChanLocker) Lock(key interface{}, acquireFn func()) bool {
	waitCh := make(chan struct{})
	actual, locked := c.LoadOrStore(key, waitCh)
	if !locked {
		acquireFn()
		c.Delete(key)
		close(waitCh)
		return true
	}
	<- actual.(chan struct{})
	return false
}

func SetRead(id uint) {
  readState[id] = time.Now()
}


func GetRead(id uint) (time.Time, error) {
  lastRead, ok := readState[id]
  if !ok {
    return time.Time{},  ErrCacheItemNotFound
  }
  return lastRead, nil
}
//	tick := time.NewTicker(60 * time.Second)
//	defer tick.Stop()
//
//	updater := make(chan Item)
//	defer close (updater)
//
//	mutex = sync.RWMutex{}
//
//	go startUpdateDaemon(tick)
//
//
//}

//func startUpdateDaemon(ticker *time.Ticker, updater chan Item) {
//	handler := NewHandlerWithDefaultConfig(context.Background())
//
//	for {
//		select {
//		case <- ticker.C:
//			mutex.Lock()
//			for item := range handler.Subscribe(NewStories) {
//				NewStoriesCache.Items[item.ID()] = item
//			}
//			mutex.Unlock()
//			case <- updater:
//				return
//
//		}
//	}
//}
