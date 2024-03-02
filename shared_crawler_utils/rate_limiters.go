package shared_crawler_utils

import "time"


func InitRequestGuard()  chan  struct{} {
	request_guard := make(chan struct{},1)
	request_guard<- struct{}{}
	return request_guard
}
/*
	Must be called before every SEC request with the global sec_request_guard chan
	Ensures that the 10 requests/second rule isn't violated

*/
var sec_request_guard = InitRequestGuard()
func ConsumeSECRequest() {
	// Waits for a request to be available
	<-sec_request_guard

	// Makes a request available after 110 milliseconds
	go func(sec_request_guard chan struct{}) {
		// Theoretically could be 100 ms, but this causes 429, too many request errors
		time.Sleep(106 * time.Millisecond)
		sec_request_guard<-struct{}{}
	}(sec_request_guard)
}

/*
	Must be called before every EOD request EOD gives us 1000/minute meaning 60 milliseconds per request
	Uses global channel to cause calls from different threads to wait
*/
var eod_request_guard = InitRequestGuard()
func ConsumeEODRequest() {
	// Waits for a request to be available
	<-eod_request_guard

	// Makes a request available after 110 milliseconds
	go func(eod_request_guard chan struct{}) {
		// Must be 60 milliseconds and adds 6 millisecond buffer to avoid differences in network latency causing a limit violation
		time.Sleep(66 * time.Millisecond)
		eod_request_guard<-struct{}{}
	}(eod_request_guard)
}