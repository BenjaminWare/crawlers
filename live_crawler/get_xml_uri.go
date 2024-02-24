package live_crawler

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
	crawler_utils "insiderviz.com/crawlers/shared_crawler_utils"
)


func GetXMLUri(link string) (string, error) {
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
	crawler_utils.ConsumeSECRequest()

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
