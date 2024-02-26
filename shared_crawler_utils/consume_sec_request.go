package shared_crawler_utils

import "time"

// Global channel is a thread safe queue that each worker waits on
var request_guard = InitRequestGuard()

func InitRequestGuard()  chan  struct{} {
	request_guard := make(chan struct{},1)
	request_guard<- struct{}{}
	return request_guard
}
/*
	Must be called before every SEC request with the global request_guard chan
	Ensures that the 10 requests/second rule isn't violated

*/
func ConsumeSECRequest() {
	print(&request_guard,"\n")
	// Waits for a request to be available
	<-request_guard

	// Makes a request available after 110 milliseconds
	go func(request_guard chan struct{}) {
		// Theoretically could be 100 ms, but this causes 429, too many request errors
		time.Sleep(110 * time.Millisecond)
		request_guard<-struct{}{}
	}(request_guard)
}