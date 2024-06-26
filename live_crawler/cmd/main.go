package main

import (
	"database/sql"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"

	. "insiderviz.com/crawlers/live_crawler"
	. "insiderviz.com/crawlers/utils"
)

var (
	conn *sql.DB
)

func main() {
	// Ensures nothing is written to console
	os.Stdout = nil
	// Runs lambda
	lambda.Start(handler)

}

func handler(event events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	conn = CreateMySQLConnection(os.Getenv("CONNECTION_STRING"))
	success := LiveCrawl(conn)

	message := "SUCCESS: Live Crawler got all forms"
	if !success {
		message = "FAILURE: Live Crawler didn't reach overlap"
	}
	// LiveCrawl(conn)
	// Return http response based on if crawl succeeded
	response := events.APIGatewayProxyResponse{
		StatusCode: 200,
		Body:       message,
		// Body:       "\"Live Crawler Finished\"",
	}
	return response, nil
}
