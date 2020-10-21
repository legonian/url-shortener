package database

import (
	"testing"
)

func TestCreateCache(t *testing.T) {
	CACHE_LIMIT = 2
	shortcuts := []string{"123", "456", "789"}
	validData := Data{
		OK:         true,
		FullURL:    "https://www.FullURL.com/",
		ViewsCount: 1111,
	}
	err := Init()
	expect(t, err, nil)

	validData.ShortURL = shortcuts[0]
	AddCache(validData)

	newData := CheckCache(shortcuts[0], NotViewed)
	expect(t, newData.ShortURL, validData.ShortURL)

	newData = CheckCache(shortcuts[0], IsViewed)
	expect(t, newData.ViewsCount, validData.ViewsCount+1)

	validData.ShortURL = shortcuts[1]
	AddCache(validData)
	validData.ShortURL = shortcuts[2]
	AddCache(validData)

	newData = CheckCache(shortcuts[0], IsViewed)
	expect(t, newData.OK, false)
	newData = CheckCache(shortcuts[1], IsViewed)
	expect(t, newData.OK, false)
	newData = CheckCache(shortcuts[2], IsViewed)
	expect(t, newData.OK, true)
}
