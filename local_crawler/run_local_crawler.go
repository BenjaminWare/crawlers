package local_crawler

import (
	"database/sql"
	"fmt"
	"sync"
	"time"

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
func RunLocalCrawler(submissions_folder string, start string,end string,offset int, stride int, conn *sql.DB) {
	num_threads := 20 //Specifies how many thread we should try and use, each thread handles one form
	// Ensures only num_threads are created, by using a channel of empty structs
	thread_guard := make(chan struct{} , num_threads)
	var wg sync.WaitGroup
	fileNames := GetFilenamesInDirectory(submissions_folder)
	


	//Loop through all the files, which are all the forms filed for companies and insiders, we only care about companies
	var formsCompleted int32 = 0
	currentFiles  := map[int]struct{}{} //set of integers which are the file numbers that are running, if the program fails you can safely restart it from the lowest number in this set
	var currentFilesMutex sync.Mutex
	start_time := time.Now().UnixMilli();
	
	for i :=offset; i <len(fileNames);i+=stride {
		fileName := fileNames[i]
		//forms has all the acc_nums and urls needed to get all the forms for a specific issuer
		forms,cik := parseSubmissionsFileJSON(submissions_folder+"/"+fileName, start,end)

		//This is an issuer with atleast one form
		if len(forms) > 0 {
			// Crawls the issuer data sequentially so that if the program crashes all issuers will still be present
			if cik != "" {
				 CrawlIssuersByCIK(conn,[]string{cik},1)
			}

			thread_guard <- struct{}{} // would block if guard channel is already filled
			wg.Add(1)
			currentFilesMutex.Lock()
			currentFiles[i] = struct{}{}
			currentFilesMutex.Unlock()
			go func(forms []FormJsonEntry,i int,currentFiles map[int]struct{},currentFilesMutex *sync.Mutex) {
				
				//Handles all the forms for this issuer
				issuerWorker(forms,conn,&formsCompleted,int64(start_time),"Ben maple6leaf@gmail.com")

				//Handles the goroutine ending
				<-thread_guard
				

				//Write the shared map to remove the file that just finished, needs lock to be thread safe
				(*currentFilesMutex).Lock()
				delete(currentFiles,i)

				fmt.Println("---------------------------------------------")
				fmt.Printf("Finished %d Files Currently Working: ",i)
				smallestFile := 2147483647
				for k,_ := range currentFiles {
					fmt.Printf("%d ",k)
					if k < smallestFile {
						smallestFile = k
					}
				}
				(*currentFilesMutex).Unlock()
				fmt.Printf("Smallest File(Where you should restart): %d\n",smallestFile)
				fmt.Println("---------------------------------------------")
				defer wg.Done()
			}(forms,i,currentFiles,&currentFilesMutex)
		}

	}
	wg.Wait()



}
