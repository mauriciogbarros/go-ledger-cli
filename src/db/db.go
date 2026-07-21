package db

import (
	"database/sql"
	"time"

	"go.mod/internal/account"
	"go.mod/internal/accountType"
	"go.mod/internal/chart"
	"go.mod/internal/entry"
	"go.mod/internal/id"
	"go.mod/internal/ledger"
	_ "modernc.org/sqlite"
)

func InitializeLedger(db *sql.DB, ledger *ledger.Ledger) error {
	var stmt = `CREATE TABLE IF NOT EXISTS account_types (
		id TEXT PRIMARY KEY,
		name VARCHAR(255) NOT NULL UNIQUE,
		side INTEGER NOT NULL,
		ref_counter INTEGER NOT NULL
	);`
	_, err := db.Exec(stmt)
	if err != nil {
		return err
	}

	var count int
	stmt = `SELECT COUNT(*) FROM account_types;`
	err = db.QueryRow(stmt).Scan(&count)
	if err != nil {
		return err
	}

	var accountTypes *[]accountType.AccountType
	if count > 0 {
		accountTypes, err = GetAccountTypes(db)
		if err != nil {
			return err
		}
	} else {
		accountTypes = accountType.CreateDefaultAccountTypes()
		for _, at := range *accountTypes {
			stmt, err := db.Prepare("INSERT INTO account_types(id, name, side, ref_counter) values(?, ?, ?, ?)")
			if err != nil {
				return err
			}

			sId := at.GetId().String()
			name := at.GetName()
			side := at.GetSideNum()
			refCounter := at.GetRefCounter()
			_, err = stmt.Exec(sId, name, side, refCounter)
			if err != nil {
				return err
			}
		}
	}
	ledger.SetAccountTypes(accountTypes)

	return nil
}

func InitializeChartOfAccounts(db *sql.DB, chart *chart.ChartOfAccounts) error {
	var stmt = `CREATE TABLE IF NOT EXISTS accounts (
		id TEXT PRIMARY KEY,
		ref INTEGER NOT NULL UNIQUE,
		name VARCHAR(255) NOT NULL UNIQUE,
		account_type_id TEXT NOT NULL,
		FOREIGN KEY(account_type_id) REFERENCES account_types(id)
	);`
	_, err := db.Exec(stmt)
	if err != nil {
		return err
	}

	var accounts *[]account.Account
	accounts, err = GetAccounts(db)
	if err != nil {
		return err
	}
	chart.SetAccounts(accounts)
	ledger := chart.GetLedger()
	ledger.MapAccounts(accounts)
	
	return nil
}

func InitializeJournal(db *sql.DB) error {
	var stmt = `CREATE TABLE IF NOT EXISTS entries (
		id TEXT PRIMARY KEY,
		date TEXT NOT NULL,
		debit_account_id TEXT NOT NULL,
		credit_account_id TEXT NOT NULL,
		cents INTEGER NOT NULL,
		explanation VARCHAR(255) NOT NULL,
		posted INTEGER NOT NULL,
		FOREIGN KEY(debit_account_id) REFERENCES accounts(id),
		FOREIGN KEY(credit_account_id) REFERENCES accounts(id)
	)`
	_, err := db.Exec(stmt)
	if err != nil {
		return err
	}
	
	return nil
}

func GetAccountTypes(db *sql.DB) (*[]accountType.AccountType, error) {
	rows, err := db.Query("SELECT id, name, side, ref_counter FROM account_types")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var accountTypes []accountType.AccountType
	for rows.Next() {
		var sId, name string
		var side, refCounter int
		err := rows.Scan(&sId, &name, &side, &refCounter)
		if err != nil {
			return nil, err
		}
		accountType, err := accountType.CreateAccountTypeFromDb(sId, name, side, refCounter)
		if err != nil {
			return nil, err
		}

		accountTypes = append(accountTypes, accountType)
	}
	if rows.Err() != nil {
		return nil, rows.Err()
	}

	return &accountTypes, nil
}

