package apod

import (
	"testing"
	"time"
)

func TestGetEntry(t *testing.T) {
	var tests = []struct {
		apiURL     string
		date       time.Time
		requestHD  bool
		apiKey     string
		shouldPass bool
	}{
		{DefaultAPIURL, time.Now(), true, DefaultAPIKey, true},
		{DefaultAPIURL, time.Now(), false, DefaultAPIKey, true},
		{DefaultAPIURL, time.Now(), true, "", false},
		{"http://127.0.0.1", time.Now(), DefaultRequestHD, DefaultAPIKey, false},
		{"#http#:#//#127#.0#.0#.1#:", time.Now(), DefaultRequestHD, DefaultAPIKey, false},
	}

	for index, test := range tests {
		_, err := GetEntry(test.apiURL, test.date, test.requestHD, test.apiKey)
		if err != nil && test.shouldPass {
			t.Errorf("GetEntry test %d failed but should have passed.", index+1)
		} else if err == nil && !test.shouldPass {
			t.Errorf("GetEntry test %d passed but should have failed.", index+1)
		}
	}
}

func TestGetEntryForDate(t *testing.T) {
	var tests = []struct {
		date       time.Time
		shouldPass bool
	}{
		{time.Now(), true},
		{time.Now().Add(24 * time.Hour), false},
		{time.Date(1970, 1, 1, 0, 0, 0, 0, time.Local), false},
	}

	for index, test := range tests {
		_, err := GetEntryForDate(test.date)
		if err != nil && test.shouldPass {
			t.Errorf("GetEntryForDate test %d failed but should have passed.", index+1)
		} else if err == nil && !test.shouldPass {
			t.Errorf("GetEntryForDate test %d passed but should have failed.", index+1)
		}
	}
}

func TestGetEntryForToday(t *testing.T) {
	if _, err := GetEntryForToday(); err != nil {
		t.Error("GetEntryForToday test should have passed.")
	}
}
