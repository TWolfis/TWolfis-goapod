package GoApod

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
)

/*
APOD or Astronomical Picture of the Day is a picture or video uploaded by NASA every single day
This file contains the code to work with this API using Go
*/

const (
	Scheme   = "https"
	Host     = "api.nasa.gov"
	apodPath = "/planetary/apod"
)

var (
	// ApiKey can be overwritten with the users personal ApiKey
	ApiKey = "DEMO_KEY"
)

// APOD contains the variables associated with the APOD API
type APOD struct {

	// Date of the APOD image to retrieve defaults to today
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

	// Thumbs Return the URL of video thumbnail. If an APOD is not a video, this parameter is ignored.
	// defaults to false
	Thumbs bool

	Response  APODResponse
	Responses []APODResponse
}

// APODResponse holds the response JSON object received by a successful call to the APOD API
type APODResponse struct {
	Date           string `json:"date"`
	Explanation    string `json:"explanation"`
	Hdurl          string `json:"hdurl"`
	MediaType      string `json:"media_type"`
	ServiceVersion string `json:"service_version"`
	Title          string `json:"title"`
	URL            string `json:"url"`
}

// composeQuery creates a query from a struct by marshalling it to json
func (a *APOD) composeQuery() (*http.Request, error) {
	//compose URL from path
	url := Scheme + "://" + Host + apodPath

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
	if a.Date != "" && a.StartDate == "" && a.EndDate == "" && a.Count == 0 {
		//If date is set and StartDate, EndDate and Count are not set
		//Add date to query
		query.Add("date", a.Date)

	} else if a.Date == "" && a.StartDate != "" && a.EndDate != "" && a.Count == 0 {
		//If Date is not set, StartDate and EndDate are set and Count is not set
		//Add StartDate and EndDate to query
		query.Add("start_date", a.StartDate)
		query.Add("end_date", a.EndDate)

	} else if a.Date == "" && a.StartDate == "" && a.EndDate == "" && a.Count > 0 {
		//If Date, StartDate and EndDate are not set but Count is set
		//Add Count to Query
		query.Add("count", strconv.Itoa(a.Count))
	}

	if a.Thumbs {
		//If Thumbs is set
		//Add Thumbs to query
		query.Add("thumbs", "True")
	}
	req.URL.RawQuery = query.Encode()
	return req, nil
}

func (a *APOD) MakeRequest() (string, []byte, error) {

	//create request
	request, err := a.composeQuery()
	if err != nil {
		return "", []byte{}, err
	}

	//create client and make request to the API
	client := &http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		return "", []byte{}, err
	}

	defer resp.Body.Close()

	reader, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", []byte{}, err
	}

	return request.URL.String(), reader, nil
}

// Unwrap takes the json object received from the API call and unwraps it to an array of APODResponse structs
func (a *APOD) Unwrap(resp []byte) error {

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

// DownloadImage downloads the APOD in either hd or normal definition
func DownloadImage(a APODResponse, hdURL bool) error {

	var src string

	if hdURL {
		if a.Hdurl != "" {
			src = a.Hdurl
		}
	} else if hdURL == false {
		if a.URL != "" {
			src = a.URL
		}
	} else {
		return errors.New("no url found")
	}

	directory, _ := os.Getwd()
	filename := a.Title
	path := directory + "/" + filename + ".jpg"

	fmt.Println("Downloading image to: ", path)
	//make request to src to fetch the file
	resp, err := http.Get(src)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	//create an empty file to write the jpg to
	file, err := os.Create(path)
	if err != nil {
		return err
	}

	defer file.Close()

	//write bytes to jpg file
	_, err = io.Copy(file, resp.Body)
	if err != nil {
		return err
	}

	return nil
}
