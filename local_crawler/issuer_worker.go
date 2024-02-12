package local_crawler

import (
	"database/sql"
	"fmt"
	"sync/atomic"
	"time"

	. "insiderviz.com/crawlers/shared_crawler_functions"
)

/*
Does the work of loading and saving one issuers worth of forms, while respecting the global 10 request a second to the sec limit
*/
func issuerWorker(forms []FormJsonEntry, request_guard chan struct{}, conn *sql.DB, counter *int32, start_time int64, userAgent string) {

	for _, form := range forms {

		//Makes request to the SEC and gets rawform4 back, which has all needed info
		xml,err := FromURLLoadForm4XML(form.Url, form.AccNum, userAgent, request_guard)

		//Uses rawForm4 struct to populate database
		if err != nil {
			print(err)
		} else {
			SaveForm(conn, xml)
			atomic.AddInt32(counter, 1)
			elapsed_time := float32(time.Now().UnixMilli()-start_time) / 1000.0
			fmt.Printf("Completed %d forms in %0.2f seconds %0.2fs per form, URL: %s\n", *counter, elapsed_time, elapsed_time/float32(*counter), form.Url)
		}

	}
}
