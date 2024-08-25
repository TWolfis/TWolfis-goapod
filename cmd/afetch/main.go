package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/TWolfis/goapod"
)

var (
	date      string
	startdate string
	enddate   string
	download  bool
	dstFile   string
	hdurl     bool
)

func VerifyDate(date string) error {
	_, err := time.Parse("2006-1-1", date)
	return err
}

func PrintApod(a *goapod.Apod) {
	if len(a.Responses) == 0 {
		ar := a.Response
		fmt.Printf("Title: %s\nDate: %s\nExplanation: %s\nURL: %s",
			ar.Title, ar.Date, ar.Explanation, ar.URL)
	} else if len(a.Responses) > 1 {
		for _, ar := range a.Responses {
			fmt.Printf("Title: %s\nDate: %s\nExplanation: %s\nURL: %s\n\n",
				ar.Title, ar.Date, ar.Explanation, ar.URL)
		}
	} else {
		fmt.Println("No APOD for this date")
	}
}

func SaveImage(a *goapod.Apod, dstFile string, hdurl bool) error {
	if len(a.Responses) == 0 {
		ar := a.Response
		if dstFile == "" {
			dstFile = strings.ToLower(ar.Title) + ".jpg"
			img, err := ar.FetchImage(hdurl)
			if err != nil {
				return err
			}
			os.WriteFile(dstFile, img, 0644)
		}
	} else if len(a.Responses) > 1 {
		for _, ar := range a.Responses {
			dstFile = strings.ToLower(ar.Title) + ".jpg"
			img, err := ar.FetchImage(hdurl)
			if err != nil {
				return err
			}
			os.WriteFile(dstFile, img, 0644)
		}
	} else {
		return fmt.Errorf("no APOD for this date")
	}

	return nil
}

func main() {

	flag.StringVar(&date, "d", "", "date for APOD in format YYYY-MM-DD")
	flag.StringVar(&startdate, "sd", "", "start date for APOD range in format YYYY-MM-DD")
	flag.StringVar(&enddate, "ed", "", "start date for APOD range in format YYYY-MM-DD")
	flag.BoolVar(&download, "dl", download, "download image")
	flag.StringVar(&dstFile, "df", "", "destination file for downloaded image")
	flag.BoolVar(&hdurl, "hd", hdurl, "use hd image url")
	flag.Parse()

	a := goapod.Apod{}

	if date != "" {
		err := VerifyDate(date)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		a.Date = date
	} else if startdate != "" && enddate != "" {
		err := VerifyDate(startdate)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		a.StartDate = startdate

		err = VerifyDate(enddate)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		a.EndDate = enddate
	}

	a.Fetch()
	PrintApod(&a)

	// download image
	if download {
		SaveImage(&a, dstFile, hdurl)
	}
}
