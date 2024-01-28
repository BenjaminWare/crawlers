package shared_crawler_functions

import (
	"context"
	"database/sql"
	"encoding/xml"
	"fmt"
	"sort"
	"strings"
)

func setIssuer(tx *sql.Tx, cik, name string) error {
	issuerSql := `
	insert into issuer (cik, name, sic, sic_description, ein, state_of_incorporation, phone, sector, industry)
	values (?, ?, '', '', '', '', '', '', '')
	ON DUPLICATE KEY UPDATE cik=cik
	`
	_, err := tx.Exec(issuerSql, cik, name)

	return err
}

func setReporter(tx *sql.Tx, cik, name string) error {
	reporterSql := `
	insert into reporter (cik, name)
	values (?, ?)	
	ON DUPLICATE KEY UPDATE cik=cik
	`
	_, err := tx.Exec(reporterSql, cik, name)

	return err
}

func saveNonDerivTransaction(tx *sql.Tx, acc_num string, data NonDerivativeTransaction) int {
	transactionSql := `
	insert into non_derivative_transaction (
		acc_num,
		security_title,
		transaction_date,
		transaction_form_type,
		transaction_code,
		equity_swap_involved,
		transaction_shares,
		transaction_price_per_share,
		transaction_acquired_disposed_code,
		post_transaction_amounts_shares,
		ownership_nature,
		is_holding
	) values (
		?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?
	);`

	result, err := tx.Exec(transactionSql, acc_num, data.SecurityTitle.Value, substringIfLongEnough(data.TransactionDate.Value, 0, 10),
		data.TransactionCoding.TransactionFormType,
		data.TransactionCoding.TransactionCode,
		data.TransactionCoding.EquitySwapInvolved,
		data.TransactionAmounts.TransactionShares.Value,
		data.TransactionAmounts.TransactionPricePerShare.Value,
		data.TransactionAmounts.TransactionAcquiredDisposedCode.Value,
		data.PostTransactionAmounts.SharesOwnedFollowingTransaction.Value,
		data.OwnershipNature.NatureOfOwnership.Value,
		false,
	)
	if err != nil {
		fmt.Println("Error saving transaction1:", err)
	}
	id, err := result.LastInsertId()
	if err != nil {
		fmt.Println("Error saving transaction1:", err)
	}
	return int(id)
}
func saveNonDerivHolding(tx *sql.Tx, acc_num string, data NonDerivativeHolding) int {
	transactionSql := `
	insert into non_derivative_transaction (
		acc_num,
		security_title,
		transaction_date,
		transaction_form_type,
		transaction_code,
		equity_swap_involved,
		transaction_shares,
		transaction_price_per_share,
		transaction_acquired_disposed_code,
		post_transaction_amounts_shares,
		ownership_nature,
		is_holding
	) values (
		?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?
	)`

	result, err := tx.Exec(transactionSql, acc_num, data.SecurityTitle.Value,
		nil,
		nil,
		nil,
		nil,
		nil,
		nil,
		nil,
		data.PostTransactionAmounts.SharesOwnedFollowingTransaction.Value,
		data.OwnershipNature.NatureOfOwnership.Value,
		true,
	)
	if err != nil {
		fmt.Println("Error saving transaction2:", err)
	}
	id, err := result.LastInsertId()
	if err != nil {
		fmt.Println("Error saving transaction2:", err)
	}
	return int(id)
}

