package database

import (
	"io/ioutil"
	"log"
	"testing"
)

var (
	newData   Data
	validData Data = Data{
		OK:         true,
		ShortURL:   "ShortURL",
		FullURL:    "https://www.FullURL.com/",
		ViewsCount: 123,
	}
	invalidData Data = Data{
		OK:         false,
		ShortURL:   "ShortURL",
		FullURL:    "https://www.FullURL.com/",
		ViewsCount: 123,
	}
)

func init() {
	log.SetOutput(ioutil.Discard)
}

func TestCreateCache(t *testing.T) {
	err := Model.Init()
	expect(t, err, nil)

	err = AddToCache(validData)
	expect(t, err, nil)
}

func TestGetCache(t *testing.T) {
	err := Model.Init()
	expect(t, err, nil)

	newData = CheckCache(validData.ShortURL, false)
	expect(t, newData.ShortURL, validData.ShortURL)

	newData = CheckCache(validData.ShortURL, true)
	expect(t, newData.ViewsCount, validData.ViewsCount+1)
}

func expect(t *testing.T, varToTest interface{}, expected interface{}) {
	if varToTest != expected {
		t.Fatalf("variable value is %v, expected %v", varToTest, expected)
	}
}
