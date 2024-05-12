package main

import (
	"database/sql"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"

	stock_day_crawler "insiderviz.com/crawlers/stock_day_crawler"
	utils "insiderviz.com/crawlers/utils"
)

var (
	conn *sql.DB
)

func main() {
	os.Stdout = nil
	lambda.Start(handler)
}

func handler(event events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	conn = utils.CreateMySQLConnection(os.Getenv("CONNECTION_STRING"))
	stock_day_crawler.CrawlTodaysStockDay(conn, 10)
	// Return http response based on if crawl succeeded
	response := events.APIGatewayProxyResponse{
		StatusCode: 200,
		Body:       "Success",
	}
	return response, nil
}
