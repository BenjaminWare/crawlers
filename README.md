# crawlers

A collection of crawlers used to crawl SEC data into SQL databases expects a .env file at the root with a connection_string variable set.

## Packages

1. local-crawler - expects a downloaded version of the submissions.zip folder (found at https://www.sec.gov/edgar/sec-api-documentation) has functions to crawl all forms, issuers, and stock days useful for get all historical forms. Must be run once when massive schema changes are made that can't be migrated.
2. live-crawler - Reads from the SEC form4 rss feed to get current day forms. These feed only has the 2000 most recent forms so it must run frequently to not miss forms. It will read issuers that are currently saved and only crawl issuers that aren't already in the DB for example if a company IPOs and is brand new.
3. stock_day_crawler - Gets todays stock data information for all companies, should run daily. Uses the EOD historical data endpoint.
4. issuer_update_crawler - Updates issuer and ticker tables as information about issuers changes. Doesn't need to run to often as changes to companies are more rare (once a week should work).
5. shared_crawler_utils - Has helper function used across multiple crawlers also includes sub packages for issuer and stock_day crawling specifically. (the stock_day_crawler is only responsible for todays stock_data the other crawlers get historical stock_data)
   Setting up Live-Crawler:
   To build the live crawler docker image from the root directory (The file is named Dockerfile.live_crawler anticipating the creation of other dockerfiles):
   docker build -f Dockerfile.live_crawler --platform linux/amd64 -t docker-image:test .

Create an ECR_REPO in AWS then run these docker commands to populate it:
docker tag docker-image:test <ECR_REPO_URI>:latest
docker push <ECR_REPO_URI>:latest

Create a lambda function from the container image you created in ECR

### Helpful Hints About AWS

When populating an aurora db from a sqldump:

1. Make sure the DB is publicly accessbible in AWS RDS
2. Make sure your dump file has GTID set to OFF (This is in the advanced settings in Workbench data export)
3. Make sure to include create schema in your dump file if trying to load from one dump file

### Helpful Notes

1. EOD API doesn't accept 0000-00-00 as a date use 0001-01-01 to get first day
