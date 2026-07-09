package storage

import (
	"go.mod/db"
	"go.mod/internal/account"
	"go.mod/internal/ui"
)

var accounts []account.Account

func AddAccount() error {
	name, tAccount, err := ui.MenuNewAccount()
	if err != nil {
		return err
	}

	newAccount := account.NewAccount(name, tAccount)
	err = db.CreateAccount(newAccount)
	if err != nil {
		return err
	}
	
	accounts = append(accounts, newAccount)
	return nil
}

func GetAccounts() []account.Account {
	return accounts
}

func GetAccountsByType(t account.AccountType) []account.Account {
	var tAccounts []account.Account
	for _, a := range accounts {
		if a.Type == t {
			tAccounts = append(tAccounts, a)
		}
	}
	return tAccounts
}