package shared_crawler_utils

import (
	"bytes"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"
	"unicode"
)

func FromURLLoadForm4XML(url, accNum, userAgent string) (RawForm4,error) {
	var form4 RawForm4

	// create the http request
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		panic(err)
	}

	req.Header.Set("User-Agent", userAgent)
	req.Host = "www.sec.gov"

	// create the http client
	client := &http.Client{}

	var resp *http.Response

	//Attempt the request 10 times on account of random 503s and connected parties not responding in time
	requestSuccess := false
	for i := 0; i < 3 && !requestSuccess; i++ { 
		//Only make request to SEC if one is available obeying 10/sec rule
		ConsumeSECRequest()
		resp, err = client.Do(req)
		//404 is allowed because it means the form doesn't exist
		if err == nil && (resp.StatusCode == 200 || resp.StatusCode == 404){
			requestSuccess = true
		}
	}
	
	if requestSuccess {
		// get data from res body
		data, err := io.ReadAll(resp.Body)
		if err != nil {
			panic(err)
		}

		// parse the xml
		xml.Unmarshal(data, &form4)	
		if(len(form4.Footnotes.Footnote) > 0) {
			//Trims whitespace around Footnote ids " F2 " -> "F2"
			FootnotesWhitespaceless := make([]Footnote,len(form4.Footnotes.Footnote))
			for i, footnote := range form4.Footnotes.Footnote {
				FootnotesWhitespaceless[i] = Footnote{
					XMLName: footnote.XMLName,
					Text: footnote.Text,
					FootnoteId: strings.TrimSpace(footnote.FootnoteId),
				}
			}
			form4.Footnotes.Footnote = FootnotesWhitespaceless
			decoder := xml.NewDecoder(bytes.NewReader(data))
			walkXML(xml.Name{}, decoder, form4,[]int{-1,-1,-1,-1},-1)
		}
		today := time.Now()
		form4.AccessionNumber = accNum
		form4.Url = url
		form4.DateAdded = today.Format("2006-01-02")    

		return form4,nil
	} else {
		// Form couldn't be parsed
		return form4, errors.New("Status Code: " + resp.Status + "Form: " + url)
	}
}

/*
*	Walks XML and writes the name of fields with the FootnoteId in snake case such
*	that a parent field with FootnoteId attr id="F2" replaces the second value in the names array
*	<exerciseDate>
			YYYY-MM-DD
		<FootnoteId id="F2"/>		--> names[1] = exercise_date
	</exerciseDate>
	@indexCounter holds - NDH, DH, NDT, DT indicies of the 4 transaction type FootnoteId arrays, should all start at -1
	@activeFlag indicates which transaction type we are in for purposes of associating Footnote with a transaction
			-1 = none, 0 = NDH, 1 = DH, 2 =NDT, 3 = DT (Same order as params)
*/
func walkXML(parent xml.Name, decoder *xml.Decoder, form4 RawForm4,indexCounter []int,activeFlag int) {
	for {
		token, err := decoder.Token()
		if err != nil {
			if err == io.EOF {
				break
			}
			log.Fatal(err)
		}

		switch t := token.(type) {
		case xml.StartElement:
			switch t.Name.Local {
				case "FootnoteId" : 
			
					for _, attr := range t.Attr {
						if attr.Name.Local == "id" {
							//Value is of form F#, # is the number
							attrVal := attr.Value 
							attrVal = strings.TrimSpace(attrVal) //" F2 " -> "F2"
							
							//TODO the field_referenced saved is always snakeCase because I assume all the SQL columns are named correctly :)
							//IF one is wrong we'll need a mapping between snake_case attributes -> sql column names
							i := form4.DerivativeTable.FootnoteIds
							fmt.Println(i)
							switch activeFlag {
							case -1:
								//Footnote_inst isn't in any transaction
								if form4.Footnote_inst == nil {
									form4.Footnote_inst = make(map[string][]string)
								}
								form4.Footnote_inst[attrVal] = append(form4.Footnote_inst[attrVal], camelToSnakeCase(parent.Local))
							//By contract if the activeFlag isn't -1 the associated index variable must be set
							case 0: //NDH
								if form4.NonDerivativeTable.NonDerivativeHoldings[indexCounter[0]].Footnote_inst == nil {
									form4.NonDerivativeTable.NonDerivativeHoldings[indexCounter[0]].Footnote_inst = make(map[string][]string)
								}
								form4.NonDerivativeTable.NonDerivativeHoldings[indexCounter[0]].Footnote_inst[attrVal] = append(form4.NonDerivativeTable.NonDerivativeHoldings[indexCounter[0]].Footnote_inst[attrVal],camelToSnakeCase(parent.Local))
							case 1: //DH
								if form4.DerivativeTable.DerivativeHoldings[indexCounter[1]].Footnote_inst == nil {
									form4.DerivativeTable.DerivativeHoldings[indexCounter[1]].Footnote_inst = make(map[string][]string)
								}
								form4.DerivativeTable.DerivativeHoldings[indexCounter[1]].Footnote_inst[attrVal] = append(form4.DerivativeTable.DerivativeHoldings[indexCounter[1]].Footnote_inst[attrVal],camelToSnakeCase(parent.Local))
							case 2: //NDT
								if form4.NonDerivativeTable.NonDerivativeTransactions[indexCounter[2]].Footnote_inst == nil {
									form4.NonDerivativeTable.NonDerivativeTransactions[indexCounter[2]].Footnote_inst = make(map[string][]string)
								}
								form4.NonDerivativeTable.NonDerivativeTransactions[indexCounter[2]].Footnote_inst[attrVal] = append(form4.NonDerivativeTable.NonDerivativeTransactions[indexCounter[2]].Footnote_inst[attrVal] ,camelToSnakeCase(parent.Local))
							case 3: //DT
								if form4.DerivativeTable.DerivativeTransactions[indexCounter[3]].Footnote_inst == nil {
									form4.DerivativeTable.DerivativeTransactions[indexCounter[3]].Footnote_inst = make(map[string][]string)
								}
								form4.DerivativeTable.DerivativeTransactions[indexCounter[3]].Footnote_inst[attrVal] = append(form4.DerivativeTable.DerivativeTransactions[indexCounter[3]].Footnote_inst[attrVal] ,camelToSnakeCase(parent.Local))
							}
						}
						}
				//When entering a transaction
				//indexCounter stays increment because we always want to know how many transactions we've seen
				//activeFlag is only set on the child because we are in the kind of transaciton
				case "nonDerivativeHolding":
					indexCounter[0]++
					walkXML(t.Name, decoder,form4,indexCounter,0)
				case "derivativeHolding":
					indexCounter[1]++
					walkXML(t.Name, decoder,form4,indexCounter,1)
				case "nonDerivativeTransaction":
					indexCounter[2]++
					walkXML(t.Name, decoder,form4,indexCounter,2)
				case "derivativeTransaction":
					indexCounter[3]++
					walkXML(t.Name, decoder,form4,indexCounter,3)

				default: 
					walkXML(t.Name, decoder,form4,indexCounter, activeFlag)
				}
		case xml.EndElement:
			if t.Name == parent {
				return
			}
		}
	}
}
func camelToSnakeCase(input string) string {
	var result []rune
	for i, char := range input {
		if i > 0 && unicode.IsUpper(char) {
			result = append(result, '_')
		}
		result = append(result, unicode.ToLower(char))
	}
	return string(result)
}

