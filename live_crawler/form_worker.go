package live_crawler

import (
	. "insiderviz.com/crawlers/shared_crawler_utils"
)

/*
	Given the @acc_num,@link to a forms html page will parse the form into a RawForm4 and pass to the output channel
	@acc_num,@link are a forms acc_num and link
	@request_guard global channel used to prevent SEC from recieving more than 10 requests a second
	@output channel RawForm4 is placed into
*/
func FormWorker(acc_num string, link string,output chan RawForm4) {
	// Given a link to an index finds the correct XML link, requires a request to the SEC
	xmlUri, err := GetXMLUri(link)
	if err != nil {
		panic(err)
	}
	xml,err := FromURLLoadForm4XML(xmlUri,acc_num, "Ware benwareohio@gmail.com")

	if err == nil {
		output <- xml
	} else {
		print(err)
	}

}
