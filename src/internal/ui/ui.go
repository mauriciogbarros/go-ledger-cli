package ui

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"slices"
	"strconv"
	"strings"
	"time"

	"go.mod/internal/account"
	"go.mod/internal/chart"
	"go.mod/internal/journal"
)

var reader = bufio.NewReader(os.Stdin)

func MenuNewAccountName() (string, error) {
	fmt.Print("Account name: ")
	name, err := reader.ReadString('\n')
	name = strings.TrimSpace(name)
	if err != nil {
		return "", err
	}
	if len(name) == 0 {
		return "", errors.New("Account name cannot be empty")
	}
	if len(name) > account.MaxNameLength {
		return "", fmt.Errorf("Account name too long (max %d characters)", account.MaxNameLength)
	}

	return name, nil
}

func MenuNewAccountType() (int, error) {
	fmt.Println("Choose the account type:")
	fmt.Println("1. Asset")
	fmt.Println("2. Liability")
	fmt.Println("3. Equity")
	fmt.Println("4. Revenue")
	fmt.Println("5. Expense")
	fmt.Println("────────────────────────")
	fmt.Print("Choice: ")

	input, err := reader.ReadString('\n')
	if err != nil {
		return 0, err
	}
	input = strings.TrimSpace(input)
	choice, err := strconv.Atoi(input)
	if err != nil {
		return 0, errors.New("Invalid input")
	}
	if choice < 1 || choice > 5 {
		return 0, errors.New("Invalid account type")
	}

	return choice, nil
}

func DisplayJournal(journal journal.Journal, fromDate time.Time, toDate time.Time) {
	fmt.Println(journal)
	// TODO: move displaying entries from journal.String() to here
	// TODO: implement date filter
}

func MenuGetAccount(chart *chart.ChartOfAccounts) (int, error) {
	accounts := chart.GetAccounts()
	width := 1 + 3 + 3 + account.MaxNameLength + 3 + 9 + 1
	fmt.Println("Accounts")
	fmt.Println(strings.Repeat("─", width))
	var refs []int
	for _, account := range accounts {
		refs = append(refs, account.GetRef())
		fmt.Println(account.String())
	}
	fmt.Println(strings.Repeat("─", width))
	fmt.Print("Enter account Ref: ")
	input, err := reader.ReadString('\n')
	if err != nil {
		return 0, err
	}
	input = strings.TrimSpace(input)
	ref, err := strconv.Atoi(input)
	if err != nil {
		return 0, err
	}
	index := slices.Index(refs, ref)
	if index < 0 {
		return 0, errors.New("Invalid account reference")
	}
	return ref, nil
}

func MenuGetExplanation() (string, error) {
	fmt.Print("Explanation: ")
	explanation, err := reader.ReadString('\n')
	explanation = strings.TrimSpace(explanation)
	if err != nil {
		return "", err
	}
	if len(explanation) == 0 {
		return "", errors.New("Explanation cannot be empty")
	}
	return explanation, nil
}