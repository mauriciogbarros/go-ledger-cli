package cli

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"

	"go.mod/db"
	"go.mod/internal/account"
	"go.mod/internal/chart"
	"go.mod/internal/currency"
	"go.mod/internal/entry"
	"go.mod/internal/journal"
	"go.mod/internal/ui"
)

func Run() (string, error) {
	var err error
	err = db.Initialize()
	if err != nil {
		return "", err
	}
	var chart chart.ChartOfAccounts
	var journal journal.Journal
	err = chart.NewChartOfAccounts("Chart of Accounts")
	if err != nil {
		return "", err
	}
	err = journal.NewJournal("General Journal", &chart)
	if err != nil {
		return "", err
	}

	args := os.Args[1:]
	fmt.Println()
	if len(args) == 0 {
		return "Usage: ledger <command> [args]", nil
	}

	switch args[0] {

	case "commands":
		fmt.Println("Available commands:")
		fmt.Println("───────────────────")
		fmt.Println("commands           - Show this help")
		fmt.Println("view-journal       - Show journal entries")
		fmt.Println("view-chart         - Show chart of accounts")
		fmt.Println("new-account [args] - Add a new account to the chart")
		fmt.Println("new-entry <amount> - Add a new entry to the journal")
		fmt.Println()
		return "Ok", nil

	case "view-chart":
		fmt.Println(chart)
		return "Ok", nil

	case "view-journal":
		fmt.Println(journal)
		return "Ok", nil

	case "new-account":
		msg, err := runNewAccount(&chart, args[1:])
		return msg, err

	case "new-entry":
		msg, err := runNewEntry(&chart, &journal, args[1:])
		return msg, err
		
	default:
		return "Unknown command" + args[0], nil
	}
}

func runNewAccount(chart *chart.ChartOfAccounts, args []string) (string, error) {
	var name string
	var accType int
	var err error
	switch len(args) {
	case 0:
		name, err = ui.MenuNewAccountName()
		if err != nil {
			return "", err
		}
		accType, err = ui.MenuNewAccountType()
		if err != nil {
			return "", err
		}
	case 1:
		if len(args[0]) <= 0 || len(args[0]) > account.MaxNameLength {
			return "", errors.New("Invalid accocunt name")
		}
		name = args[0]
		accType, err = ui.MenuNewAccountType()
		if err != nil {
			return "", err
		}
	case 2:
		if len(args[0]) <= 0 || len(args[0]) > account.MaxNameLength {
			return "", errors.New("Invalid accocunt name")
		}
		name = args[0]
		accType, err = strconv.Atoi(args[1])
		if err != nil {
			t := strings.ToLower(args[1])
			switch t {
			case "asset":
				accType = 1
			case "liability":
				accType = 2
			case "equity":
				accType = 3
			case "revenue":
				accType = 4
			case "expense":
				accType = 5
			default:
				return "", errors.New("Invalid account type")
			}
		}
	default:
		return "", errors.New("Usage: ledger new-account [name] [type]")
	}
	newAccount, err := chart.AddAccount(name, accType)
	if err != nil {
		return "", err
	}
	return "Account created: " + strconv.Itoa(newAccount.GetRef()) + " - " + newAccount.GetName(), nil
}

func runNewEntry(chart *chart.ChartOfAccounts, journal *journal.Journal, args []string) (string, error) {
	if len(args) != 1 {
		return "Usage: ledger new-entry <ammount>", nil
	}
	amount64, err := strconv.ParseFloat(args[0], 64)
	if err != nil {
		return "", err
	}
	amount := currency.Convert64(amount64)
	debitAccountRef, err := ui.MenuGetAccount(chart)
	if err != nil {
		return "", err
	}
	creditAccountRef, err := ui.MenuGetAccount(chart)
	if err != nil {
		return "", err
	}
	if debitAccountRef == creditAccountRef {
		return "", errors.New("Invalid entry: debit and credit accounts must be different.")
	}		
	explanation, err := ui.MenuGetExplanation()
	fmt.Println(len(explanation))
	if err != nil {
		return "", err
	}
	dr, err := chart.GetAccountByRef(debitAccountRef)
	if err != nil {
		return "", err
	}
	drName := dr.GetName()
	cr, err := chart.GetAccountByRef(creditAccountRef)
	if err != nil {
		return "", err
	}
	crName := cr.GetName()
	newEntry := entry.NewEntry(debitAccountRef, creditAccountRef, amount, explanation)
	err = journal.AddEntry(newEntry)
	if err != nil {
		return "", err
	}

	return "Entry created: Dr: " + drName + ", Cr: " + crName + " = $" + amount.String(), nil
}