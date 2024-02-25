package shared_crawler_utils

import (
	"context"
	"database/sql"
	"fmt"
	"math"
	"sort"
	"strings"
)


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

/*
	Puts one form in the db this includes writing the form table and all tables depending on form (transactions, footnotes etc)
	It also saves the reporter found on the form (because reporter are just cik,name if they get more complicated move this functionality)

	Returns true when the form is already present in the DB (this indicates the live crawler should stop)

*/
func SaveForm(Conn *sql.DB,form RawForm4) bool{

		// 2.8147498e+14 largest dollar amount we will hold, so it fits in sql and is consistent value, no floating point weirdness on powers of two
		var MONEY_MAX float32 =float32(math.Pow(2,48))
		form_duplicate := false
		// start the transaction
		tx, err := Conn.BeginTx(context.TODO(),nil)
		if err != nil {
			fmt.Println("Error starting transaction:", err)
			panic(err)
		}
		// Defer a rollback in case anything fails.
		defer tx.Rollback()

		// Prepare form4 data 
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

		// clamp numeric values so they fit in sql field
		if net_total > MONEY_MAX {
			net_total = MONEY_MAX
		}
		if net_shares > MONEY_MAX {
			net_shares = MONEY_MAX
		}

		transaction_codes_slice :=  (strings.Split(transaction_codes, ""))
		sort.Strings(transaction_codes_slice)
		transaction_codes = strings.Join(transaction_codes_slice,"")

		// Inserts form4, if rowsAffected on the response == 1 it succeed and everything else should be inserted
		var formResponse sql.Result
		formSql := `
		insert into form (acc_num ,period_of_report, rpt_is_director, rpt_is_officer, rpt_is_ten_percent_owner, rpt_is_other, rpt_officer_title, rpt_other_text, issuer_cik, reporter_cik, xml_url, pdf_url,net_shares,net_total,transaction_codes)
		values (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
		ON DUPLICATE KEY UPDATE acc_num=acc_num
		`
		formResponse, err = tx.Exec(formSql, form.AccessionNumber, form.PeriodOfReport[:10], form.ReportingOwners[0].ReportingOwnerRelationship.IsDirector,
			form.ReportingOwners[0].ReportingOwnerRelationship.IsOfficer, form.ReportingOwners[0].ReportingOwnerRelationship.IsTenPercentOwner,
			form.ReportingOwners[0].ReportingOwnerRelationship.IsOther, form.ReportingOwners[0].ReportingOwnerRelationship.OfficerTitle,
			form.ReportingOwners[0].ReportingOwnerRelationship.OtherText, form.IssuerCIK, form.ReportingOwners[0].ReportingOwnerId.RptOwnerCik, form.Url, "",net_shares,net_total,transaction_codes)
			
		if err != nil {
			fmt.Println(net_total)
			fmt.Println("Error saving form:", err)
		}
		rowsAffected,rowsAffectedErr := formResponse.RowsAffected()

		// Only input all the other tables if the form was successfuly inserted
		if rowsAffectedErr == nil && rowsAffected == 1{
			
			err = setReporter(tx, form.ReportingOwners[0].ReportingOwnerId.RptOwnerCik, form.ReportingOwners[0].ReportingOwnerId.RptOwnerName)
			if err != nil {
				fmt.Println("Error saving reporter:", err)
				panic(err)
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
