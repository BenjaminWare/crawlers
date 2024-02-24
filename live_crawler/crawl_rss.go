package live_crawler

import (
	"net/http"
	"strconv"
	"time"

	"github.com/mmcdole/gofeed"
)

func crawlRSSFeed() []*gofeed.Feed {
	parser := gofeed.NewParser()
	client := &http.Client{}
	feeds := make([]*gofeed.Feed, 21)
	populateFeeds(feeds, parser, client)
	return feeds

}


func populateFeeds(feeds []*gofeed.Feed, parser *gofeed.Parser, client *http.Client) {
	
	limit := 100
	lastTime := time.Now().UnixMilli()
	timeCounter := 0

	i := 0
	for skip := 0; skip <= 2000; skip += limit {
		//TODO replace this with ConsumeSECRequest()
		if timeCounter == 9 {
			curTime := time.Now().UnixMilli()
			dif := curTime - lastTime
			if dif < 1000 {
				time.Sleep(time.Duration(1000-dif) * time.Millisecond)
			}
			lastTime = time.Now().UnixMilli()
			timeCounter = 0
		} else {
			timeCounter += 1
		}
		url := "https://www.sec.gov/cgi-bin/browse-edgar?action=getcurrent&CIK=&type=4&company=&dateb=&owner=only&start=" + strconv.Itoa(skip) + "&count=" + strconv.Itoa(limit) + "&output=atom"
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			panic(err)

		}

		req.Header.Set("User-Agent", "Davis bmdavis419@gmail.com")
		req.Host = "www.sec.gov"

		// Make the request using http.Client
		resp, err := client.Do(req)
		if err != nil {
			panic(err)
		}
		if resp.StatusCode == 429 {
			//sendFailureEmail("429 when getting rss!")
			panic("429 on rss")
		}

		// // Parse the response body using gofeed.Parser
		feed, err := parser.Parse(resp.Body)
		if err != nil {
			panic(err)
		}
		feeds[i] = feed
		i++

	}
}