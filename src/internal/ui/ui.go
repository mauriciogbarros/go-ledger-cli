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
	accountTypes := ledger.GetAccountTypes()
	options := make([]int, 0)
	fmt.Println("Choose the account type:")
	for _, at := range *accountTypes {
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

func InputDate(dateField string) (time.Time, error) {
	fmt.Printf("%s date (YYYY-MM-DD HH:MM:SS): ", dateField)
	dateString, err := reader.ReadString('\n')
	if err != nil {
		return time.Time{}, err
	}
	dateString = strings.TrimSpace(dateString)
	if len(dateString) == 0 {
		return time.Time{}, nil
	}
	date, err := time.Parse(time.DateTime, dateString)
	if err != nil {
		return time.Time{}, err
	}
	return date, nil
}

func InputAccountRef(ledger *ledger.Ledger, side string) (int, error) {
	accounts := ledger.GetAccounts()
	width := 1 + 3 + 3 + account.MaxNameLength + 3 + 9 + 1
	var menu strings.Builder
	menu.WriteString(" Ref   Accounts\n")
	menu.WriteString("─")
	menu.WriteString(strings.Repeat("─", 3))
	menu.WriteString("─┬─")
	menu.WriteString(strings.Repeat("─", account.MaxNameLength))
	menu.WriteString("─┬─")
	menu.WriteString(strings.Repeat("─", 9 + 1))
	menu.WriteString("\n")
	var refs []int
	for _, account := range *accounts {
		refs = append(refs, account.GetRef())
		menu.WriteString(account.String())
		menu.WriteString("\n")
	}

	menu.WriteString(strings.Repeat("─", width))
	menu.WriteString("\n")
	fmt.Fprintf(&menu, "Enter %s account Reef: ", side)
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