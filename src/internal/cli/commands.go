package cli

import (
	"database/sql"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"go.mod/internal/database"
	"go.mod/internal/entry"
	"go.mod/internal/ledger"
	"go.mod/internal/ui"
)

func Run() (string, error) {
	args := os.Args[1:]
	fmt.Println()
	if len(args) == 0 {
		var msg string
		msg += "Usage: ledger <command> [args]\n"
		msg += "Enter \"ledger help\" to view commands."
		return "", errors.New(msg)
	}

	if args[0] == "help" {		
		return runHelp()
	}

	var err error
	var db *sql.DB
	db, err = sql.Open("sqlite", "./internal/database/ledger.db")
	if err != nil {
		return "", err
	}
	defer db.Close()
	err = database.Initialize(db)
	if err != nil {
		return "", err
	}

	ledger := ledger.NewLedger()
	ledger.CreateChart()
	ledger.CreateJournal()
	accountTypes, err := database.GetAccountTypes(db)
	if err != nil {
		return "", nil
	}
	ledger.SetChartAccountTypes(accountTypes)

	switch args[0] {
	case "view-ledger":
		return runViewLedger(db, ledger, args[1:])

	case "view-chart":
		return runViewChart(db, ledger, args[1:])
	
	case "view-types":
		return runViewTypes(db, ledger, args[1:])

	case "view-journal":
		return runViewJournal(db, ledger, args[1:])

	case "view-accounts":
		return runViewAccounts(db, ledger, args[1:])

	case "new-account":
		return runNewAccount(db, ledger, args[1:])

	case "new-entry":
		return runNewEntry(db, ledger, args[1:])
		
	default:
		return "Unknown command" + args[0], nil
	}
}

func runHelp() (string, error) {
		fmt.Println("Available commands:")
		fmt.Println("───────────────────")
		fmt.Println("help               - Show this help")
		fmt.Println("view-ledger [arg]  - View general ledger")
		fmt.Println("view-chart [arg]   - View chart of accounts")
		fmt.Println("view-types [arg]   - View account types information")
		fmt.Println("view-accounts [arg] - View account information")
		fmt.Println("view-journal [arg] - View journal entries")
		fmt.Println("new-account [arg]  - Add a new account to the chart")
		fmt.Println("new-entry [arg]    - Add a new entry to the journal")
		fmt.Println()

		return "Ok", nil
}

func runViewLedger(db *sql.DB, ledger *ledger.Ledger, args []string) (string, error) {
	accounts, err := database.GetAccounts(db)
	if err != nil {
		return "", err
	}
	err = ledger.SetAccounts(accounts)
	if err != nil {
		return "", err
	}

	switch len(args) {
	case 0:
		fmt.Println("To view entries from start of journal, press \"return\"")
		fromDate, err := ui.InputDate("From")
		if err != nil {
			return "", err
		}
		fmt.Println("To view entries until the end of journal, press \"return\"")
		toDate, err := ui.InputDate("To")
		if err != nil {
			return "", err
		}
		entries, err := database.GetEntriesPostedBetweenDates(db, ledger, fromDate, toDate)
		if err != nil {
			return "", err
		}
		ledger.SetJournalEntries(entries)
		output, err := ledger.String()
		if err != nil {
			return "", err
		}
		fmt.Println(output)
	
	case 1:
		switch args[0] {
		case "help":
			var msg strings.Builder
			msg.WriteString("View general ledger:\n")
			msg.WriteString("Argument options:\n")
			msg.WriteString("- <empty>: the complete general ledger, between dates (optionally)\n")
			msg.WriteString("- help: show this help\n")
			msg.WriteString("- type: view the ledger for an account type, between dates (optionally)\n")
			msg.WriteString("- account: view the ledger for an account, between dates (optionally)\n")

		case "type":
			fmt.Println("Ledger for account type - TBI")

		case "account":
			fmt.Println("Ledger for accoun - TBI")
		}

	default:
		return "", errors.New("Usage: ledger view-ledger [arg]")
	}

	return "Ok", nil
}

func runViewChart(db *sql.DB, ledger *ledger.Ledger, args [] string) (string, error) {
	switch len(args) {
	case 0:
		accounts, err := database.GetAccounts(db)
		if err != nil {
			return "", err
		}
		err = ledger.SetAccounts(accounts)
		if err != nil {
			return "", err
		}
		output, err := ledger.ViewChart(0)
		if err != nil {
			return "", err
		}
		fmt.Println(output)
		return "Ok", nil

	case 1:
		switch args[0] {
		case "help":
			var msg strings.Builder
			msg.WriteString("View chart of acccounts\n")
			msg.WriteString("Argument options:\n")
			msg.WriteString("- <empty>: view the complete chart of accounts\n")
			msg.WriteString("- help: show this help\n")
			msg.WriteString("- type: view the chart for a specific account type\n")
			fmt.Println(msg.String())
			return "Ok", nil

		case "type":
			refPrefix, err := ui.InputAccountTypeRefPrefix(ledger)
			if err != nil {
				return "", err
			}
			accountType, err := ledger.GetChart().GetAccountTypeByRefPrefix(refPrefix)
			if err != nil {
				return "", err
			}
			id := accountType.GetId()
			accounts, err := database.GetTypeAccounts(db, id)
			if err != nil {
				return "", err
			}
			err = ledger.SetAccounts(accounts)
			if err != nil {
				return "", err
			}
			output, err := ledger.ViewChart(refPrefix)
			if err != nil {
				return "", err
			}
			fmt.Println(output)
		
		default:
			return "", errors.New("Invalid argument")
		}
	default:
		return "", errors.New("Usage: ledger view-chart [arg]")
	}

	return "Ok", nil
}

