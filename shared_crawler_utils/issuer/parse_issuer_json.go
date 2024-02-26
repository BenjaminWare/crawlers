package issuer

import (
	"encoding/csv"
	"encoding/json"
	"io"
	"net/http"
	"os"

	utils "insiderviz.com/crawlers/shared_crawler_utils"
)


func parseIssuerJSON(cik string) utils.Issuer {
	var issuer utils.Issuer
	client := &http.Client{}

	url := "https://data.sec.gov/submissions/CIK"+cik+".json"
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		panic(err)

	}

	req.Header.Set("User-Agent", "Davis bmdavis419@gmail.com")
	req.Host = "data.sec.gov"

	// Only make the request if the global guard allows it
	utils.ConsumeSECRequest()
	// Make the request using http.Client
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	if resp.StatusCode == 429 {
		print(url)
		panic("429 on rss")
	}
	// // Parse the response body using gofeed.Parser
	data,err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	json.Unmarshal(data,&issuer)
	
	issuer.Sector,issuer.Industry,err = getSectorAndIndustry(issuer.Tickers)
	if err != nil {
		panic(err)
	}
	return issuer
}

// Private helper functions

func getSectorAndIndustry(tickers []string) (string, string, error) {
	// load the csv file
	file, err := os.Open("data/Sectors.csv")
	if err != nil {
		return "", "", err
	}
	defer file.Close()

	// read in the csv
	reader := csv.NewReader(file)
	data, err := reader.ReadAll()
	if err != nil {
		panic(err)
	}

	m := createSectorMap(data)

	// search for each ticker within the entries
	sector, industry := "", ""
	for _, ticker := range tickers {
		// search the map for the ticker
		if val, ok := m[ticker]; ok {
			sector, industry = val[0], val[1]
			break
		}
	}

	return sector, industry, nil
}

func createSectorMap(data [][]string) map[string][]string {
	m := make(map[string][]string)

	// iterate over the data
	for i := 0; i < len(data); i++ {
		// each line has TICKER,SECTOR,INDUSTRY
		curLine := data[i]

		m[curLine[0]] = []string{curLine[1], curLine[2]}
	}

	return m
}
