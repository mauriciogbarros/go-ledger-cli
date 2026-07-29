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
	"go.mod/internal/ledger"
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

func InputAccountTypeRefPrefix(ledger *ledger.Ledger) (int, error) {
	if ledger == nil {
		return 0, errors.New("ledger is nil")
	}
	accountTypes := ledger.GetAccountTypes()
	if accountTypes == nil {
		return 0, errors.New("no account types available")
	}
	options := make([]int, 0)
	fmt.Println("Choose the account type:")
	for _, at := range *accountTypes {
		if at == nil {
			continue
		}
		options = append(options, at.GetRefPrefix())
		fmt.Printf("%d. %s\n", at.GetRefPrefix(), at.GetName())
	}
	fmt.Println(strings.Repeat("─", 3 + accountType.MaxNameLength))
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

	exists := false
	for _, o := range options {
		if choice == o {
			exists = true
			break
		}
	}
	if !exists {
		return 0, errors.New("Invalid choice.")
	}
	return choice, nil
}

func InputAmountF64() (float64, error) {
	fmt.Print("Amount: ")
	amountString, err := reader.ReadString('\n')
	if err != nil {
		return 0, err
	}
	amountString = strings.TrimSpace(amountString)
	amountString = strings.ReplaceAll(amountString, ",", "")
	amountF64, err := strconv.ParseFloat(amountString, 64)
	if err != nil {
		return 0, err
	}
	return amountF64, nil
}

func InputDate(dateField string) (time.Time, error) {
	fmt.Printf("%s date (YYYY-MM-DD): ", dateField)
	dateString, err := reader.ReadString('\n')
	if err != nil {
		return time.Time{}, err
	}
	dateString = strings.TrimSpace(dateString)
	if len(dateString) == 0 {
		return time.Time{}, nil
	}
	parts := strings.Split(dateString, "-")
	if len(parts) != 3 {
		return time.Time{}, fmt.Errorf("invalid date format, use YYYY-MM-DD")
	}
	layout := "2006-"
	if len(parts[1]) == 1 {
		layout += "1-"
	} else {
		layout += "01-"
	}
	if len(parts[2]) == 1 {
		layout += "2"
	} else {
		layout += "02"
	}
	date, err := time.Parse(layout, dateString)
	if err != nil {
		return time.Time{}, fmt.Errorf("invalid date format, use YYYY-MM-DD")
	}
	return date, nil
}

func InputAccountRef(ledger *ledger.Ledger) (int, error) {
	if ledger == nil {
		return 0, errors.New("ledger is nil")
	}
	accounts := ledger.GetAccounts()
	if accounts == nil {
		return 0, errors.New("no accounts available")
	}
	var menu strings.Builder
	menu.WriteString(" Ref   Accounts\n")
	menu.WriteString("─")
	menu.WriteString(strings.Repeat("─", 1 + 3))
	menu.WriteString("─┬─")
	menu.WriteString(strings.Repeat("─", account.MaxNameLength))
	menu.WriteString("─")
	menu.WriteString("\n")
	var refs []int
	for _, account := range *accounts {
		if account == nil {
			continue
		}
		ref := account.GetRef()
		refs = append(refs, ref)
		fmt.Fprintf(&menu, " %d │ %s\n", ref, account.GetName())
	}
	menu.WriteString("─")
	menu.WriteString(strings.Repeat("─", 1 + 3))
	menu.WriteString("─┴─")
	menu.WriteString(strings.Repeat("─", account.MaxNameLength))
	menu.WriteString("─")
	fmt.Println(menu.String())
	fmt.Print("Enter ref: ")
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

func InputText(fieldName string) (string, error) {
	fmt.Printf("%s: ", fieldName)
	explanation, err := reader.ReadString('\n')
	explanation = strings.TrimSpace(explanation)
	if err != nil {
		return "", err
	}
	if len(explanation) == 0 {
		explanation = ""
	}
	return explanation, nil
}