func GetAccounts(db *sql.DB) (*[]account.Account, error) {
	rows, err := db.Query("SELECT id, ref, name, account_type_id FROM accounts")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	accounts := []account.Account{}
	for rows.Next() {
		var sId, name, sAccTypeId string
		var ref int
		err := rows.Scan(&sId, &ref, &name, &sAccTypeId)
		if err != nil {
			return nil, err
		}

		account, err := account.CreateAccountFromDb(sId, ref, name, sAccTypeId)
		if err != nil {
			return nil, err
		}

		accounts = append(accounts, account)
	}
	if rows.Err() != nil {
		return nil, rows.Err()
	}
	return &accounts, nil
}

func GetEntries(db *sql.DB, chart *chart.ChartOfAccounts, fromDate time.Time, toDate time.Time) (*[]*entry.Entry, error) {
	var rows *sql.Rows
	var err error
	if fromDate.IsZero() && toDate.IsZero() {
		rows, err = db.Query("SELECT id, date, debit_account_id, credit_account_id, cents, explanation, posted FROM entries")
	} else if fromDate.IsZero() {
		rows, err = db.Query("SELECT id, date, debit_account_id, credit_account_id, cents, explanation, posted FROM entries WHERE date < ?", toDate.Format(time.RFC3339))
	} else if toDate.IsZero() {
		rows, err = db.Query("SELECT id, date, debit_account_id, credit_account_id, cents, explanation, posted FROM entries WHERE date > ?", fromDate.Format(time.RFC3339))
	} else { 
		rows, err = db.Query("SELECT id, date, debit_account_id, credit_account_id, cents, explanation, posted FROM entries WHERE date BETWEEN ? AND ?", fromDate.Format(time.RFC3339), toDate.Format(time.RFC3339))
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var entries []*entry.Entry
	for rows.Next() {
		var sId, sDate, explanation, sDrId, sCrId string
		var cents, intPosted int
		err := rows.Scan(&sId, &sDate, &sDrId, &sCrId, &cents, &explanation, &intPosted)
		if err != nil {
			return nil, err
		}

		debitAccount, err := chart.GetAccountByStringId(sDrId)
		if err != nil {
			return nil, err
		}

		creditAccount, err := chart.GetAccountByStringId(sCrId)
		if err != nil {
			return nil, err
		}

		entry, err := entry.NewEntryFromDb(sId, sDate, debitAccount, creditAccount, cents, explanation, intPosted)
		if err != nil {
			return nil, err
		}

		entries = append(entries, entry)
	}
	if rows.Err() != nil {
		return nil, rows.Err()
	}

	return &entries, nil
}

func AddAccount(db *sql.DB, account *account.Account) error {
	stmt, err := db.Prepare("INSERT INTO accounts(id, ref, name, account_type_id) values(?,?,?,?)")
	if err != nil {
		return err
	}
	defer stmt.Close()

	sId := account.GetId().String()
	ref := account.GetRef()
	name := account.GetName()
	sAccTypeId := account.GetAccountTypeId().String()
	_, err = stmt.Exec(sId, ref, name, sAccTypeId)
	if err != nil {
		return err
	}

	return nil
}

func UpdateAccountTypeRefCounter(db *sql.DB, accountTypeId id.Id, refCounter int) error {
	sAccTypeId := accountTypeId.String()
	_, err := db.Exec("UPDATE account_types SET ref_counter = ? WHERE id = ?", refCounter, sAccTypeId)
	if err != nil {
		return err
	}

	return nil
}

func AddEntry(db *sql.DB, newEntry entry.Entry) error {
	stmt, err := db.Prepare("INSERT INTO entries(id, date, debit_account_id, credit_account_id, cents, explanation, posted) values(?, ?, ?, ?, ?, ?, ?)")
	if err != nil {
		return err
	}
	defer stmt.Close()

	sId := newEntry.GetId().String()
	date := newEntry.GetDate().Format(time.RFC3339)
	sDrId := newEntry.GetDebitAccount().GetId().String()
	sCrId := newEntry.GetCreditAccount().GetId().String()
	amount := int(newEntry.GetAmount())
	explanation := newEntry.GetExplanation()
	posted := newEntry.GetPostedInt()
	_, err = stmt.Exec(sId, date, sDrId, sCrId, amount, explanation, posted)
	if err != nil {
		return err
	}

	return nil
}