func saveDerivTransaction(tx *sql.Tx, acc_num string, data DerivativeTransaction) int {
	transactionSql := `
	insert into derivative_transaction (
		acc_num,
		security_title,
		conversion_or_exercise_price,
		transaction_date,
		transaction_form_type,
		transaction_code,
		equity_swap_involved,
		transaction_shares,
		transaction_price_per_share,
		transaction_acquired_disposed_code,
		exercise_date,
		expiration_date,
		underlying_security_title,
		underlying_security_shares,
		post_transaction_amounts_shares,
		ownership_nature,
		is_holding
	)
	values (
		?, ?, ?, ?, ?, ?, ?, ?, ?, ?,
		?, ?, ?, ?, ?, ?, ?
	)`

	result, err := tx.Exec(transactionSql, acc_num, data.SecurityTitle.Value, data.ConversionOrExercisePrice.Value,
		substringIfLongEnough(data.TransactionDate.Value, 0, 10), data.TransactionCoding.TransactionFormType,
		data.TransactionCoding.TransactionCode,
		data.TransactionCoding.EquitySwapInvolved,
		data.TransactionAmounts.TransactionShares.Value,
		data.TransactionAmounts.TransactionPricePerShare.Value,
		data.TransactionAmounts.TransactionAcquiredDisposedCode.Value,
		substringIfLongEnough(data.ExerciseDate.Value, 0, 10), substringIfLongEnough(data.ExpirationDate.Value, 0, 10),
		data.UnderlyingSecurity.UnderlyingSecurityTitle.Value,
		data.UnderlyingSecurity.UnderlyingSecurityShares.Value,
		data.PostTransactionAmounts.SharesOwnedFollowingTransaction.Value,
		data.OwnershipNature.NatureOfOwnership.Value,
		false,
	)
	if err != nil {
		fmt.Println("Error saving transaction3:", err)
	}
	id, err := result.LastInsertId()
	if err != nil {
		fmt.Println("Error saving transaction3:", err)
	}
	return int(id)
}

func substringIfLongEnough(str string, start int, end int) string {
	if len(str) >= end {
		return str[start:end]
	} else {
		return str
	}
}

func saveDerivHolding(tx *sql.Tx, acc_num string, data DerivativeHolding) int {
	transactionSql := `
	insert into derivative_transaction (
		acc_num,
		security_title,
		conversion_or_exercise_price,
		transaction_date,
		transaction_form_type,
		transaction_code,
		equity_swap_involved,
		transaction_shares,
		transaction_price_per_share,
		transaction_acquired_disposed_code,
		exercise_date,
		expiration_date,
		underlying_security_title,
		underlying_security_shares,
		post_transaction_amounts_shares,
		ownership_nature,
		is_holding
	)
	values (
		?, ?, ?, ?, ?, ?, ?, ?, ?, ?,
		?, ?, ?, ?, ?, ?, ?
	)`
	result, err := tx.Exec(transactionSql, acc_num, data.SecurityTitle.Value, data.ConversionOrExercisePrice.Value,
		nil,
		nil,
		nil,
		nil,
		nil,
		nil,
		nil,
		substringIfLongEnough(data.ExerciseDate.Value, 0, 10), substringIfLongEnough(data.ExpirationDate.Value, 0, 10),
		data.UnderlyingSecurity.UnderlyingSecurityTitle.Value,
		data.UnderlyingSecurity.UnderlyingSecurityShares.Value,
		data.PostTransactionAmounts.SharesOwnedFollowingTransaction.Value,
		data.OwnershipNature.NatureOfOwnership.Value,
		true,
	)
	if err != nil {
		fmt.Println("Error saving transaction4:", err)
	}
	id, err := result.LastInsertId()
	if err != nil {
		fmt.Println("Error saving transaction4:", err)
	}
	return int(id)

}
func saveFootnote(tx *sql.Tx, acc_num string, text string) int {
	footnoteSql := `
	insert into footnote (acc_num,text)
	values (?, ?)	
	`
	result, err := tx.Exec(footnoteSql, acc_num, text)
	if err != nil {
		fmt.Println("Error saving footnote:", err)
	}
	id, err := result.LastInsertId()
	if err != nil {
		fmt.Println("Error saving footnote:", err)
	}
	return int(id)
}

