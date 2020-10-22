package database

import (
	"testing"
)

func TestCacheFunctions(t *testing.T) {
	CACHE_LIMIT = 2
	shortcuts := []string{"123", "456", "789"}
	testData := Data{
		OK:         true,
		FullURL:    "https://www.example.com/",
		ViewsCount: 1111,
	}
	err := Init()
	expect(t, err, nil)

	testData.ShortURL = shortcuts[0]
	AddCache(testData)

	cacheData := CheckCache(shortcuts[0], NotViewed)
	expect(t, cacheData.ShortURL, testData.ShortURL)

	cacheData = CheckCache(shortcuts[0], IsViewed)
	expect(t, cacheData.ViewsCount, testData.ViewsCount+1)

	testData.ShortURL = shortcuts[1]
	AddCache(testData)
	testData.ShortURL = shortcuts[2]
	AddCache(testData)

	cacheData = CheckCache(shortcuts[0], IsViewed)
	expect(t, cacheData.OK, false)
	cacheData = CheckCache(shortcuts[1], IsViewed)
	expect(t, cacheData.OK, false)
	cacheData = CheckCache(shortcuts[2], IsViewed)
	expect(t, cacheData.OK, true)
}
