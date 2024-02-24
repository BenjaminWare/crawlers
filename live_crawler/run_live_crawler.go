package live_crawler

import (
	"context"
	"database/sql"
	"fmt"
	"sync"
	"time"

	"github.com/go-co-op/gocron"

	. "insiderviz.com/crawlers/shared_crawler_utils"
)

/*
	Runs a live crawl of the SEC RSS feed into the provided @conn sql DB
*/
func LiveCrawl(conn *sql.DB) bool{
	
	//Context allows the SaveForm thread to cancel the LoadForm thread when a duplicate is found
	ctx, cancel := context.WithCancel(context.Background())
	// create the channel to pass through data
	data := make(chan RawForm4)
	// Creates second thread for getForm4RSS
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		getForm4RSS(data, ctx)

	}()
	// When SaveForm finds a duplicate getForm4RSS is stopped by cancel()
	duplicate := false
	for form := range data {
		if SaveForm(conn, form) {
			fmt.Println("Duplicate found breaking")
			duplicate = true
			cancel()
		}
	}
	cancel()

	// Both Loading and Saving are finished
	wg.Wait()

	// A duplicate being found indicates success as there is an overlap between forms in the feed and the DB
	if !duplicate {
		fmt.Println("FAILURE: WE DIDN'T REACH ALL FORMS")
		SendFailureEmail("THE LIVE CRAWLER FAILED DESPAIR, THE DAY IT FAILED IS TODAY") 
	} else {
		fmt.Println("SUCCESS: WE GOT EVERY FORM")
	}
	return duplicate
}


/*
	Runs the live crawler against the https://www.sec.gov/cgi-bin/browse-edgar?action=getcurrent&CIK=&type=4&company=&dateb=&owner=only&start=0&count=100&output=atom rss feed
	Every 5 minutes from 12am-6pm monday through friday

	The live crawler finds, parses and stores all the new forms into the given sql db at @conn
	When a form is detected as already being in the db the crawling stops 
*/
func RunLiveCrawler(conn *sql.DB) {
	timezone, err := time.LoadLocation("America/New_York")
	if err != nil {
		panic(err)
	}

	s := gocron.NewScheduler(timezone)
	// Schedule the job to run every 5 minutes
	s.Cron("*/5 0-18 * * 1-5").Do(func() {
		LiveCrawl(conn)
	})
	s.StartBlocking()
}
