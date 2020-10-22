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

var (
	CACHE_DURATION_S, _ = time.ParseDuration("30s")
	CACHE_LIMIT         = 10
)

// Check does varieble with Item type is expired
func (item Item) isExpired() bool {
	if item.Expiration == 0 {
		return false
	}
	return item.Expiration < time.Now().UnixNano()
}

// Get data from cache
func CheckCache(key string, increment int) Data {
	store.mu.RLock()
	defer store.mu.RUnlock()

	if store.items[key].Content.OK == false {
		return Data{OK: false}
	}

	store.items[key] = Item{
		Content:    store.items[key].Content,
		Counter:    store.items[key].Counter + increment,
		Expiration: time.Now().Add(CACHE_DURATION_S).UnixNano(),
	}

	item := store.items[key]
	item.Content.ViewsCount += item.Counter
	log.Printf("-- CACHE GET %v    view_count += %v\n", key, increment)
	return item.Content
}

// Save data to cache
func AddCache(newData Data) {
	store.mu.RLock()
	defer store.mu.RUnlock()

	if CACHE_LIMIT <= len(store.items) {
		for key, item := range store.items {
			if 0 < item.Counter {
				GetData(item.Content.ShortURL, item.Counter)
			}
			delete(store.items, key)
			log.Printf("-- CACHE DELETED %s", key)
		}
	} else {
		for key, item := range store.items {
			isHasViewToSave := 0 < item.Counter
			expired := item.isExpired()

			if expired && isHasViewToSave {
				GetData(item.Content.ShortURL, item.Counter)
			}
			if expired {
				log.Printf("-- CACHE EXPIRED %s    with %vviews",
					item.Content.ShortURL,
					item.Counter)
				delete(store.items, key)
			}
		}
	}

	store.items[newData.ShortURL] = Item{
		Content:    newData,
		Counter:    0,
		Expiration: time.Now().Add(CACHE_DURATION_S).UnixNano(),
	}
	log.Printf("-- CACHE SET %s\n", newData.ShortURL)
}

// Remove all cache items, but if item hold non zero value save it to database
func ClearCache() {
	store.mu.RLock()
	defer store.mu.RUnlock()

	for key, item := range store.items {
		if 0 < item.Counter {
			GetData(item.Content.ShortURL, item.Counter)
		}
		delete(store.items, key)
		log.Printf("-- CACHE DELETED %s", key)
	}
}
