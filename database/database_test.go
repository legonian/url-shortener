package database

import (
	"testing"
)

func TestDatabase(t *testing.T) {
	validURL := "https://www.example.com/"

	if err := Init(); err != nil {
		t.Fatalf("error when initializing database: %s", err)
	}

	newData := CreateData(validURL)
	if newData.OK != true {
		t.Error("CreateData() with valid URL return error")
	}
	if newData.FullURL != validURL {
		t.Errorf("CreateData() with valid URL return data with FullURL = %s, "+
			"expect %s", newData.FullURL, validURL)
	}
	if newData.ViewsCount != 0 {
		t.Errorf("CreateData() with valid URL return data with ViewsCount = %s, "+
			"expect 0", newData.FullURL)
	}

	getData := GetData(newData.ShortURL, IsViewed)
	if getData.OK != true {
		t.Error("data from GetData() return error")
	}
	if getData.ShortURL != newData.ShortURL {
		t.Errorf("data from GetData() return ShortURL = %s, expect %s",
			newData.FullURL, validURL)
	}
	if getData.FullURL != newData.FullURL {
		t.Errorf("data from GetData() return FullURL = %s, expect %s",
			newData.FullURL, validURL)
	}
	if getData.ViewsCount != 1 {
		t.Errorf("data from GetData() return ViewsCount = %s, expect 1",
			newData.FullURL)
	}
}
