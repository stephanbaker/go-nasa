//Package apod assists in the retrieveal of data from NASA's
//Astronomy Picture of the Day (APOD) API.  For more information
//regarding this API, please see https://api.nasa.gov/api.html#apod
package apod

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
)

//APOD is a struct representing a single NASA astronomical photo of the day entry.
type APOD struct {
	Copyright      string `json:"copyright"`
	Date           string `json:"date"`
	Explanation    string `json:"explanation"`
	HDURL          string `json:"hdurl"`
	MediaType      string `json:"media_type"`
	ServiceVersion string `json:"service_version"`
	Title          string `json:"title"`
	URL            string `json:"url"`
}

//DefaultAPIURL represents the default base URL for this API
const DefaultAPIURL = "https://api.nasa.gov/planetary/apod"

//DefaultAPIKey represents the default API key for this APOD API
const DefaultAPIKey = "DEMO_KEY"

//DefaultRequestHD represents the default for whether or not to request hd urls
const DefaultRequestHD = true

//GetEntryForToday queries the NASA APOD API and returns the today's entry.
func GetEntryForToday() (*APOD, error) {
	return GetEntry(DefaultAPIURL, time.Now(), DefaultRequestHD, DefaultAPIKey)
}

//GetEntryForDate queries the NASA APOD API and returns the result for the given
//parameters.  This function expects the date of the entry to be provided.
func GetEntryForDate(date time.Time) (*APOD, error) {
	return GetEntry(DefaultAPIURL, date, DefaultRequestHD, DefaultAPIKey)
}

//GetEntry queries the NASA APOD API and returns the result for the given
//parameters.  This function expects the api url, date of the entry, whether or not to provide
//the high-defenition URL for entry, and the API key to be used for the request.
func GetEntry(apiurl string, date time.Time, hd bool, apikey string) (*APOD, error) {
	//Get our base url
	apodURL, err := url.Parse(apiurl)
	if err != nil {
		return nil, err
	}

	//Encode our query parameters
	parameters := url.Values{}
	parameters.Set("api_key", apikey)
	parameters.Set("date", date.Format("2006-01-02"))
	if hd {
		parameters.Set("hd", "true")
	} else {
		parameters.Set("hd", "false")
	}
	apodURL.RawQuery = parameters.Encode()

	//Form the request
	request, err := http.NewRequest("GET", apodURL.String(), nil)
	if err != nil {
		return nil, err
	}

	//Send the http request
	client := http.DefaultClient
	response, err := client.Do(request)
	if err != nil {
		return nil, err
	}

	if response.StatusCode != http.StatusOK {
		return nil, errors.New("Bad response status.")
	}

	//Read the response body
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	//Parse the JSON response
	var apod = &APOD{}
	err = json.Unmarshal(body, apod)
	if err != nil {
		return nil, err
	}

	return apod, err
}