/*
Saves a footnote instance to sql as defined by a footnote id being placed on a field, in xml this is </footnoteId id="F3"> in html it looks like (3) in a field
@ndt_id, @dt_id a footnote_inst can be part of at most one transaction so at most one of these isn't nil by contract
*/
func saveFootnoteInst(tx *sql.Tx, acc_num string, f_id int, dt_id interface{}, ndt_id interface{}, field_referenced string) {
	FootnoteInstSql := `
	insert into footnote_inst (acc_num, footnote_id,dt_id,ndt_id,field_referenced)
	values (?, ?,?,?,?)	
	`
	_, err := tx.Exec(FootnoteInstSql, acc_num, f_id, dt_id, ndt_id, field_referenced)
	if err != nil {
		fmt.Println("Error FootnoteInst:", err)
	}
}

// Puts one RawForm4 in the db this includes entering any associated transactions, footnotes or footnote_inst
// Returns true when the form was already in the db and false otherwise that includes successfully entry or error
func SaveForm(Conn *sql.DB,form RawForm4) bool{
		form_duplicate := false
		// start the transaction
		tx, err := Conn.BeginTx(context.TODO(),nil)
		if err != nil {
			fmt.Println("Error starting transaction:", err)
			panic(err)
		}
		sqlStr := `SELECT acc_num FROM form WHERE acc_num=?`
		row := tx.QueryRow(sqlStr, form.AccessionNumber)
		acc_num_result := "not needed..."
		err = row.Scan(&acc_num_result)
		//Some forms have no content on them and the length of form.ReportingOwners is a good test
		if err != nil && len(form.ReportingOwners) > 0{
			// save issuer and reporter
			err = setIssuer(tx, form.Issuer.IssuerCIK, form.Issuer.IssuerName)
			if err != nil {
				fmt.Println("Error saving issuer:", err)
				panic(err)
			}
			
			err = setReporter(tx, form.ReportingOwners[0].ReportingOwnerId.RptOwnerCik, form.ReportingOwners[0].ReportingOwnerId.RptOwnerName)
			if err != nil {
				fmt.Println("Error saving reporter:", err)
				panic(err)
			}
		
			var net_shares float32 = 0.0
			var net_total float32 = 0.0
			transaction_codes := ""
			for _, trans := range form.NonDerivativeTable.NonDerivativeTransactions {
				if !strings.Contains(transaction_codes,trans.TransactionCoding.TransactionCode) {
					transaction_codes += trans.TransactionCoding.TransactionCode
				}
				
				if trans.TransactionAmounts.TransactionAcquiredDisposedCode.Value == "D" {
					net_shares -= trans.TransactionAmounts.TransactionShares.Value
					net_total -= trans.TransactionAmounts.TransactionPricePerShare.Value * trans.TransactionAmounts.TransactionShares.Value
				} else {
					net_shares += trans.TransactionAmounts.TransactionShares.Value
					net_total += trans.TransactionAmounts.TransactionPricePerShare.Value * trans.TransactionAmounts.TransactionShares.Value
				}
			}
			transaction_codes_slice :=  (strings.Split(transaction_codes, ""))
			sort.Strings(transaction_codes_slice)
			transaction_codes = strings.Join(transaction_codes_slice,"")
			// save the base form
			formSql := `
			insert into form (acc_num, period_of_report, rpt_is_director, rpt_is_officer, rpt_is_ten_percent_owner, rpt_is_other, rpt_officer_title, rpt_other_text, issuer_cik, reporter_cik, xml_url, pdf_url,net_shares,net_total,transaction_codes)
			values (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
			ON DUPLICATE KEY UPDATE acc_num=acc_num
			`
			_, err = tx.Exec(formSql, form.AccessionNumber, form.PeriodOfReport[:10], form.ReportingOwners[0].ReportingOwnerRelationship.IsDirector,
				form.ReportingOwners[0].ReportingOwnerRelationship.IsOfficer, form.ReportingOwners[0].ReportingOwnerRelationship.IsTenPercentOwner,
				form.ReportingOwners[0].ReportingOwnerRelationship.IsOther, form.ReportingOwners[0].ReportingOwnerRelationship.OfficerTitle,
				form.ReportingOwners[0].ReportingOwnerRelationship.OtherText, form.Issuer.IssuerCIK, form.ReportingOwners[0].ReportingOwnerId.RptOwnerCik, form.Url, "",net_shares,net_total,transaction_codes)
			if err != nil {
				fmt.Println("Error saving form:", err)
			}
				

		//save footnotes
		//Hold onto the id of footnotes inserted, in case they are in a transaction, so we can put the id in the join table
		footnoteSQLIds := map[string]int{} // F# to SQLId
		for _, footnote := range form.Footnotes.Footnote {

			id := saveFootnote(tx, form.AccessionNumber, footnote.Text)
			footnoteSQLIds[footnote.FootnoteId] = id
		}
		//saves footnote instances that don't belong to a particular transaction
		for footnote_id, fields_referenced := range form.Footnote_inst {
			for _, field_referenced := range fields_referenced {
				saveFootnoteInst(tx, form.AccessionNumber, footnoteSQLIds[footnote_id], nil, nil, field_referenced)
			}
		}
		// save the derivative transactions
		for _, transaction := range form.DerivativeTable.DerivativeTransactions {
			id := saveDerivTransaction(tx, form.AccessionNumber, transaction)

			for footnote_id, fields_referenced := range transaction.Footnote_inst {
				for _, field_referenced := range fields_referenced {
					saveFootnoteInst(tx, form.AccessionNumber, footnoteSQLIds[footnote_id], id, nil, field_referenced)
				}

			}
		}

		// save the derivative holdings, same table as derivative transactions with a flag
		for _, holding := range form.DerivativeTable.DerivativeHoldings {
			id := saveDerivHolding(tx, form.AccessionNumber, holding)
			for footnote_id, fields_referenced := range holding.Footnote_inst {
				for _, field_referenced := range fields_referenced {
					saveFootnoteInst(tx, form.AccessionNumber, footnoteSQLIds[footnote_id], id, nil, field_referenced)
				}
			}
		}

		// save the non derivative transactions
		for _, transaction := range form.NonDerivativeTable.NonDerivativeTransactions {
			id := saveNonDerivTransaction(tx, form.AccessionNumber, transaction)
			for footnote_id, fields_referenced := range transaction.Footnote_inst {
				for _, field_referenced := range fields_referenced {
					saveFootnoteInst(tx, form.AccessionNumber, footnoteSQLIds[footnote_id], nil, id, field_referenced)
				}
			}
		}
		// save the non-derivative holdings, same table as derivative transactions with a flag
		for _, holding := range form.NonDerivativeTable.NonDerivativeHoldings {
			id := saveNonDerivHolding(tx, form.AccessionNumber, holding)
			for footnote_id, fields_referenced := range holding.Footnote_inst {
				for _, field_referenced := range fields_referenced {
					saveFootnoteInst(tx, form.AccessionNumber, footnoteSQLIds[footnote_id], nil, id, field_referenced)
				}
			}

		}

		} else {
		// Form was already in the DB
		form_duplicate = true
		}

	err = tx.Commit()
	if err != nil {
		panic(err)
	}
	return form_duplicate

}

