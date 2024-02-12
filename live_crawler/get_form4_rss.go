package live_crawler

import (
	"context"
	"fmt"
	"strings"

	. "insiderviz.com/crawlers/shared_crawler_functions"
)

type formJsonEntry struct {
	Url    string `json:"url"`
	AccNum string `json:"acc_num"`
}

// Parses rss feed given date,skip and limit params
// Stores unique accNums in seenAccNums and returns an array of Entries that have accNum and URL for XML of forms, and a boolean indicating the current date is still valid
// If the date is invalid that means the RSS feed began returning values from the next day, we've run out of forms for today
func getForm4RSS(output chan RawForm4,ctx context.Context) {
	fmt.Println("Get Form 4 Start")

	// Ensures only 10 requests a second are made to SEC
	request_guard := make(chan struct{} , 1)
	request_guard<-struct{}{}
	// Controls the number of FormWorkers that can run at a time, smaller number
	// thread_guard := 

	//Gets all 2000 entries in the RSS feed, this isn't 1-to-1 with forms as some acc_nums are in there multiple times
	fmt.Println("Crawl Feeds Start")
	feeds := crawlRSSFeed()
	fmt.Println("Crawl Feeds End")

	// Set to ensure only unique acc_nums are crawled
	seenAccNums := map[string]struct{}{}
	//The RSS feed requires us to get items in 21 blocks of 100
	for j := 0; j < 21; j++ {
		fmt.Printf("Parsing block %d\n",j*100)
		feed := feeds[j]
		for i := 0; i < len(feed.Items); i++ {
			// Gets link from form RSS entry
			item := feed.Items[i]
			link := item.Link
			formType := item.Categories[0] 
			linkParts := strings.Split(link, "/")
			link = ""
			for j := 0; j < len(linkParts)-1; j++ {
				link += linkParts[j] + "/"
			}

			// Gets accession number from RSS entry
			accNum := linkParts[7]
			accNumWithDashes := accNum[:10] + "-" + accNum[10:12] + "-" + accNum[12:]

			// Verifies the form in question is unique and Form4
			_, ok := seenAccNums[accNumWithDashes]
			if !ok && formType == "4" {
				seenAccNums[accNumWithDashes] = struct{}{}
					select {
						// Closes channel once all goroutines are done (wg.Wait()) and ctx is done
						case <-ctx.Done():
							close(output) 
							return
						default:
							FormWorker(accNumWithDashes,link,request_guard,output)
					}


			}

		}
	}
	close(output) 
	fmt.Println("Get Form 4 End")
}