func runViewTypes(db *sql.DB, ledger *ledger.Ledger, args []string) (string, error) {
	switch len(args) {
	case  0:		
		accounts, err := database.GetAccounts(db)
		if accounts == nil {
			return "", errors.New("no accounts available")
		}
		if err != nil {
			return "", err
		}
		err = ledger.SetAccounts(accounts)
		if err != nil {
			return "", err
		}
		fmt.Println(ledger.ViewAccountTypes())

	case 1:
		fmt.Println("View Account Type")
		fmt.Println("─────────────────")
		switch args[0] {
		case "help":
			var msg strings.Builder
			msg.WriteString("View account types information\n")
			msg.WriteString("Argument options:\n")
			msg.WriteString("<empty>: for all account types\n")
			msg.WriteString("type: for a specific account type\n")
			fmt.Println(msg.String())

		case "type":
			refPrefix, err := ui.InputAccountTypeRefPrefix(ledger)
			if err != nil {
				return "", err
			}
			accountType, err := ledger.GetChart().GetAccountTypeByRefPrefix(refPrefix)
			if err != nil {
				return "", err
			}
			accounts, err := database.GetTypeAccounts(db, accountType.GetId())
			if err != nil {
				return "", err
			}
			if err = ledger.SetAccounts(accounts); err != nil {
				return "", err
			}
			output, err := ledger.ViewAccountType(refPrefix)
			if err != nil {
				return "", err
			}
			fmt.Println(output)

		default:
			return "", errors.New("Invalid argument")
		}
		
	default:
		fmt.Println("View Account Type")
		fmt.Println("─────────────────")
		return "", errors.New("Usage: ledger view-types [arg]")
	}
	
	return "Ok", nil
}

func runViewAccounts(db *sql.DB, ledger *ledger.Ledger, args []string) (string, error) {
	switch len(args) {
	case 0:
		accounts, err := database.GetAccounts(db)
		if err != nil {
			return "", nil
		}

		err = ledger.SetAccounts(accounts)
		if err != nil {
			return "", nil
		}

		output, err := ledger.ViewAccounts()
		if err != nil {
			return "", err
		}

		fmt.Println(output)

	case 1:
		fmt.Println("View Accounts")
		fmt.Println("─────────────")
		switch args[0] {
		case "hellp":
			var msg strings.Builder
			msg.WriteString("Argument options:\n")
			msg.WriteString("- <empty> - show a listing of all accounts\n")
			msg.WriteString("- help    - show this help\n")
			msg.WriteString("- ref     - show the information of an account\n")
			fmt.Println(msg.String())

		case "ref":
			accounts, err := database.GetAccounts(db)
			if err != nil {
				return "", nil
			}

			err = ledger.SetAccounts(accounts)
			if err != nil {
				return "", nil
			}

			ref, err := ui.InputAccountRef(ledger)
			if err != nil {
				return "", nil
			}

			output, err := ledger.ViewAccountByRef(ref)
			if err != nil {
				return "", nil
			}
			
			fmt.Println(output)
		
		default:
			fmt.Println("View Accounts")
			fmt.Println("─────────────")
			return "", errors.New("Invalid argument")
		}

	default:
		return "", errors.New("Usage: ledger view-account [arg]")
	}

	return "Ok", nil
}

