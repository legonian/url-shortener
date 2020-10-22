package database

import (
	"testing"
)

func TestDatabase(t *testing.T) {
	validURL := "https://www.example.com/"

	err := Init()
	expect(t, err, nil)

	newData := CreateData(validURL)
	expect(t, newData.OK, true)
	expect(t, newData.FullURL, validURL)
	expect(t, newData.ViewsCount, 0)

	getData := GetData(newData.ShortURL, IsViewed)
	expect(t, getData.OK, newData.OK)
	expect(t, getData.ShortURL, newData.ShortURL)
	expect(t, getData.FullURL, newData.FullURL)
	expect(t, getData.ViewsCount, 1)
}

func expect(t *testing.T, varToTest interface{}, expected interface{}) {
	if varToTest != expected {
		t.Fatalf("Variable value is '%v', expected '%v'", varToTest, expected)
	}
}
