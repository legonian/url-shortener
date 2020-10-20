package database

import (
	"log"
	"sync"
	"time"
)

// Cache data types
type (
	Item struct {
		Content    Data
		Counter    int
		Expiration int64
	}
	Storage struct {
		items map[string]Item
		mu    *sync.RWMutex
	}
)

var store = Storage{
	items: make(map[string]Item),
	mu:    &sync.RWMutex{},
}

const (
	CACHE_DURATION = "30s"
	CACHE_LIMIT    = 10
)

func (item Item) isExpired() bool {
	if item.Expiration == 0 {
		return false
	}
	return item.Expiration < time.Now().UnixNano()
}

// Get data from cache
func CheckCache(key string, increment int) Data {
	log.Printf("-- Cache getting... (increase view count by %v)\n", increment)
	store.mu.RLock()
	defer store.mu.RUnlock()

	if store.items[key].Content.OK == false {
		return Data{OK: false}
	}

	if 0 < increment {
		store.items[key] = Item{
			Content:    store.items[key].Content,
			Counter:    store.items[key].Counter + increment,
			Expiration: store.items[key].Expiration,
		}
	}

	item := store.items[key]
	if item.isExpired() {
		delete(store.items, key)
		item.Content = GetData(item.Content.ShortURL, item.Counter)
	}
	item.Content.ViewsCount += item.Counter
	return item.Content
}

// Save data to cache
func AddCache(newData Data) error {
	store.mu.RLock()
	defer store.mu.RUnlock()

	if CACHE_LIMIT < len(store.items) {
		store.items = make(map[string]Item)
	}
	for key, item := range store.items {
		isHasViewToSave := 0 < item.Counter

		if item.isExpired() && isHasViewToSave {
			GetData(item.Content.ShortURL, item.Counter)
		}
		if item.isExpired() {
			delete(store.items, key)
		}
	}

	duration, err := time.ParseDuration(CACHE_DURATION)
	if err == nil {
		store.items[newData.ShortURL] = Item{
			Content:    newData,
			Counter:    0,
			Expiration: time.Now().Add(duration).UnixNano(),
		}
		log.Printf("-- Cache set (%s, expired in %s)\n", newData.ShortURL, CACHE_DURATION)
	}
	return err
}