type RawForm4 struct {
	XMLName                  xml.Name           `xml:"ownershipDocument" json:"-"`
	DateAdded                string             `xml:"-" json:"DateAdded"`
	SchemaVersion            string             `xml:"schemaVersion"`
	DocumentType             string             `xml:"documentType"`
	PeriodOfReport           string             `xml:"periodOfReport"`
	DateOfOriginalSubmission string             `xml:"dateOfOriginalSubmission"`
	NoSecuritiesOwned        bool               `xml:"noSecuritiesOwned"`
	AccessionNumber          string             `xml:"accessionNumber"`
	Url                      string             `xml:"url"`
	Issuer                   Issuer             `xml:"issuer"`
	ReportingOwners          []ReportingOwner   `xml:"reportingOwner" json:"ReportingOwners"`
	NonDerivativeTable       NonDerivativeTable `xml:"nonDerivativeTable"`
	DerivativeTable          DerivativeTable    `xml:"derivativeTable"`
	OwnerSignatures          OwnerSignature     `xml:"ownerSignature"`
	Footnotes                Footnotes          `xml:"footnotes"`
	//Only has instances not in a transaction
	Footnote_inst            map[string][]string 
}

type Footnotes struct {
	XMLName   xml.Name `xml:"footnotes" json:"-"`
	Footnote []Footnote `xml:"footnote"`

}
type Footnote struct {
	XMLName   xml.Name `xml:"footnote" json:"-"`
	Text 	   string `xml:",chardata"`
	FootnoteId string `xml:"id,attr"`
}

