package database

import (
	"database/sql"
	"time"

	"go.mod/internal/account"
	"go.mod/internal/accountType"
	"go.mod/internal/entry"
	"go.mod/internal/id"
	"go.mod/internal/ledger"
	_ "modernc.org/sqlite"
)

func Initialize(db *sql.DB) error {
	var err error
	err = initializeAccountTypes(db)
	if err != nil {
		return err
	}

	err = initializeAccounts(db)
	if err != nil {
		return err
	}

	err = initializeJournal(db)
	if err != nil {
		return err
	}

	return nil
}

func initializeAccountTypes(db *sql.DB) error {
	var stmt = `CREATE TABLE IF NOT EXISTS account_types (
		id TEXT PRIMARY KEY,
		name VARCHAR(255) NOT NULL UNIQUE,
		ref_prefix INTEGER NOT NULL UNIQUE
	);`
	_, err := db.Exec(stmt)
	if err != nil {
		return err
	}

	return nil
}

func initializeAccounts(db *sql.DB) error {
	var stmt = `CREATE TABLE IF NOT EXISTS accounts (
		id TEXT PRIMARY KEY,
		ref INTEGER NOT NULL UNIQUE,
		name VARCHAR(255) NOT NULL UNIQUE,
		description TEXT,
		account_type_id TEXT NOT NULL,
		FOREIGN KEY(account_type_id) REFERENCES account_types(id)
	);`
	_, err := db.Exec(stmt)
	if err != nil {
		return err
	}
	
	return nil
}

