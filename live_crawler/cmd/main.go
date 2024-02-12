package main

import (
	"database/sql"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"

	. "insiderviz.com/crawlers/live_crawler"
	. "insiderviz.com/crawlers/shared_crawler_functions"
)

var (
	conn *sql.DB
)
func main() {
	lambda.Start(handler)  

}

func handler(event events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse,error) {
	conn = CreateMySQLConnection(os.Getenv("connection_string"))
	success := LiveCrawl(conn)

	message := "SUCCESS: Live Crawler got all forms"
	if !success {
		message = "FAILURE: Live Crawler didn't reach overlap"
	}
	// LiveCrawl(conn)
	// Return http response based on if crawl succeeded
	response := events.APIGatewayProxyResponse{
		StatusCode: 200,
		Body :message,
		// Body:       "\"Live Crawler Finished\"",
	}
	return response, nil
}