type Issuer struct {
	XMLName             xml.Name `xml:"issuer" json:"-"`
	IssuerCIK           string   `xml:"issuerCik"`
	IssuerName          string   `xml:"issuerName"`
	IssuerTradingSymbol string   `xml:"issuerTradingSymbol"`
}

type ReportingOwner struct {
	XMLName                    xml.Name                   `xml:"reportingOwner"`
	ReportingOwnerId           ReportingOwnerId           `xml:"reportingOwnerId"`
	ReportingOwnerAddress      ReportingOwnerAddress      `xml:"reportingOwnerAddress" `
	ReportingOwnerRelationship ReportingOwnerRelationship `xml:"reportingOwnerRelationship" `
}

type ReportingOwnerId struct {
	XMLName      xml.Name `xml:"reportingOwnerId"`
	RptOwnerCik  string   `xml:"rptOwnerCik"`
	RptOwnerCcc  string   `xml:"rptOwnerCcc"`
	RptOwnerName string   `xml:"rptOwnerName"`
}

type ReportingOwnerAddress struct {
	XMLName                  xml.Name `xml:"reportingOwnerAddress"`
	RptOwnerStreet1          string   `xml:"rptOwnerStreet1"`
	RptOwnerStreet2          string   `xml:"rptOwnerStreet2"`
	RptOwnerCity             string   `xml:"rptOwnerCity"`
	RptOwnerState            string   `xml:"rptOwnerState"`
	RptOwnerZipCode          string   `xml:"rptOwnerZipCode"`
	RptOwnerStateDescription string   `xml:"rptOwnerStateDescription"`
}

type ReportingOwnerRelationship struct {
	XMLName           xml.Name `xml:"reportingOwnerRelationship"`
	IsDirector        bool     `xml:"isDirector"`
	IsOfficer         bool     `xml:"isOfficer"`
	IsTenPercentOwner bool     `xml:"isTenPercentOwner"`
	IsOther           bool     `xml:"isOther"`
	OfficerTitle      string   `xml:"officerTitle"`
	OtherText         string   `xml:"otherText"`
}

type NonDerivativeTable struct {
	XMLName                   xml.Name                   `xml:"nonDerivativeTable"`
	NonDerivativeTransactions []NonDerivativeTransaction `xml:"nonDerivativeTransaction" `
	NonDerivativeHoldings     []NonDerivativeHolding     `xml:"nonDerivativeHolding"`
}

type NonDerivativeTransaction struct {
	XMLName                xml.Name               `xml:"nonDerivativeTransaction"`
	SecurityTitle          SecurityTitle          `xml:"securityTitle" `
	TransactionDate        TransactionDate        `xml:"transactionDate" `
	DeemedExecutionDate    DeemedExecutionDate    `xml:"deemedExecutionDate" `
	TransactionCoding      TransactionCoding      `xml:"transactionCoding" `
	TransactionTimeliness  TransactionTimeliness  `xml:"transactionTimeliness" `
	TransactionAmounts     TransactionAmounts     `xml:"transactionAmounts" `
	PostTransactionAmounts PostTransactionAmounts `xml:"postTransactionAmounts" `
	OwnershipNature        OwnershipNature        `xml:"ownershipNature" `
	Footnote_inst 		 map[string][]string 	//Maps F# to field_name
}

