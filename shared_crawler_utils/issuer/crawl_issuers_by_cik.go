package issuer

import (
	"database/sql"
	"fmt"
	"sync"
)

func CrawlIssuersByCIK(conn *sql.DB, ciks []string,num_threads int) {
	// Controls number of threads spun up
	thread_guard := make(chan struct{} , num_threads)
	// Makes sure the function waits until all routines finish
	var wg sync.WaitGroup

	for i, cik := range ciks {
		thread_guard <- struct{}{}
		wg.Add(1)
		go func(cik string, count int) {
			defer wg.Done()
			issuer := parseIssuerJSON(cik)
			// Uses 0 padded version of cik
			issuer.Cik = cik
			saveIssuer(conn, issuer)
			fmt.Printf("Finished%d/%d: %s\t%s\n",count,len(ciks),issuer.Name,issuer.Cik)
			<-thread_guard
		}(cik,i)

	}

	wg.Wait()
}