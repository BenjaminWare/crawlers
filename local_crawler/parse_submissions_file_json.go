package local_crawler

import (
	"encoding/json"
	"os"
	"strings"
)

type FormJsonEntry struct {
	Url    string `json:"url"`
	AccNum string `json:"acc_num"`
}

//Given a json filename and start/end dates
//If the file is an Issuer gets all acc_nums from between start/end
//If its a reporter returns empty
// Returns the forms to crawl and the cik if it is an issuer or empty string otherwise
func parseSubmissionsFileJSON(fileName, startDate string,endDate string) ([]FormJsonEntry, string) {
	typeOfEntity := "Both"
	entries := make([]FormJsonEntry, 0)


	// read the json file
	jsonFile, err := os.Open(fileName)
	if err != nil {
		panic(err)
	}
	defer jsonFile.Close()

	// decode the json file
	var apiResponse SecApiResponse
	decoder := json.NewDecoder(jsonFile)
	err = decoder.Decode(&apiResponse)
	if err != nil {
		return entries, ""
	}

	// setup the cik
	cik := apiResponse.Cik
	cik = strings.Repeat("0",10 - len(cik)) + cik

	// check the type of entity
	if apiResponse.IsIssuer == 1 && apiResponse.IsReporter == 0 {
		typeOfEntity = "Issuer"
	} else if apiResponse.IsIssuer == 0 && apiResponse.IsReporter == 1 {
		typeOfEntity = "Reporter"
	}

	if typeOfEntity == "Reporter" {
		return entries,""
	}

	// loop through the filings
	recent := apiResponse.Filings.Recent
	for i := range recent.AccNums {
		// make sure the date is large enough
		if recent.ReportDates[i] >= startDate  && recent.ReportDates[i] < endDate{
			// check the form type
			if recent.Forms[i] == "4" {
				// get the acc number for the url
				accNumUrl := strings.ReplaceAll(apiResponse.Filings.Recent.AccNums[i], "-", "")

				// get the primary document for the url
				primaryDoc := strings.SplitAfter(recent.PrimaryDocs[i], "/")
				fileNameLen := len( primaryDoc[len(primaryDoc)-1])
				
				if(fileNameLen >= 3 && primaryDoc[len(primaryDoc)-1][fileNameLen-3:fileNameLen] == "xml") {
						// build the url
					uri := "https://www.sec.gov/Archives/edgar/data/" + apiResponse.Cik + "/" + accNumUrl + "/" + primaryDoc[len(primaryDoc)-1]

					entries = append(entries, FormJsonEntry{
						Url:    uri,
						AccNum: recent.AccNums[i],
					})
				}
				
			}
		}
	}

	return entries,cik

}

type SecApiResponse struct {
	Cik                             string   `json:"cik"`
	EntityType                      string   `json:"entityType"`
	Sic                             string   `json:"sic"`
	SicDescription                  string   `json:"sicDescription"`
	IsReporter                      int      `json:"insiderTransactionForOwnerExists"`
	IsIssuer                        int      `json:"insiderTransactionForIssuerExists"`
	Name                            string   `json:"name"`
	Ein                             string   `json:"ein"`
	Tickers                         []string `json:"tickers"`
	Exchanges                       []string `json:"exchanges"`
	Description                     string   `json:"description"`
	Website                         string   `json:"website"`
	InvestorWebsite                 string   `json:"InvestorWebsite"`
	Category                        string   `json:"category"`
	FiscalYearEnd                   string   `json:"fiscalYearEnd"`
	StateOfIncorporation            string   `json:"stateOfIncorporation"`
	StateOfIncorporationDescription string   `json:"stateOfIncorporationDescription"`
	Phone                           string   `json:"phone"`
	Flags                           string   `json:"flags"`
	Filings                         filings  `json:"filings"`
}

type filings struct {
	Recent recent `json:"recent"`
	Files  []file `json:"files"`
}

type recent struct {
	AccNums     []string `json:"accessionNumber"`
	Dates       []string `json:"filingDate"`
	ReportDates []string `json:"reportDate"`
	Forms       []string `json:"form"`
	PrimaryDocs []string `json:"primaryDocument"`
}

type file struct {
	Name        string `json:"name"`
	FilingCount int    `json:"filingCount"`
	FilingFrom  string `json:"filingFrom"`
	FilingTo    string `json:"filingTo"`
}