type NonDerivativeHolding struct {
	XMLName                xml.Name               `xml:"nonDerivativeHolding"`
	SecurityTitle          SecurityTitle          `xml:"securityTitle" `
	PostTransactionAmounts PostTransactionAmounts `xml:"postTransactionAmounts" `
	OwnershipNature        OwnershipNature        `xml:"ownershipNature" `
	Footnote_inst 		 map[string][]string 	//Maps F# to field_name
}

type DerivativeTable struct {
	XMLName                xml.Name                `xml:"derivativeTable"`
	DerivativeTransactions []DerivativeTransaction `xml:"derivativeTransaction" `
	DerivativeHoldings     []DerivativeHolding     `xml:"derivativeHolding" `
	FootnoteIds				[]int
}

type DerivativeTransaction struct {
	XMLName                   xml.Name `xml:"derivativeTransaction"`
	DerivativeTableID         uint
	SecurityTitle             SecurityTitle             `xml:"securityTitle"`
	ConversionOrExercisePrice ConversionOrExercisePrice `xml:"conversionOrExercisePrice" `
	TransactionDate           TransactionDate           `xml:"transactionDate"`
	DeemedExecutionDate       DeemedExecutionDate       `xml:"deemedExecutionDate" `
	TransactionCoding         TransactionCoding         `xml:"transactionCoding"`
	TransactionTimeliness     TransactionTimeliness     `xml:"transactionTimeliness"`
	TransactionAmounts        TransactionAmounts        `xml:"transactionAmounts"`
	ExerciseDate              ExerciseDate              `xml:"exerciseDate"`
	ExpirationDate            ExpirationDate            `xml:"expirationDate" `
	UnderlyingSecurity        UnderlyingSecurity        `xml:"underlyingSecurity"`
	PostTransactionAmounts    PostTransactionAmounts    `xml:"postTransactionAmounts" `
	OwnershipNature           OwnershipNature           `xml:"ownershipNature"`
	Footnote_inst 		 	  map[string][]string 	//Maps F# to field_name
}

type DerivativeHolding struct {
	XMLName                   xml.Name `xml:"derivativeHolding"`
	DerivativeTableID         uint
	SecurityTitle             SecurityTitle             `xml:"securityTitle" `
	ConversionOrExercisePrice ConversionOrExercisePrice `xml:"conversionOrExercisePrice" `
	ExerciseDate              ExerciseDate              `xml:"exerciseDate" `
	ExpirationDate            ExpirationDate            `xml:"expirationDate" `
	UnderlyingSecurity        UnderlyingSecurity        `xml:"underlyingSecurity" `
	PostTransactionAmounts    PostTransactionAmounts    `xml:"postTransactionAmounts" `
	OwnershipNature           OwnershipNature           `xml:"ownershipNature" `
	Footnote_inst 		 	  map[string][]string 	//Maps F# to field_name
}

type UnderlyingSecurity struct {
	XMLName                  xml.Name                 `xml:"underlyingSecurity"`
	UnderlyingSecurityTitle  UnderlyingSecurityTitle  `xml:"underlyingSecurityTitle" `
	UnderlyingSecurityShares UnderlyingSecurityShares `xml:"underlyingSecurityShares" `
	UnderlyingSecurityValue  UnderlyingSecurityValue  `xml:"underlyingSecurityValue" `
}

type UnderlyingSecurityShares struct {
	XMLName xml.Name `xml:"underlyingSecurityShares"`
	Value   float32   `xml:"value"`
}

type UnderlyingSecurityValue struct {
	XMLName xml.Name `xml:"underlyingSecurityValue"`
	Value   string   `xml:"value"`
}

type UnderlyingSecurityTitle struct {
	XMLName xml.Name `xml:"underlyingSecurityTitle"`
	Value   string   `xml:"value"`
}

