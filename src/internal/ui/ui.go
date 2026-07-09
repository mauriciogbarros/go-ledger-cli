package ui

import (
	"errors"
	"fmt"
	"time"

	"go.mod/internal/account"
	"go.mod/internal/journal"
)

func MenuNewAccount() (string, int, error) {
	var name string
	fmt.Print("Account name: ")
	_, err := fmt.Scanf("%s", &name)
	if err != nil {
		return "", 0, err
	}
	if len(name) == 0 {
		return "", 0, errors.New("Account name cannot be empty")
	}
	if len(name) > account.MaxNameLength {
		return "", 0, fmt.Errorf("Account name too long (max %d characters)", account.MaxNameLength)
	}
	var choice int
	fmt.Println("Choose the account type:")
	fmt.Println("1. Asset")
	fmt.Println("2. Liability")
	fmt.Println("3. Equity")
	fmt.Println("4. Revenue")
	fmt.Println("5. Expense")
	fmt.Println("────────────────────────")
	fmt.Print("Choice: ")
	_, err = fmt.Scanf("%d", &choice)
	if err != nil {
		return "", 0, err
	}
	if choice < 1 || choice > 5 {
		return "", 0, errors.New("Invalid account type")
	}

	return name, choice, nil
}

func DisplayJournal(journal journal.Journal, fromDate time.Time, toDate time.Time) {
	fmt.Println(journal)
	// TODO: move displaying entries from journal.String() to here
	// TODO: implement date filter
}

func DisplayAccounts(accounts []account.Account) {
	
}