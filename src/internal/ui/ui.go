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
	"go.mod/internal/accountType"
	"go.mod/internal/chart"
	"go.mod/internal/id"
)

var reader = bufio.NewReader(os.Stdin)

func InputAccountName() (string, error) {
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

func InputAccountType(chart *chart.ChartOfAccounts) (id.Id, error) {
	ledger := chart.GetLedger()
	accountTypes := ledger.GetAccountTypes()
	options := make([]int, 0)
	fmt.Println("Choose the account type:")
	for _, at := range *accountTypes {
		option := at.GetRefGroup() / 1000
		options = append(options, option)
		fmt.Printf("%d. %s\n", option, at.GetName())
	}
	fmt.Println(strings.Repeat("─", 3 + accountType.MaxNameLength))
	fmt.Print("Choice: ")

	input, err := reader.ReadString('\n')
	if err != nil {
		return id.Id{}, err
	}
	input = strings.TrimSpace(input)
	choice, err := strconv.Atoi(input)
	if err != nil {
		return id.Id{}, errors.New("Invalid input")
	}

	isValidChoice := false
	for _, option := range options {
		if choice == option {
			isValidChoice = true
			break
		}
	}
	if !isValidChoice {
		return id.Id{}, errors.New("Invalid choice")
	}

	accountType, err := ledger.GetAccountTypeByRef(choice * 1000)
	if err != nil {
		return id.Id{}, err
	}
	
	accountTypeId := accountType.GetId()
	return accountTypeId, nil
}

func InputDate(from_to string) (time.Time, error) {
	fmt.Printf("%s date (YYYY-MM-DD): ", from_to)
	dateString, err := reader.ReadString('\n')
	if err != nil {
		return time.Time{}, err
	}
	dateString = strings.TrimSpace(dateString)
	if len(dateString) == 0 {
		return time.Time{}, nil
	}
	date, err := time.Parse("2006-01-02", dateString)
	if err != nil {
		return time.Time{}, err
	}
	return date, nil
}

func InputAccountRef(chart *chart.ChartOfAccounts, side int) (int, error) {
	accounts := chart.GetAccounts()

	width := 1 + 3 + 3 + account.MaxNameLength + 3 + 9 + 1
	var menu string = ""
	menu += " Ref   Accounts\n"
	menu += "─"
	menu += strings.Repeat("─", 3)
	menu += "─┬─"
	menu += strings.Repeat("─", account.MaxNameLength)
	menu += "─┬─"
	menu += strings.Repeat("─", 9 + 1)
	menu += "\n"
	var refs []int
	for _, account := range *accounts {
		refs = append(refs, account.GetRef())
		menu += account.String()
		menu += "\n"
	}

	menu += strings.Repeat("─", width)
	fmt.Println(menu)
	if side == 0 {
		fmt.Print("Enter debit account Ref: ")
	} else {
		fmt.Print("Enter credit account Ref: ")
	}
	input, err := reader.ReadString('\n')
	fmt.Println()
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

func InputExplanation() (string, error) {
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