func runViewJournal(db *sql.DB, ledger *ledger.Ledger, args []string) (string, error) {
	var entries *[]*entry.Entry
	var err error
	var fromDate time.Time
	var toDate time.Time

	accounts, err := database.GetAccounts(db)
	if err != nil {
		return "", err
	}
	err = ledger.SetAccounts(accounts)
	if err != nil {
		return "", err
	}

	switch len(args) {
	case 0:
		entries, err = database.GetEntries(db, ledger)
		if err != nil {
			return "", err
		}
		ledger.SetJournalEntries(entries)
		fmt.Println(ledger.ViewJournal())

	case 1:
		switch args[0] {
		case "help":
			var msg strings.Builder
			msg.WriteString("View Journal\n")
			msg.WriteString("────────────")
			msg.WriteString("Argument options:\n")
			msg.WriteString("- <empty>    - Show all entries\n")
			msg.WriteString("- help       - Show this help\n")
			msg.WriteString("- dates      - Show entries from date to date\n")
			msg.WriteString("- posted     - Show all entries posted\n")
			msg.WriteString("- not-posted - Show all entries not posted\n")
			fmt.Println(msg.String())

		case "dates":
			fmt.Println("To view entries from start of journal, press \"return\"")
			fromDate, err = ui.InputDate("From")
			if err != nil {
				return "", err
			}
			fmt.Println("To view entries until the end of journal, press \"return\"")
			toDate, err = ui.InputDate("To")
			if err != nil {
				return "", err
			}

			entries, err = database.GetEntriesBetweenDates(db, ledger, fromDate, toDate)
			if err != nil {
				return "", err
			}
			ledger.SetJournalEntries(entries)
			fmt.Println(ledger.ViewJournal())

		case "posted":
			entries, err = database.GetEntriesPosted(db, ledger, true)
			if err != nil {
				return "", err
			}
			ledger.SetJournalEntries(entries)
			fmt.Println(ledger.ViewJournal())
			
		case "not-posted":
			entries, err = database.GetEntriesPosted(db, ledger, false)
			if err != nil {
				return "", err
			}
			ledger.SetJournalEntries(entries)
			fmt.Println(ledger.ViewJournal())

		default:
			return "", errors.New("Invalid argument")
		}
	default:
		return "", errors.New("Usage: ledger view-journal [arg]")
	}

	return "Ok", nil
}

func runNewAccount(db *sql.DB, ledger *ledger.Ledger, args []string) (string, error) {
	var name, description string
	var atRefPrefix int

	accounts, err := database.GetAccounts(db)
	if err != nil {
		return "", err
	}
	
	err = ledger.SetAccounts(accounts)
	if err != nil {
		return "", err
	}

	fmt.Println("Create New Account")
	fmt.Println("──────────────────")
	switch len(args) {
	case 0:
		name, err = ui.InputAccountName()
		if err != nil {
			return "", err
		}

		atRefPrefix, err = ui.InputAccountTypeRefPrefix(ledger)
		if err != nil {
			return "", err
		}

		description, err = ui.InputText("Description")
		if err != nil {
			return "", err
		}

	case 1:
		if args[0] == "help" {
			var msg strings.Builder
			msg.WriteString("Create a new account:\n")
			msg.WriteString("Argument options:\n")
			msg.WriteString("- help: show this help\n")
			msg.WriteString("- <name>: create new account with <name>\n")
			msg.WriteString("- <empty>: create new account\n")
			fmt.Println(msg.String())
			return "Ok", nil
		} else {
			name = args[0]
			atRefPrefix, err = ui.InputAccountTypeRefPrefix(ledger)
			if err != nil {
				return "", err
			}

			description, err = ui.InputText("Description")
			if err != nil {
				return "", err
			}
		}

	default:
		return "", errors.New("Usage: ledger new-account [arg]")
	}

	newAccount, err := ledger.CreateAccount(atRefPrefix, name, description)
	if err != nil {
		return "", err
	}

	err = database.AddAccount(db, newAccount)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("Account created: %s (%d)", newAccount.GetName(), newAccount.GetRef()), nil
}

func runNewEntry(db *sql.DB, ledger *ledger.Ledger, args []string) (string, error) {
	var amountF64 float64
	var err error
	fmt.Println("New Journal Entry")
	fmt.Println("─────────────────")
	switch len(args) {
	case 0:
		amountF64, err = ui.InputAmountF64()
		if err != nil {
			return "", err
		}

	case 1:
		amountF64, err = strconv.ParseFloat(args[0], 64)
		if err != nil {
			return "", err
		}

	default:
		return "", errors.New("Usage: ledger new-entry <amount")
	}

	date, err := ui.InputDate("Entry")
	if err != nil {
		return "", err
	}

	accounts, err := database.GetAccounts(db)
	if err != nil {
		return "", err
	}
	
	err = ledger.SetAccounts(accounts)
	if err != nil {
		return "", err
	}

	fmt.Println("Debit account")
	fmt.Println("─────────────")
	debitAccountRef, err := ui.InputAccountRef(ledger)
	if err != nil {
		return "", err
	}

	fmt.Println("Credit account")
	fmt.Println("──────────────")
	creditAccountRef, err := ui.InputAccountRef(ledger)
	if err != nil {
		return "", err
	}

	if debitAccountRef == creditAccountRef {
		return "", errors.New("Invalid entry: debit and credit accounts must be different.")
	}

	explanation, err := ui.InputText("Explanation")
	if err != nil {
		return "", err
	}

	newEntry, err := ledger.CreateJournalEntry(date, debitAccountRef, creditAccountRef, amountF64, explanation)
	if err != nil {
		return "", err
	}

	err = database.AddEntry(db, &newEntry)
	if err != nil {
		return "", err
	}

	debitAccount := newEntry.GetDebitAccount()
	creditAccount := newEntry.GetCreditAccount()
	if debitAccount == nil || creditAccount == nil {
		return "", errors.New("entry is missing debit or credit account")
	}
	return "Entry created: Dr = " + debitAccount.GetName() + "; Cr = " + creditAccount.GetName() + "; $" + newEntry.GetAmount().String(), nil
}
