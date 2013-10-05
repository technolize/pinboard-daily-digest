package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

const BaseUrl = "https://api.pinboard.in/v1/posts"
const Token = "YOUR_TOKEN"

type Result struct {
	Posts []Post
}

type Post struct {
	Description string
	Href        string
	Tags        string
}

func digest(date time.Time) error {
	var err error

	env := Environ()
	if env["PINBOARD_TOKEN"] == "" {
		fmt.Println("ENV PINBOARD_TOKEN not found.")
		os.Exit(1)
	}

	v := url.Values{}
	v.Set("dt", date.Format("2006-01-02"))
	v.Set("format", "json")
	v.Set("auth_token", env["PINBOARD_TOKEN"])

	url := fmt.Sprintf("%s/%s?%s", BaseUrl, "get", v.Encode())
	resp, err := http.Get(url)
	if err != nil {
		return err
	}

	defer resp.Body.Close()
	var body []byte
	if body, err = ioutil.ReadAll(resp.Body); err != nil {
		return err
	}

	var res Result
	if err := json.Unmarshal(body, &res); err != nil {
		return err
	}

	fmt.Printf("%s\n\n", date.Format(time.ANSIC))
	for _, p := range res.Posts {
		tag := ""
		if p.Tags != "" {
			tag = "tags: " + strings.Join(strings.Split(p.Tags, " "), ", ") + "\n"
		}
		fmt.Printf("%s\n%s\n%s\n", p.Description, p.Href, tag)
	}

	return err
}

func today() error {
	y := time.Now().AddDate(0, 0, -1)
	return digest(y)
}

func usage() {
	fmt.Printf(`usage of %s:
options:
    -d specific date digest
    -t today's digest
`, os.Args[0])
	os.Exit(1)
}

func errorHandler(e error) {
	fmt.Println(os.Stderr, "error: ", e)
	os.Exit(1)
}

func main() {
	var (
		err  error
		date string
		td   bool
	)

	flag.StringVar(&date, "d", "", "specific date digest")
	flag.BoolVar(&td, "t", false, "today's digest")
	flag.Parse()

	if len(date) > 0 {
		d, err := time.Parse("2006-01-02", date)
		if err != nil {
			errorHandler(err)
		}
		err = digest(d)
	} else if td {
		err = today()
	} else {
		usage()
	}

	if err != nil {
		errorHandler(err)
	}
}
