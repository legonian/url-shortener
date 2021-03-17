package database

import (
	"testing"
)

func TestCache(t *testing.T) {
	CACHE_LIMIT = 2
	testData := Data{
		OK:         true,
		FullURL:    "https://www.example.com/",
		ViewsCount: 1111,
	}
	shortcuts := []string{"123", "456", "789"}

	if err := Init(); err != nil {
		t.Fatalf("error when initializing database: %s", err)
	}

	testData.ShortURL = shortcuts[0]
	AddCache(testData)

	cacheData := CheckCache(shortcuts[0], NotViewed)
	if cacheData.ShortURL != testData.ShortURL {
		t.Errorf("cached ShortURL = %v, expect %v",
			cacheData.ShortURL, testData.ShortURL)
	}

	cacheData = CheckCache(shortcuts[0], IsViewed)
	if cacheData.ViewsCount != testData.ViewsCount+1 {
		t.Errorf("cache not increment views: ViewsCount = %v, expect %v",
			cacheData.ViewsCount, testData.ViewsCount+1)
	}

	testData.ShortURL = shortcuts[1]
	AddCache(testData)

	testData.ShortURL = shortcuts[2]
	AddCache(testData)

	cacheData = CheckCache(shortcuts[0], IsViewed)
	if cacheData.OK {
		t.Errorf("first cached data not cleared then reach the limit")
	}

	cacheData = CheckCache(shortcuts[1], IsViewed)
	if cacheData.OK {
		t.Errorf("second cached data not cleared then reach the limit")
	}

	cacheData = CheckCache(shortcuts[2], IsViewed)
	if cacheData.OK != true {
		t.Errorf("third cached data not accessed from cache")
	}
}
