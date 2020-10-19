package handler

import (
	"database/sql"
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
	item.Content.ViewsCount += s.items[key].Counter
	if item.Expired() {
		delete(s.items, key)
	}
	return item.Content
}

func (s Storage) GetWithIncrement(key string, db *sql.DB) Data {
	s.mu.RLock()
	defer s.mu.RUnlock()
	s.items[key] = Item{
		Content:    s.items[key].Content,
		Counter:    s.items[key].Counter + 1,
		Expiration: s.items[key].Expiration,
	}

	item := s.items[key]
	if item.Expired() {
		delete(s.items, key)
		q := fmt.Sprintf("select * from get_full_url('%s',%v)", item.Content.ShortURL, item.Counter)
		getQuery(db, q)
	}
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
func CheckCache(short_url string, db *sql.DB) string {
	cache := store.GetWithIncrement(short_url, db)

	if cache.OK {
		log.Print("Cache get\n")
		return cache.FullURL
	}
	return ""
}

//Set data to cache
func AddToCache(newData Data) error {
	if d, err := time.ParseDuration(CACHE_DURATION); err == nil {
		log.Printf("> Cache new: for route /%s. Expired in %s\n", newData.ShortURL, CACHE_DURATION)
		store.Set(newData, d)
		return nil
	} else {
		log.Printf("> Cache err: %s\n", err)
		return err
	}
}
