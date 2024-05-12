package live_crawler

import (
	"context"
	"database/sql"
	"fmt"
	"sync"
	"time"

	"github.com/go-co-op/gocron"
	. "insiderviz.com/crawlers/utils"
	issuer_utils "insiderviz.com/crawlers/utils/issuer"
)

/*
	Runs a live crawl of the SEC RSS feed into the provided @conn sql DB
*/
func LiveCrawl(conn *sql.DB) bool{
	
	//Context allows the SaveForm thread to cancel the LoadForm thread when a duplicate is found
	ctx, cancel := context.WithCancel(context.Background())
	// create the channel to pass through data
	data := make(chan RawForm4)


	// Reads issuer_ciks currently in db and stores them in @seen_ciks map
	seen_ciks := make(map[string]struct{},0)
	issuer_sql := `select cik from issuer`
	result, err := conn.Query(issuer_sql)
	if err != nil {
		panic(err)
	}
	for result.Next() {
		var cik string
		result.Scan(&cik)
		seen_ciks[cik] = struct{}{}
	}

	// Creates second thread for getForm4RSS
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		getForm4RSS(data,ctx)

	}()

	// When SaveForm finds a duplicate getForm4RSS is stopped by cancel()
	duplicate := false
	for form := range data {
		// Must crawl issuer first if it is new, this should be very rare
		_, cik_already_saved := seen_ciks[form.IssuerCIK]
		if !cik_already_saved {
			fmt.Println("NEW CIK FOUND SAVING ",form.IssuerCIK)
			// Crawls the issuer that isn't in the db including all of its stock data
			issuer_utils.CrawlIssuersByCIK(conn,[]string{form.IssuerCIK},true,"0001-01-01",1)
		}
		//Saves form to db terminating when the form is a duplicate
		if SaveForm(conn, form) {
			fmt.Println("Duplicate found breaking")
			duplicate = true
			cancel()
		}
	}
	cancel()

	// Both Loading and Saving are finished
	wg.Wait()
	fmt.Println("Saving Issuers")
	


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
	Every 5 minutes from 6am-10pm monday through friday

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
	s.Cron("*/5 6-22 * * 1-5").Do(func() {
		LiveCrawl(conn)
	})
	s.StartBlocking()
}
