package local_crawler

import (
	"database/sql"

	. "insiderviz.com/crawlers/shared_crawler_utils"
	. "insiderviz.com/crawlers/shared_crawler_utils/issuer"
)

/*
	Runs the local crawler over a downloaded submissions folder into the provided database
	@submissions_folder - folder in which all the CIK files are, these are expected to be from https://www.sec.gov/edgar/sec-api-documentation, submissions.zip folder
	@start/end - date to start/end at as a string in YYYY-MM-DD
	@offset - which file to start at in the submissions folder, useful if the crawler crashes at any point to restart from somewhere
	@stride - how many files to jump over, if running the crawler on multiple machines should be set to how many machines are running and offset should be adjusted
	@conn - mySQL database pointer to put forms into

*/
func RunIssuerCrawl(submissions_folder string, conn *sql.DB) {
	num_threads := 20 //Specifies how many thread we should try and use, each thread handles one form

	fileNames := GetFilenamesInDirectory(submissions_folder)

	// Crawls Issuers, this is done before any of the forms to respect foreign keys and because the SEC 10 req/sec seems to trigger at less than 10 req/sec when using more than one endpoint
	////////////////////////////////
	ciks :=parseSubmissionsFolderToIssuerCiks(submissions_folder,fileNames)
	CrawlIssuersByCIK(conn,ciks,false, "", num_threads)
	/////////////////////////////////



}
