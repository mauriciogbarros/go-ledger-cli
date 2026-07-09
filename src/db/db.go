package db

import (
	"database/sql"
	"strconv"

	"go.mod/internal/account"
	"go.mod/internal/entry"
	"go.mod/internal/id"
	_ "modernc.org/sqlite"
)

func Initialize () error {
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
		id TEXT PRIMARY KEY,
		name VARCHAR(` + strconv.Itoa(account.MaxNameLength) +`) NOT NULL,
		type TEXT NOT NULL
	);`
	_, err = db.Exec(stmt)
	if err != nil {
		return err
	}

	stmt = `CREATE TABLE IF NOT EXISTS entries (
		id TEXT PRIMARY KEY,
		date TEXT NOT NULL,
		debit_account_id TEXT NOT NULL,
		credit_account_id TEXT NOT NULL,
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

	_, err = stmt.Exec(newAccount.ID, newAccount.Name, newAccount.Type)
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
		var stringId, name, type_ string
		err := rows.Scan(&stringId, &name, &type_)
		if err != nil {
			return nil, err
		}
		parsedId, err := id.FromString(stringId)
		if err != nil {
			return nil, err
		}
		accounts = append(accounts, account.Account{
			ID:   id.Id(parsedId),
			Name: name,
			Type: account.AccountType(0),
		})
	}

	return accounts, nil
}