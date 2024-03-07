package stock_data

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

func parseStockData(client *http.Client,ticker string,startDate string) ([]stockDay,error){
	uri := fmt.Sprintf("https://eodhistoricaldata.com/api/eod/%s?fmt=json&from=%s&api_token=%s", ticker, startDate, os.Getenv("EOD_TOKEN"))
	stockData := make([]stockDay,0)
	// create the http request
	req, err := http.NewRequest("GET", uri, nil)
	if err != nil {
		return stockData, err
	}
	resp, err := client.Do(req)
	if err != nil || resp.StatusCode != 200 {
		return stockData, err
	}

	// convert response body to byte array
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return stockData, err
	}
	// parse the response
	var apiRes []eodDayEntry
	err = json.Unmarshal(bodyBytes, &apiRes)
	if err != nil {
		return stockData, err
	}

	// convert the response to the correct format
	for _, entry := range apiRes {
		stockData = append(stockData, stockDay{
			Ticker: ticker,
			Date:   entry.Date,
			Close:  entry.AdjustedClose,
			Volume: entry.Volume,
		})
	}


	return stockData,nil
}

type stockDay struct {
	Ticker		  string 
	Date          string
	Close         float64 
	Volume		  int64
}
type eodDayEntry struct {
	Date          string  `json:"date"`
	Open          float64 `json:"open"`
	High          float64 `json:"high"`
	Low           float64 `json:"low"`
	Close         float64 `json:"close"`
	AdjustedClose float64 `json:"adjusted_close"`
	Volume        int64   `json:"volume"`
}
