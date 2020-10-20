package database

import (
	"fmt"
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
	CACHE_DURATION = "20s"
	// CACHE_LIMIT    = 2
)

func (item Item) Expired() bool {
	if item.Expiration == 0 {
		return false
	}
	return item.Expiration < time.Now().UnixNano()
}

func (s Storage) cache(key string, increment int) Data {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if s.items[key].Content.OK == false {
		log.Printf("-- Cache not found: %s", key)
		return Data{OK: false}
	}

	if 0 < increment {
		s.items[key] = Item{
			Content:    s.items[key].Content,
			Counter:    s.items[key].Counter + increment,
			Expiration: s.items[key].Expiration,
		}
	}

	item := s.items[key]
	if item.Expired() {
		delete(s.items, key)
		q := fmt.Sprintf("select * from get_full_url('%s',%v)",
			item.Content.ShortURL,
			item.Counter,
		)
		item.Content = Model.GetQuery(q)
	}

	item.Content.ViewsCount += item.Counter
	log.Printf("-- Cache Found: %v", item.Content)
	return item.Content
}

func (s Storage) setCache(d Data, duration time.Duration) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.items[d.ShortURL] = Item{
		Content:    d,
		Counter:    0,
		Expiration: time.Now().Add(duration).UnixNano(),
	}
}

// Get data from cache
func CheckCache(short_url string, isIter bool) Data {
	log.Printf("-- Cache getting... (iterating? = %v)\n", isIter)
	if isIter {
		cache := store.cache(short_url, 1)
		return cache
	} else {
		cache := store.cache(short_url, 0)
		return cache
	}
}

//Set data to cache
func AddCache(newData Data) error {
	refreshCache()
	if d, err := time.ParseDuration(CACHE_DURATION); err == nil {
		log.Printf("-- Cache set (%s, expired in %s)\n", newData.ShortURL, CACHE_DURATION)
		store.setCache(newData, d)
		return nil
	} else {
		log.Printf("-- Cache err: %s\n", err)
		return err
	}
}

func refreshCache() {
	store.mu.RLock()
	defer store.mu.RUnlock()

	// if CACHE_LIMIT < len(store.items) {
	// 	log.Print("-- Cache limit: cleare all cache\n")
	// 	store.items = make(map[string]Item)
	// }

	for i, item := range store.items {
		isExpired := item.Expired()
		isHasViewToSave := 0 < item.Counter

		if isExpired && isHasViewToSave {
			q := fmt.Sprintf("select * from get_full_url('%s',%v)",
				item.Content.ShortURL,
				item.Counter,
			)
			Model.GetQuery(q)
		}

		if isExpired {
			log.Printf("-- Cache %s expired and cleared\n", i)
			delete(store.items, i)
		}
	}
}
