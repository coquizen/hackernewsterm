package hackernews

import (
	"github.com/patrickmn/go-cache"
	"strconv"
	"time"
)

type Store interface {
	Item(uint) (Item, error)
	// SetReadTimeStamp(uint)
	// GetReadTimeStamp(uint) time.Time
}

type store struct {
   firebase *FirebaseClient
}

type cachedStore struct {
	cache *cache.Cache
	firebase *FirebaseClient
}


func NewStore(firebaseClient *FirebaseClient) Store {
  return firebaseClient
}

func NewCachedStore(cache *cache.Cache, firebaseClient *FirebaseClient) Store {
	return &cachedStore{
		cache,
		firebaseClient,
	}
}
//
//func (s *store) Fetch(request RequestType) ([]Item, error) {
//
//}

func (s *cachedStore) Item(ID uint) (Item, error) {
	var key =  strconv.FormatUint(uint64(ID), 10)
	cachedItem, found := s.cache.Get(key)
	if found {
		return cachedItem.(Item), nil
	}
	live, err := s.firebase.Item(ID)
	if err != nil {
		return nil, err
	}
	s.cache.Set(key, live, 2* time.Minute)
	return live, nil
}

func (s *cachedStore) SetRead(id uint) {
	readState[id] = time.Now()
}

func (s *cachedStore) GetRead(id uint) time.Time {
	return readState[id]
}