type ExpirationDate struct {
	XMLName xml.Name `xml:"expirationDate"`
	Value   string   `xml:"value"`
}

type ExerciseDate struct {
	XMLName xml.Name `xml:"exerciseDate"`
	Value   string   `xml:"value"`
}

type ConversionOrExercisePrice struct {
	XMLName xml.Name `xml:"conversionOrExercisePrice"`
	Value   float32  `xml:"value"`
}

type SecurityTitle struct {
	XMLName xml.Name `xml:"securityTitle"`
	Value   string   `xml:"value"`
}

type TransactionDate struct {
	XMLName xml.Name `xml:"transactionDate"`
	Value   string   `xml:"value"`
}

type DeemedExecutionDate struct {
	XMLName xml.Name `xml:"deemedExecutionDate"`
	Value   string   `xml:"value"`
}

type TransactionCoding struct {
	XMLName             xml.Name `xml:"transactionCoding"`
	TransactionFormType string   `xml:"transactionFormType" `
	TransactionCode     string   `xml:"transactionCode"`
	EquitySwapInvolved  bool     `xml:"equitySwapInvolved" `
}

type TransactionTimeliness struct {
	XMLName xml.Name `xml:"transactionTimeliness"`
	Value   string   `xml:"value"`
}

type TransactionAmounts struct {
	XMLName                         xml.Name                        `xml:"transactionAmounts"`
	TransactionShares               TransactionShares               `xml:"transactionShares" `
	TransactionTotalValue           TransactionTotalValue           `xml:"transactionTotalValue" `
	TransactionPricePerShare        TransactionPricePerShare        `xml:"transactionPricePerShare" `
	TransactionAcquiredDisposedCode TransactionAcquiredDisposedCode `xml:"transactionAcquiredDisposedCode" `
}

type TransactionShares struct {
	XMLName xml.Name `xml:"transactionShares"`
	Value   float32  `xml:"value"`
}

type TransactionTotalValue struct {
	XMLName xml.Name `xml:"transactionTotalValue"`
	Value   string   `xml:"value"`
}

type TransactionPricePerShare struct {
	XMLName xml.Name `xml:"transactionPricePerShare"`
	Value   float32  `xml:"value"`
}

type TransactionAcquiredDisposedCode struct {
	XMLName xml.Name `xml:"transactionAcquiredDisposedCode"`
	Value   string   `xml:"value"`
}

type PostTransactionAmounts struct {
	XMLName                         xml.Name                        `xml:"postTransactionAmounts"`
	SharesOwnedFollowingTransaction SharesOwnedFollowingTransaction `xml:"sharesOwnedFollowingTransaction" `
	ValueOwnedFollowingTransaction  ValueOwnedFollowingTransaction  `xml:"valueOwnedFollowingTransaction" `
}

type SharesOwnedFollowingTransaction struct {
	XMLName xml.Name `xml:"sharesOwnedFollowingTransaction"`
	Value   float32  `xml:"value"`
}

type ValueOwnedFollowingTransaction struct {
	XMLName xml.Name `xml:"valueOwnedFollowingTransaction"`
	Value   string   `xml:"value"`
}

type OwnershipNature struct {
	XMLName                   xml.Name                  `xml:"ownershipNature"`
	DirectOrIndirectOwnership DirectOrIndirectOwnership `xml:"directOrIndirectOwnership" `
	NatureOfOwnership         NatureOfOwnership         `xml:"natureOfOwnership" `
}

type DirectOrIndirectOwnership struct {
	XMLName xml.Name `xml:"directOrIndirectOwnership"`
	Value   string   `xml:"value"`
}

type NatureOfOwnership struct {
	XMLName xml.Name `xml:"natureOfOwnership"`
	Value   string   `xml:"value"`
}

type OwnerSignature struct {
	XMLName       xml.Name `xml:"ownerSignature"`
	SignatureName string   `xml:"signatureName"`
	SignatureDate string   `xml:"signatureDate"`
}
