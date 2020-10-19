package handler

import (
	"fmt"
	"sync"
	"time"

	_ "github.com/lib/pq"
)

type (
	// Single cache data
	Item struct {
		Content    Data
		Expiration int64
	}
	// All cache data
	Storage struct {
		items map[string]Item
		mu    *sync.RWMutex
	}
)

var store = Storage{
	items: make(map[string]Item),
	mu:    &sync.RWMutex{},
}

const CACHE_DURATION = "20s"

func (item Item) Expired() bool {
	if item.Expiration == 0 {
		return false
	}
	return time.Now().UnixNano() > item.Expiration
}

func (s Storage) Get(key string) Data {
	s.mu.RLock()
	defer s.mu.RUnlock()

	item := s.items[key]
	if item.Expired() {
		delete(s.items, key)
		return Data{OK: false}
	}
	return item.Content
}

func (s Storage) Set(d Data, duration time.Duration) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.items[d.ShortURL] = Item{
		Content:    d,
		Expiration: time.Now().Add(duration).UnixNano(),
	}
}

// Get data from cache
func CheckCache(short_url string) string {
	content := store.Get(short_url)

	if content.OK {
		fmt.Print("Cache get\n")
		return content.FullURL
	}
	return ""
}

//Set data to cache
func AddToCache(newData Data) error {
	if d, err := time.ParseDuration(CACHE_DURATION); err == nil {
		fmt.Printf("> Cache new: for route /%s. Expired in %s\n", newData.ShortURL, CACHE_DURATION)
		store.Set(newData, d)
		return nil
	} else {
		fmt.Printf("> Cache err: %s\n", err)
		return err
	}
}
