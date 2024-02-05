package main

import (
	"database/sql"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

var (
	conn *sql.DB
)
func main() {
	// region := os.Getenv("AWS_REGION")
	// awsSession, err := session.NewSession(&aws.Config{
	// 	Region: aws.String(region)},)

	conn = CreateMySQLConnection("root:root@tcp(127.0.0.1:3306)/insiderviz-crawler2") 
	lambda.Start(handler)
}

func handler(req events.APIGatewayProxyRequest) /*(*events.APIGatewayProxyResponse,error)*/ {
	LiveCrawl(conn)
}