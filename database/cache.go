package database

import (
	"fmt"
	"log"
	"sync"
	"time"

	_ "github.com/lib/pq"
)

type (
	// Single cache data
	Item struct {
		Content    Data
		Counter    int
		Expiration int64
	}
	// All cache data
	Storage struct {
		items map[string]Item
		mu    *sync.RWMutex
	}
)

var Store = Storage{
	items: make(map[string]Item),
	mu:    &sync.RWMutex{},
}

const CACHE_DURATION = "20s"

func (item Item) Expired() bool {
	if item.Expiration == 0 {
		return false
	}
	return item.Expiration < time.Now().UnixNano()
}

func (s Storage) Get(key string, increment int) Data {
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
	item.Content.ViewsCount += item.Counter

	if item.Expired() {
		delete(s.items, key)
		q := fmt.Sprintf("select * from get_full_url('%s',%v)",
			item.Content.ShortURL,
			item.Counter,
		)
		item.Content = Model.GetQuery(q)
	}
	log.Printf("-- Cache Found: %v", item.Content)
	return item.Content
}

func (s Storage) Set(d Data, duration time.Duration) {
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
		cache := Store.Get(short_url, 1)
		return cache
	} else {
		cache := Store.Get(short_url, 0)
		return cache
	}
}

//Set data to cache
func AddToCache(newData Data) error {
	if d, err := time.ParseDuration(CACHE_DURATION); err == nil {
		log.Printf("-- Cache set (%s, expired in %s)\n", newData.ShortURL, CACHE_DURATION)
		Store.Set(newData, d)
		return nil
	} else {
		log.Printf("-- Cache err: %s\n", err)
		return err
	}
}
