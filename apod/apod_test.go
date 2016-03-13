package apod

import (
	"testing"
	"time"
)

func TestGetEntry(t *testing.T) {
	var apiKey = DefaultAPIKey
	var tests = []struct {
		apiKey     string
		apiURL     string
		date       time.Time
		requestHD  bool
		shouldPass bool
	}{
		{apiKey, DefaultAPIURL, time.Now(), true, true},
		{apiKey, DefaultAPIURL, time.Now(), false, true},
		{"", DefaultAPIURL, time.Now(), true, false},
		{apiKey, "http://127.0.0.1", time.Now(), DefaultRequestHD, false},
		{apiKey, "#http#:#//#127#.0#.0#.1#:", time.Now(), DefaultRequestHD, false},
	}

	for index, test := range tests {
		_, err := GetEntry(test.apiKey, test.apiURL, test.date, test.requestHD)
		if err != nil && test.shouldPass {
			t.Errorf("GetEntry test %d failed but should have passed.", index+1)
		} else if err == nil && !test.shouldPass {
			t.Errorf("GetEntry test %d passed but should have failed.", index+1)
		}
	}
}

func TestGetEntryForDate(t *testing.T) {
	var apiKey = DefaultAPIKey
	var tests = []struct {
		apiKey     string
		date       time.Time
		shouldPass bool
	}{
		{apiKey, time.Now(), true},
		{apiKey, time.Now().Add(24 * time.Hour), false},
		{apiKey, time.Date(1970, 1, 1, 0, 0, 0, 0, time.Local), false},
	}

	for index, test := range tests {
		_, err := GetEntryForDate(test.apiKey, test.date)
		if err != nil && test.shouldPass {
			t.Errorf("GetEntryForDate test %d failed but should have passed.", index+1)
		} else if err == nil && !test.shouldPass {
			t.Errorf("GetEntryForDate test %d passed but should have failed.", index+1)
		}
	}
}

func TestGetEntryForToday(t *testing.T) {
	if _, err := GetEntryForToday(DefaultAPIKey); err != nil {
		t.Error("GetEntryForToday test should have passed." + err.Error())
	}
}

func TestDownloadAPOD(t *testing.T) {
	apod, err := GetEntryForToday(DefaultAPIKey)
	if err != nil {
		t.Errorf("Unable to get entry for today. %s", err.Error())
	}

	path := "./downloads"
	_, err = DownloadAPOD(apod, path, true)
	if err != nil {
		t.Errorf("Unable to download file to path %s. %s", path, err.Error())
	}
}
