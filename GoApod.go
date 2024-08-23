package goapod

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strconv"
)

/*
Apod or Astronomical Picture of the Day is a picture or video uploaded by NASA every single day
This file contains the code to work with this API using Go
*/

const (
	Scheme   = "https"
	Host     = "api.nasa.gov"
	ApodPath = "/planetary/apod"
)

var (
	// ApiKey can be overwritten with the users personal ApiKey
	ApiKey = "DEMO_KEY"
)

// Apod contains the variables associated with the APOD API
type Apod struct {

	// Date of the Apod image to retrieve defaults to today
	Date string

	// StartDate of a date range, when requesting for a range of dates. Cannot be used with Date
	// defaults to none, cannot be used with date
	StartDate string

	// EndDate The end of the date range, when used with StartDate.
	// defaults to today
	EndDate string

	// Count If this is specified then count randomly chosen images will be returned.
	// Cannot be used with date or StartDate and EndDate.
	// defaults to none
	Count int

	// Thumbs Return the URL of video thumbnail. If an Apod is not a video, this parameter is ignored.
	// defaults to false
	Thumbs bool

	Response  ApodResponse
	Responses []ApodResponse
}

// ApodResponse holds the response JSON object received by a successful call to the APOD API
type ApodResponse struct {
	Date           string `json:"date"`
	Explanation    string `json:"explanation"`
	Hdurl          string `json:"hdurl"`
	MediaType      string `json:"media_type"`
	ServiceVersion string `json:"service_version"`
	Title          string `json:"title"`
	URL            string `json:"url"`
}

// composeQuery creates a query from a struct by marshalling it to json
func (a *Apod) composeQuery() (*http.Request, error) {
	url := "https://api.nasa.gov/planetary/apod"

	//	create get request, returns request object
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return &http.Request{}, err
	}

	// query can be used to add query options
	query := req.URL.Query()

	// add api_key to query
	query.Add("api_key", ApiKey)

	// handle various possible query parameters
	switch {
	case a.Date != "" && a.StartDate == "" && a.EndDate == "" && a.Count == 0:
		// If date is set and StartDate, EndDate, and Count are not set
		query.Add("date", a.Date)

	case a.Date == "" && a.StartDate != "" && a.EndDate != "" && a.Count == 0:
		// If Date is not set, StartDate and EndDate are set, and Count is not set
		query.Add("start_date", a.StartDate)
		query.Add("end_date", a.EndDate)

	case a.Date == "" && a.StartDate == "" && a.EndDate == "" && a.Count > 0:
		// If Date, StartDate, and EndDate are not set, but Count is set
		query.Add("count", strconv.Itoa(a.Count))
	}

	if a.Thumbs {
		// If Thumbs is set add Thumbs to query
		query.Add("thumbs", "True")
	}
	req.URL.RawQuery = query.Encode()
	return req, nil
}

// unwrap takes the json object received from the API call and unwraps it to an array of ApodResponse structs
func (a *Apod) unwrap(resp []byte) error {

	//unmarshal the received json objects
	err := json.Unmarshal(resp, &a.Response)
	if err != nil {
		err = json.Unmarshal(resp, &a.Responses)
		if err != nil {
			return err
		}
	}
	return nil
}

func (a *Apod) Fetch() error {

	// Create request
	request, err := a.composeQuery()
	if err != nil {
		return err
	}

	// Create client and make request to the API
	client := &http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	reader, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	err = a.unwrap(reader)

	return err
}

// FetchImage downloads the Apod Image in either hd or normal definition
// if hdurl is set but not available the function will default to url
func (a *ApodResponse) FetchImage(hdurl bool) ([]byte, error) {
	var src string
	if a.MediaType != "image" {
		return nil, errors.New("apodResponse is not an image")
	}

	if a.Hdurl != "" && hdurl {
		src = a.Hdurl
	} else {
		src = a.URL
	}

	//make request to src to fetch the file
	resp, err := http.Get(src)
	if err != nil {
		return []byte{}, err
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	return body, err
}

// String returns a string representation of the ApodResponse
func (a *ApodResponse) String() string {
	return "Title: " + a.Title + "\nDate: " + a.Date + "\nExplanation: " + a.Explanation + "\nURL: " + a.URL
}
