package shared_crawler_functions

import (
	"testing"
)

// Tests an empty URL expects an error to be raised
func TestFromURLLoadForm4XML_EmptyUrl(t *testing.T) {
	url := ""
	acc_num := "mishaped"
	request_guard := make(chan struct{},1)
	request_guard <- struct{}{}
	_,err := FromURLLoadForm4XML(url,acc_num,"Ben benwareohio@gmail.com",request_guard)
    if err == nil{
        t.Fatalf("Parsed empty url without raising error")
    }
}

//Tests a valid form expects no error
func TestFromURLLoadForm4XML_ValidForm(t *testing.T) {
	url := "https://www.sec.gov/Archives/edgar/data/1939261/000159396824000179/0001593968-24-000179-index.htm"
	acc_num := "mishaped"
	request_guard := make(chan struct{},1)
	request_guard <- struct{}{}
	_,err := FromURLLoadForm4XML(url,acc_num,"Ben benwareohio@gmail.com",request_guard)
    if err != nil{
        t.Fatalf("Failed to parse form")
    }
}