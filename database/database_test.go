package database

import (
	"testing"
)

func TestDatabase(t *testing.T) {
	err := Init()
	expect(t, err, nil)

	validData := Data{
		OK:         true,
		ShortURL:   "ShortURL",
		FullURL:    "https://www.FullURL.com/",
		ViewsCount: 123,
	}

	newData := AddData(validData.FullURL)
	expect(t, newData.FullURL, validData.FullURL)

	getData := GetData(newData.ShortURL, IsViewed)
	expect(t, getData.OK, newData.OK)
	expect(t, getData.ShortURL, newData.ShortURL)
	expect(t, getData.FullURL, newData.FullURL)
	expect(t, getData.ViewsCount, newData.ViewsCount+1)
}

func expect(t *testing.T, varToTest interface{}, expected interface{}) {
	if varToTest != expected {
		t.Fatalf("Variable value is '%v', expected '%v'", varToTest, expected)
	}
}
