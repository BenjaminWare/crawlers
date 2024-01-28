package live_crawler

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)


func GetXMLUri(link string, request_guard chan struct{}) (string, error) {
	// scrape the page for the XML uri
	// create the http request
	req, err := http.NewRequest("GET", link, nil)
	if err != nil {
		return "", err
	}

	req.Header.Set("User-Agent", "Ben maple6leaf@gmail.com")
	req.Host = "www.sec.gov"

	// create the http client
	client := &http.Client{}
	<-request_guard
	//After 110 milliseconds a new request is put on
	go func(request_guard chan struct{}) {
		time.Sleep(110 * time.Millisecond)
		request_guard <- struct{}{}
	}(request_guard)

	resp, err := client.Do(req)
	if resp.StatusCode == 429 {
		fmt.Printf("XML URI backed off 429")
	}
	if err != nil || resp.StatusCode != 200 {
		return "", err
	}

	// Load the HTML document
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return "", err
	}

	// find the link
	uri := ""
	doc.Find("a").Each(func(i int, s *goquery.Selection) {
		if strings.Contains(s.Text(), ".xml") {
			uri = "https://www.sec.gov" + s.AttrOr("href", "")
		}
	})

	return uri, nil
}
