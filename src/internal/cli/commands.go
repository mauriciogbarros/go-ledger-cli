package cli

import (
	"database/sql"
	"errors"
	"fmt"
	"os"
	"strconv"

	"go.mod/db"
	"go.mod/internal/account"
	"go.mod/internal/chart"
	"go.mod/internal/id"
	"go.mod/internal/ledger"
	"go.mod/internal/ui"
)

func Run() (string, error) {
	args := os.Args[1:]
	fmt.Println()
	if len(args) == 0 {
		fmt.Println("Usage: ledger <command> [args]")
		fmt.Println("Enter \"help\" to view commands.")
		return "Ok", nil
	}

	if args[0] == "help" {
		runHelp()
		return "Ok", nil
	}

	var err error
	var database *sql.DB
	database, err = sql.Open("sqlite", "./db/ledger.db")
	if err != nil {
		return "", err
	}
	defer database.Close()

	var ledger *ledger.Ledger = ledger.NewLedger("General Ledger")
	err = db.InitializeLedger(database, ledger)
	if err != nil {
		return "", err
	}

	var chart *chart.ChartOfAccounts = chart.NewChartOfAccounts("Chart of Accounts", ledger)
	err = db.InitializeChartOfAccounts(database, chart)
	if err != nil {
		return "", err
	}

	switch args[0] {
	case "view-types":
		fmt.Println(ledger.PrintTypes())
		return "Ok", nil

	case "view-chart":
		fmt.Println(chart)
		return "Ok", nil

	case "view-journal":
		return runViewJournal(database, chart)

	case "new-account":
		return runNewAccount(chart, database, args[1:])

	case "new-entry":
		return runNewEntry(chart, args[1:])
		
	default:
		return "Unknown command" + args[0], nil
	}
}

func runHelp() {
		fmt.Println("Available commands:")
		fmt.Println("───────────────────")
		fmt.Println("commands           - Show this help")
		fmt.Println("view-types         - Show account types")
		fmt.Println("view-journal       - Show journal entries")
		fmt.Println("view-chart         - Show chart of accounts")
		fmt.Println("new-account [args] - Add a new account to the chart")
		fmt.Println("new-entry <amount> - Add a new entry to the journal")
		fmt.Println()

}

func runViewJournal(database *sql.DB, chart *chart.ChartOfAccounts) (string, error) {
	fmt.Println("View journal - to be implemented")
	return "Ok", nil
}

func runNewAccount(chart *chart.ChartOfAccounts, database *sql.DB, args []string) (string, error) {
	var name string
	var accTypeId id.Id
	var err error
	switch len(args) {
	case 0:
		name, err = ui.InputAccountName()
		if err != nil {
			return "", err
		}

		accTypeId, err = ui.InputAccountType(chart)
		if err != nil {
			return "", err
		}

	case 1:
		if len(args[0]) <= 0 || len(args[0]) > account.MaxNameLength {
			return "", errors.New("Invalid accocunt name")
		}
		name = args[0]

		accTypeId, err = ui.InputAccountType(chart)
		if err != nil {
			return "", err
		}

	default:
		return "", errors.New("Usage: ledger new-account [name]")
	}

	newAccount := account.NewAccount(name, accTypeId)
	ledger := chart.GetLedger()
	accountType, err := ledger.GetAccountTypeById(accTypeId)
	if err != nil {
		return "", err
	}

	err = chart.AddAccount(newAccount)
	if err != nil {
		return "", err
	}

	err = accountType.AddAccount(newAccount)
	if err != nil {
		id := newAccount.GetId()
		chart.RemoveAccount(id)
		newAccount = nil
		return "", err
	}

	err = db.AddAccount(database, newAccount)
	if err != nil {
		id := newAccount.GetId()
		chart.RemoveAccount(id)
		newAccount = nil
		return "", err
	}

	db.UpdateAccountTypeRefCounter(database, accTypeId, accountType.GetRefCounter())
	return "Account created: " + strconv.Itoa(newAccount.GetRef()) + " - " + newAccount.GetName(), nil
}

func runNewEntry(chart *chart.ChartOfAccounts, args []string) (string, error) {
	fmt.Println("New entry - to be implemented")
	return "Ok", nil
	// if len(args) != 1 {
	// 	return "Usage: ledger new-entry <ammount>", nil
	// }
	// amount64, err := strconv.ParseFloat(args[0], 64)
	// if err != nil {
	// 	return "", err
	// }
	// amount := currency.Convert64(amount64)
	// debitAccountRef, err := ui.MenuGetAccount(chart, 0)
	// if err != nil {
	// 	return "", err
	// }
	// creditAccountRef, err := ui.MenuGetAccount(chart, 1)
	// if err != nil {
	// 	return "", err
	// }
	// if debitAccountRef == creditAccountRef {
	// 	return "", errors.New("Invalid entry: debit and credit accounts must be different.")
	// }		
	// explanation, err := ui.MenuGetExplanation()
	// fmt.Println(len(explanation))
	// if err != nil {
	// 	return "", err
	// }
	// dr, err := chart.GetAccountByRef(debitAccountRef)
	// if err != nil {
	// 	return "", err
	// }
	// drName := dr.GetName()
	// cr, err := chart.GetAccountByRef(creditAccountRef)
	// if err != nil {
	// 	return "", err
	// }
	// crName := cr.GetName()
	// newEntry := entry.NewEntry(debitAccountRef, creditAccountRef, amount, explanation)
	// err = journal.AddEntry(newEntry)
	// if err != nil {
	// 	return "", err
	// }

	// return "Entry created: Dr = " + drName + "; Cr = " + crName + "; $" + amount.String(), nil
}
