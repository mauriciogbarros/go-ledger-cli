package db

import (
	"database/sql"
	"strconv"
	"time"

	"go.mod/internal/account"
	"go.mod/internal/currency"
	"go.mod/internal/entry"
	"go.mod/internal/id"
	_ "modernc.org/sqlite"
)

func Initialize() error {
	db, err := sql.Open("sqlite", "./db/ledger.db")
	if err != nil {
		return err
	}
	defer db.Close()

	// Verify db connection
	if err := db.Ping();
	err != nil {
		return err
	}

	var stmt string
	stmt = `CREATE TABLE IF NOT EXISTS accounts (
		id INTEGER PRIMARY KEY,
		name VARCHAR(` + strconv.Itoa(account.MaxNameLength) +`) NOT NULL,
		type INTEGER NOT NULL
	);`
	_, err = db.Exec(stmt)
	if err != nil {
		return err
	}

	stmt = `CREATE TABLE IF NOT EXISTS entries (
		id TEXT PRIMARY KEY,
		date TEXT NOT NULL,
		debit_account_id INTEGER NOT NULL,
		credit_account_id INTER NOT NULL,
		amount INTEGER NOT NULL,
		explanation VARCHAR(` + strconv.Itoa(entry.MaxExplanationLength) + `) NOT NULL,
		FOREIGN KEY(debit_account_id) REFERENCES accounts(id),
		FOREIGN KEY(credit_account_id) REFERENCES accounts(id)
	)`
	_, err = db.Exec(stmt)
	if err != nil {
		return err
	}

	return nil
}

func GetAccounts() ([]account.Account, error) {
	db, err := sql.Open("sqlite", "./db/ledger.db")
	if err != nil {
		return nil, err
	}
	defer db.Close()

	rows, err := db.Query("SELECT id, name, type FROM accounts")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	accounts := []account.Account{}
	for rows.Next() {
		var id, accountType int
		var name string
		err := rows.Scan(&id, &name, &accountType)
		if err != nil {
			return nil, err
		}
		account := account.NewAccount(id, name, accountType)
		accounts = append(accounts, account)
	}

	return accounts, nil
}

func GetEntries() ([]entry.Entry, error) {
	db, err := sql.Open("sqlite", "./db/ledger.db")
	if err != nil {
		return nil, err
	}
	defer db.Close()

	rows, err := db.Query("SELECT id, date, debit_account_id, credit_account_id, amount, explanation FROM entries")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	entries := []entry.Entry{}
	for rows.Next() {
		var stringId, stringDate, explanation string
		var debitAccountId, creditAccountId, amount int
		err := rows.Scan(&stringId, &stringDate, &debitAccountId, &creditAccountId, &amount, &explanation)
		if err != nil {
			return nil, err
		}
		id, err := id.ParseId(stringId)
		if err != nil {
			return nil, err
		}
		date, err := time.Parse(time.RFC3339, stringDate)
		if err != nil {
			return nil, err
		}
		entry := entry.NewEntryFromDb(id, date, debitAccountId, creditAccountId, currency.Currency(amount), explanation)
		entries = append(entries, entry)
	}

	return entries, nil
}

func CreateAccount(newAccount account.Account) error {
	db, err := sql.Open("sqlite", "./db/ledger.db")
	if err != nil {
		return err
	}
	defer db.Close()

	stmt, err := db.Prepare("INSERT INTO accounts(id, name, type) values(?,?,?)")
	if err != nil {
		return err
	}
	defer stmt.Close()

	id := newAccount.GetRef()
	name := newAccount.GetName()
	accountType := newAccount.GetAccountTypeInt()
	_, err = stmt.Exec(id, name, accountType)
	if err != nil {
		return err
	}

	return nil
}

func CreateEntry(newEntry entry.Entry) error {
	db, err := sql.Open("sqlite", "./db/ledger.db")
	if err != nil {
		return err
	}
	defer db.Close()

	stmt, err := db.Prepare("INSERT INTO entries(id, date, debit_account_id, credit_account_id, amount, explanation) values(?, ?, ?, ?, ?, ?)")
	if err != nil {
		return err
	}
	defer stmt.Close()

	id := newEntry.GetId().String()
	date := newEntry.GetDate().Format(time.RFC3339)
	debitAccountId := newEntry.GetDebitAccountRef()
	creditAccountId := newEntry.GetCreditAccountRef()
	amount := int(newEntry.GetAmount())
	explanation := newEntry.GetExplanation()
	_, err = stmt.Exec(id, date, debitAccountId, creditAccountId, amount, explanation)
	if err != nil {
		return err
	}

	return nil
}