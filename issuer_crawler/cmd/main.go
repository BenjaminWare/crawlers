package main

import (
	"database/sql"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"

	. "insiderviz.com/crawlers/issuer_crawler"
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
	CrawlIssuers(conn, 10)
	// LiveCrawl(conn)
	// Return http response based on if crawl succeeded
	response := events.APIGatewayProxyResponse{
		StatusCode: 200,
		Body:       "SUCCESS: Recrawled issuers",
	}
	return response, nil
}
