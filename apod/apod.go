//Package apod assists in the retrieveal of data from NASA's
//Astronomy Picture of the Day (APOD) API.  For more information
//regarding this API, please see https://api.nasa.gov/api.html#apod
package apod

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
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
const DefaultAPIKey = "rtgkG6UUPtPgOO7Zk9UfDQeE0WJEhT21EnhEqMKW" //"DEMO_KEY"

//DefaultRequestHD represents the default for whether or not to request hd urls
const DefaultRequestHD = true

//GetEntryForToday queries the NASA APOD API and returns the today's entry.
func GetEntryForToday(apiKey string) (*APOD, error) {
	return GetEntry(apiKey, DefaultAPIURL, time.Now(), DefaultRequestHD)
}

//GetEntryForDate queries the NASA APOD API and returns the result for the given
//parameters.  This function expects the date of the entry to be provided.
func GetEntryForDate(apiKey string, date time.Time) (*APOD, error) {
	return GetEntry(apiKey, DefaultAPIURL, date, DefaultRequestHD)
}

//GetEntry queries the NASA APOD API and returns the result for the given
//parameters.  This function expects the api url, date of the entry, whether or not to provide
//the high-defenition URL for entry, and the API key to be used for the request.
func GetEntry(apikey string, apiurl string, date time.Time, hd bool) (*APOD, error) {
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

//DownloadAPOD checks the provided APOD struct, and if possible downloads the
//contents of the picture and saves it to the designated destination path.
func DownloadAPOD(apod *APOD, destinationPath string, hd bool) (written int64, err error) {
	if apod == nil {
		return 0, errors.New("You must provide an apod entry.")
	}
	//Verify we are handling an image entry
	if apod.MediaType != "image" {
		return 0, errors.New("DownloadAPOD will only work for APOD entries with media_type=\"image\".")
	}

	//Determine the appropriate image url
	var imageURL string
	if hd && len(apod.HDURL) > 0 {
		imageURL = apod.HDURL
	} else if len(apod.URL) > 0 {
		imageURL = apod.URL
	} else {
		return 0, errors.New("No image url was provided for this APOD entry.")
	}

	//Parse the image name from the url
	tokens := strings.Split(imageURL, "/")
	fileName := tokens[len(tokens)-1]

	//Create the directory if it does not exist
	if _, err := os.Stat(destinationPath); err != nil {
		if err := os.MkdirAll(destinationPath, 0744); err != nil {
			return 0, errors.New("Unable to write to the specified file path.")
		}
	}
	destinationPath = filepath.Join(destinationPath, fileName)

	//Check if the destination path already exists
	if _, err := os.Stat(destinationPath); err == nil {
		return 0, fmt.Errorf("File already exists at path %s.", destinationPath)
	}

	resp, err := http.Get(imageURL)
	if err != nil {
		return 0, fmt.Errorf("Unable to get file at url %s. %s", imageURL, err.Error())
	}
	defer resp.Body.Close()

	file, err := os.Create(destinationPath)
	if err != nil {
		return 0, fmt.Errorf("Unable to open destination file at path %s.  %s", destinationPath, err.Error())
	}
	written, err = io.Copy(file, resp.Body)
	if err != nil {
		return 0, fmt.Errorf("Unable to write contents of file to path %s.  %s", destinationPath, err.Error())
	}

	return written, nil
}
