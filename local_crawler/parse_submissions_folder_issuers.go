package local_crawler

import (
	"encoding/json"
	"os"
	"strings"
)

// Parses the submissions folder for just ciks so that issuers can be crawled
func parseSubmissionsFolderToIssuerCiks(folder string, fileNames []string) ([]string) {
	ciks := make([]string, 0)
	for i,fileName := range fileNames {
		if i < 100000 {
			continue
		}
		jsonFile, err := os.Open(folder + "/" + fileName)
		if err != nil {
			panic(err)
		}
		defer jsonFile.Close()

		var apiResponse SecApiResponse
		decoder := json.NewDecoder(jsonFile)
		err = decoder.Decode(&apiResponse)
		if err == nil && apiResponse.IsIssuer == 1 {
			cik := apiResponse.Cik
			cik = strings.Repeat("0",10 - len(cik)) + cik
			ciks = append(ciks, cik)
		}
	}
	return ciks

}
