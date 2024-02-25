package main

import (
	"database/sql"

	live "insiderviz.com/crawlers/live_crawler"
	local "insiderviz.com/crawlers/local_crawler"
)

// This file can run any crawler, should mainly be used for local testing as crawlers should be deployed indiviually
func main() {
	// Creates connection to local db called insiderviz-crawler2 with username and password "root"
	conn := utils.CreateMySQLConnection("root:root@tcp(127.0.0.1:3306)/insiderviz-crawler") 
	defer conn.Close()
	// cik := "0001378992"
	// issuer_utils.CrawlIssuerJSON(cik)
	// test_live(conn)
	test_local(conn)
}

func test_live(conn *sql.DB) {
	live.RunLiveCrawler(conn)
}

func test_local(conn *sql.DB) {
	folder := "./submissions/2024-01-27"

	offset := 0 // start at the first company
	stride := 1 // only go one at a time
	start := "2023-01-16" //start 
	end := "3000-01-01" // end, arbitrary date in future means get everything

	local.RunLocalCrawler(folder,start,end,offset,stride,conn)
}

var r = InitRequestGuard()

func InitRequestGuard()  chan  struct{} {
	r := make(chan struct{},1)
	r<- struct{}{}
	return r
}


// Can channels be global variables or must they be passed as arguments to go routines??
// func global_chan_test() {
// 	print("start")
// 	var counter int32 = 0
// 	var wg sync.WaitGroup
// 	for i := 0; i< 100;i++ {
// 		wg.Add(1)
// 		go func() {
// 			defer wg.Done()
// 			<-r
// 			go func(counter *int32) {
// 				// creates a new go function that lets another start after x milliseconds
// 				go func() {
// 					time.Sleep(10 * time.Millisecond)
// 					r<-struct{}{}
// 				}()
// 				atomic.AddInt32(counter, 1)
// 				time.Sleep(100 * time.Millisecond)
// 				print("thread finished")
				
// 			}(&counter)
// 		}()
// 	}
// 	wg.Wait()
// 	print(counter)
// }


