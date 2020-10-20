package database

import (
	"testing"
)

func TestCreateCache(t *testing.T) {
	validData := Data{
		OK:         true,
		ShortURL:   "ShortURL",
		FullURL:    "https://www.FullURL.com/",
		ViewsCount: 123,
	}
	err := Init()
	expect(t, err, nil)

	err = AddCache(validData)
	expect(t, err, nil)
}

func TestGetCache(t *testing.T) {
	validData := Data{
		OK:         true,
		ShortURL:   "ShortURL",
		FullURL:    "https://www.FullURL.com/",
		ViewsCount: 123,
	}
	err := Init()
	expect(t, err, nil)

	newData := CheckCache(validData.ShortURL, NotViewed)
	expect(t, newData.ShortURL, validData.ShortURL)

	newData = CheckCache(validData.ShortURL, IsViewed)
	expect(t, newData.ViewsCount, validData.ViewsCount+1)
}