func initializeJournal(db *sql.DB) error {
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

func GetAccountTypes(db *sql.DB) (*map[id.Id]*accountType.AccountType, error) {
	var count int
	var err error
	var stmt string

	stmt = `SELECT COUNT(*) FROM account_types;`
	err = db.QueryRow(stmt).Scan(&count)
	if err != nil {
		return nil, err
	}
	
	accountTypes := make(map[id.Id]*accountType.AccountType,0)
	
	if count > 0 {
		rows, err := db.Query("SELECT id, name, ref_prefix FROM account_types")
		if err != nil {
			return nil, err
		}
		defer rows.Close()
		
		for rows.Next() {
			var sId, name string
			var refPrefix int
			err := rows.Scan(&sId, &name, &refPrefix)
			if err != nil {
				return nil, err
			}
			accountType, err := accountType.NewDbAccountType(sId, name, refPrefix)
			if err != nil {
				return nil, err
			}
			accountTypes[accountType.GetId()] = accountType
		}
		if rows.Err() != nil {
			return nil, rows.Err()
		}
	} else {
		accountTypes = accountType.GetDefaultAccountTypes()
		for _, at := range accountTypes {
			stmt, err := db.Prepare("INSERT INTO account_types(id, name, ref_prefix) values(?, ?, ?)")
			if err != nil {
				return nil, err
			}

			sId := at.GetId().String()
			name := at.GetName()
			refPrefix := at.GetRefPrefix()
			_, err = stmt.Exec(sId, name, refPrefix)
			if err != nil {
				return nil, err
			}
		}
	}

	return &accountTypes, nil
}

func GetAccounts(db *sql.DB) (*[]*account.Account, error) {
	var rows *sql.Rows
	rows, err := db.Query("SELECT id, ref, name, description, account_type_id FROM accounts")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	accounts := make([]*account.Account, 0)
	for rows.Next() {
		var sId, name, sAtId string
		var description *string
		var ref int
		err := rows.Scan(&sId, &ref, &name, &description, &sAtId)
		if err != nil {
			return nil, err
		}
		desc := ""
		if description != nil {
			desc = *description
		}

		account, err := account.NewDbAccount(sId, ref, name, desc, sAtId)
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

func GetTypeAccounts(db *sql.DB, id id.Id) (*[]*account.Account, error) {
	var rows *sql.Rows
	rows, err := db.Query("SELECT id, ref, name, description, account_type_id FROM accounts WHERE account_type_id = ?", id.String())
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	accounts := make([]*account.Account, 0)
	for rows.Next() {
		var sId, name, sAtId string
		var description *string
		var ref int
		err := rows.Scan(&sId, &ref, &name, &description, &sAtId)
		if err != nil {
			return nil, err
		}
		d := ""
		if description != nil {
			d = *description
		}

		account, err := account.NewDbAccount(sId, ref, name, d, sAtId)
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

func GetEntries(db *sql.DB, ledger *ledger.Ledger) (*[]*entry.Entry, error) {
	rows, err := db.Query("SELECT id, date, debit_account_id, credit_account_id, cents, explanation, posted FROM entries")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var entries []*entry.Entry
	for rows.Next() {
		var sId, sDate, explanation, sDrId, sCrId string
		var cents, intIsPosted int
		err := rows.Scan(&sId, &sDate, &sDrId, &sCrId, &cents, &explanation, &intIsPosted)
		if err != nil {
			return nil, err
		}

		debitAccount, err := ledger.GetAccountByStringId(sDrId)
		if err != nil {
			return nil, err
		}

		creditAccount, err := ledger.GetAccountByStringId(sCrId)
		if err != nil {
			return nil, err
		}

		entry, err := entry.NewDbEntry(sId, sDate, debitAccount, creditAccount, cents, explanation, intIsPosted)
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

func GetEntriesBetweenDates(db *sql.DB, ledger *ledger.Ledger, fromDate time.Time, toDate time.Time) (*[]*entry.Entry, error) {
	var rows *sql.Rows
	var err error
	if fromDate.IsZero() && toDate.IsZero() {
		rows, err = db.Query("SELECT id, date, debit_account_id, credit_account_id, cents, explanation, posted FROM entries")
	} else if fromDate.IsZero() {
		rows, err = db.Query("SELECT id, date, debit_account_id, credit_account_id, cents, explanation, posted FROM entries WHERE date < ?", toDate.Format(time.DateTime))
	} else if toDate.IsZero() {
		rows, err = db.Query("SELECT id, date, debit_account_id, credit_account_id, cents, explanation, posted FROM entries WHERE date > ?", fromDate.Format(time.DateTime))
	} else { 
		rows, err = db.Query("SELECT id, date, debit_account_id, credit_account_id, cents, explanation, posted FROM entries WHERE date BETWEEN ? AND ?", fromDate.Format(time.DateTime), toDate.Format(time.DateTime))
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

		debitAccount, err := ledger.GetAccountByStringId(sDrId)
		if err != nil {
			return nil, err
		}

		creditAccount, err := ledger.GetAccountByStringId(sCrId)
		if err != nil {
			return nil, err
		}

		entry, err := entry.NewDbEntry(sId, sDate, debitAccount, creditAccount, cents, explanation, intPosted)
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

func GetEntriesPostedBetweenDates(db *sql.DB, ledger *ledger.Ledger, fromDate time.Time, toDate time.Time) (*[]*entry.Entry, error) {
	var rows *sql.Rows
	var err error
	if fromDate.IsZero() && toDate.IsZero() {
		rows, err = db.Query("SELECT id, date, debit_account_id, credit_account_id, cents, explanation, posted FROM entries WHERE posted = 1")
	} else if fromDate.IsZero() {
		rows, err = db.Query("SELECT id, date, debit_account_id, credit_account_id, cents, explanation, posted FROM entries WHERE posted = 1 AND date < ?", toDate.Format(time.DateTime))
	} else if toDate.IsZero() {
		rows, err = db.Query("SELECT id, date, debit_account_id, credit_account_id, cents, explanation, posted FROM entries WHERE posted = 1 AND date > ?", fromDate.Format(time.DateTime))
	} else { 
		rows, err = db.Query("SELECT id, date, debit_account_id, credit_account_id, cents, explanation, posted FROM entries WHERE posted = 1 AND date BETWEEN ? AND ?", fromDate.Format(time.DateTime), toDate.Format(time.DateTime))
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

		debitAccount, err := ledger.GetAccountByStringId(sDrId)
		if err != nil {
			return nil, err
		}

		creditAccount, err := ledger.GetAccountByStringId(sCrId)
		if err != nil {
			return nil, err
		}

		entry, err := entry.NewDbEntry(sId, sDate, debitAccount, creditAccount, cents, explanation, intPosted)
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

func GetEntriesPosted(db *sql.DB, ledger *ledger.Ledger, arePosted bool) (*[]*entry.Entry, error) {
	var intPosted int
	if arePosted {
		intPosted = 1
	} else {
		intPosted = 0
	}
	rows, err := db.Query("SELECT id, date, debit_account_id, credit_account_id, cents, explanation, posted FROM entries WHERE posted = ?", intPosted)
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

		debitAccount, err := ledger.GetAccountByStringId(sDrId)
		if err != nil {
			return nil, err
		}

		creditAccount, err := ledger.GetAccountByStringId(sCrId)
		if err != nil {
			return nil, err
		}

		entry, err := entry.NewDbEntry(sId, sDate, debitAccount, creditAccount, cents, explanation, intPosted)
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
	stmt, err := db.Prepare("INSERT INTO accounts(id, ref, name, description, account_type_id) values(?, ?, ?, ?, ?)")
	if err != nil {
		return err
	}
	defer stmt.Close()

	sId := account.GetId().String()
	ref := account.GetRef()
	name := account.GetName()
	description := account.GetDescription()
	sAtId := account.GetAccountTypeId().String()
	_, err = stmt.Exec(sId, ref, name, description, sAtId)
	if err != nil {
		return err
	}

	return nil
}

func AddEntry(db *sql.DB, entry entry.Entry) error {
	stmt, err := db.Prepare("INSERT INTO entries(id, date, debit_account_id, credit_account_id, cents, explanation, posted) values(?, ?, ?, ?, ?, ?, ?)")
	if err != nil {
		return err
	}
	defer stmt.Close()

	sId := entry.GetId().String()
	date := entry.GetDate().Format(time.DateTime)
	sDrId := entry.GetDebitAccount().GetId().String()
	sCrId := entry.GetCreditAccount().GetId().String()
	amount := int(entry.GetAmount())
	explanation := entry.GetExplanation()
	posted := 0
	if entry.IsPosted() {
		posted = 1
	}
	_, err = stmt.Exec(sId, date, sDrId, sCrId, amount, explanation, posted)
	if err != nil {
		return err
	}

	return nil
}