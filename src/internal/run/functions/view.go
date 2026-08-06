package functions

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"go.mod/internal/database"
	"go.mod/internal/input"
	"go.mod/internal/ledger"
)

func View(
	db *sql.DB,
	ledger *ledger.Ledger,
	args []string,
) error {
	if len(args) == 0 {
		return errors.New("Usage: ledger view <command> [arg]")
	}

	switch args[0] {
	case "help":
		viewHelp()
		return nil

	case "ledger":
		return viewLedger(db, ledger, args)

	case "chart":
		return viewChart(db, ledger, args)

	case "types":
		return viewTypes(db, ledger, args)

	case "type":
		return viewType(db, ledger, args)

	case "accounts":
		return viewAccounts(db, ledger, args)

	case "account":
		return viewAccount(db, ledger, args)

	case "journal":
		return viewJournal(db, ledger, args)

	default:
		return fmt.Errorf("Invalid command: %s\n", args[0])
	}
}

func viewHelp() {
	var msg strings.Builder
	msg.WriteString("Usage: ledger view <command> [arg]\n")
	msg.WriteString("Commands:\n")
	msg.WriteString("─────────\n")
	msg.WriteString("ledger [arg]  => View the General Ledger\n")
	msg.WriteString("chart [arg]   => View the Chart of Accounts\n")
	msg.WriteString("types         => View summary of all account types\n")
	msg.WriteString("type [arg]    => View details for an account type\n")
	msg.WriteString("accounts      => View summary of all accounts\n")
	msg.WriteString("account [arg] => View details for an account\n")
	msg.WriteString("journal [arg] => View the General Journal\n")
	fmt.Print(msg.String())
}

func viewLedger(db *sql.DB, ledger *ledger.Ledger, args []string) error {
	accounts, err := database.GetAccounts(db)
	if err != nil {
		return err
	}
	err = ledger.SetAccounts(accounts)
	if err != nil {
		return err
	}

	switch len(args) {
	case 1:
		return viewLedgerInfo(db, ledger)

	case 2:
		switch args[1] {
		case "help":
			viewLedgerHelp()
			return nil

		case "all":
			return viewLedgerAll(db, ledger)

		case "type":
			return viewLedgerType(db, ledger)

		case "account":
			return viewLedgerAccount(db, ledger)

		case "trial":
			return viewLedgerTrial(db, ledger)

		default:
			return fmt.Errorf("Invalid argument: %s\n", args[1])
		}

	default: 
		return errors.New("Usage: ledger view ledger [arg]")
	}
}

func viewChart(db *sql.DB, ledger *ledger.Ledger, args []string) error {
	switch len(args) {
	case 1:
		return viewChartAll(db, ledger)

	case 2:
		switch args[1] {
		case "help":
			viewChartHelp()
			return nil

		case "type":
			return viewChartType(db, ledger)

		default:
			return fmt.Errorf("Invalid argument: %s\n", args[1])
		}

	default:
		return errors.New("Usage: ledger view chart [arg]")
	}
}

func viewTypes(db *sql.DB, ledger *ledger.Ledger, args []string) error {
	switch len(args) {
	case 1:
		accounts, err := database.GetAccounts(db)
		if accounts == nil {
			return errors.New("No accounts available")
		}
		if err != nil {
			return err
		}
		err = ledger.SetAccounts(accounts)
		if err != nil {
			return err
		}
		fmt.Println(ledger.ViewAccountTypes())
		return nil

	default:
		return errors.New("Usage: ledger view types")
	}
}

func viewType(db *sql.DB, ledger *ledger.Ledger, args []string) error {
	switch len(args) {
	case 1:
		refPrefix, err := input.InputAccountTypeRefPrefix(ledger)
		if err != nil {
			return err
		}
		accountType, err := ledger.GetChart().GetAccountTypeByRefPrefix(refPrefix)
		if err != nil {
			return err
		}
		accounts, err := database.GetTypeAccounts(db, accountType.GetId())
		if err != nil {
			return err
		}
		err = ledger.SetAccounts(accounts)
		if err != nil {
			return err
		}
		output, err := ledger.ViewAccountType(refPrefix)
		if err != nil {
			return err
		}
		fmt.Println(output)
		return nil

	default:
		return errors.New("Usage: ledger view type")
	}
}

func viewAccounts(db *sql.DB, ledger *ledger.Ledger, args []string) error {
	switch len(args) {
	case 1:
		accounts, err := database.GetAccounts(db)
		if err != nil {
			return err
		}
		err = ledger.SetAccounts(accounts)
		if err != nil {
			return err
		}
		output, err := ledger.ViewAccounts()
		if err != nil {
			return err
		}
		fmt.Println(output)
		return nil

	default:
		return errors.New("Usage: ledger view accounts")
	}
}

func viewAccount(db *sql.DB, ledger *ledger.Ledger, args []string) error {
	switch len(args) {
	case 1:
		accounts, err := database.GetAccounts(db)
		if err != nil {
			return nil
		}
		err = ledger.SetAccounts(accounts)
		ref, err := input.InputAccountRef(ledger, "")
		if err != nil {
			return err
		}
		output, err := ledger.ViewAccountByRef(ref)
		if err != nil {
			return err
		}
		fmt.Println(output)
		return nil

	default:
		return errors.New("Usage: ledger view account")
	}
}

func viewJournal(db *sql.DB, ledger *ledger.Ledger, args []string) error {
	accounts, err := database.GetAccounts(db)
	if err != nil {
		return err
	}
	err = ledger.SetAccounts(accounts)
	if err != nil {
		return err
	}

	switch len(args) {
	case 1:
		return viewJournalInfo(db, ledger)

	case 2:
		switch args[1] {
		case "help":
			viewJournalHelp()
			return nil

		case "all":
			return viewJournalAll(db, ledger)

		case "dates":
			return viewJournalBetweenDates(db, ledger)

		case "posted":
			return viewJournalPosted(db, ledger, true)

		case "not-posted":
			return viewJournalPosted(db, ledger, false)

		case "entry":
			return viewJournalEntry(db, ledger)

		default:
			return fmt.Errorf("Invalid argument: %s\n", args[1])
		}

	default:
		return errors.New("Usage: ledger view journal [arg]")
	}
}

func viewLedgerInfo(db *sql.DB, ledger *ledger.Ledger) error {
	entries, err := database.GetEntries(db, ledger)
	if err != nil {
		return err
	}
	ledger.SetJournalEntries(entries)
	fmt.Println(ledger.String())
	return nil
}

func viewLedgerHelp() {
	var msg strings.Builder
	msg.WriteString("Argument options:\n")
	msg.WriteString("─────────────────\n")
	msg.WriteString("<empty> => view ledger information\n")
	msg.WriteString("help    => view this help\n")
	msg.WriteString("all     => view ledger for all accounts\n")
	msg.WriteString("type    => view ledger for an account type\n")
	msg.WriteString("account => view ledger for a single account\n")
	msg.WriteString("trial   => view trial balance\n")
	fmt.Println(msg.String())
}

func viewLedgerAll(db *sql.DB, ledger *ledger.Ledger) error {
	fmt.Println("To view ledger from start, press \"return\"")
	fromDate, err := input.InputDate("From")
	if err != nil {
		return err
	}
	fmt.Println("To view ledger until end, press \"return\"")
	toDate, err := input.InputDate("To")
	if err != nil {
		return err
	}
	entries, err := database.GetEntriesBetweenDates(db, ledger, fromDate, toDate)
	if err != nil {
		return err
	}
	ledger.SetJournalEntries(entries)
	output, err := ledger.ViewLedger()
	if err != nil {
		return err
	}
	fmt.Println(output)
	return nil
}

func viewLedgerType(db *sql.DB, ledger *ledger.Ledger) error {
	// TODO: implement view ledger for account type
	fmt.Println("View ledger for account type - TBI")
	return nil
}

func viewLedgerAccount(db *sql.DB, ledger *ledger.Ledger) error {
	// TODO: implement view ledger for account
	fmt.Println("View ledger for account - TBI")
	return nil
}

func viewLedgerTrial(db *sql.DB, ledger *ledger.Ledger) error {
	// TODO: implement view trial balance
	fmt.Println("View trial balance - TBI")
	return nil
}

func viewChartAll(db *sql.DB, ledger *ledger.Ledger) error {
	accounts, err := database.GetAccounts(db)
	if err != nil {
		return err
	}
	err = ledger.SetAccounts(accounts)
	if err != nil {
		return err
	}
	output, err := ledger.ViewChart(0)
	if err != nil {
		return err
	}
	fmt.Println(output)
	return nil
}

func viewChartHelp() {
	var msg strings.Builder
	msg.WriteString("View chart of acccounts\n")
	msg.WriteString("Argument options:\n")
	msg.WriteString("<empty> => view the complete chart of accounts\n")
	msg.WriteString("help    => show this help\n")
	msg.WriteString("type    => view the chart for a specific account type\n")
	fmt.Println(msg.String())
}

func viewChartType(db *sql.DB, ledger *ledger.Ledger) error {
	refPrefix, err := input.InputAccountTypeRefPrefix(ledger)
	if err != nil {
		return err
	}
	accountType, err := ledger.GetChart().GetAccountTypeByRefPrefix(refPrefix)
	if err != nil {
		return err
	}
	id := accountType.GetId()
	accounts, err := database.GetTypeAccounts(db, id)
	if err != nil {
		return err
	}
	err = ledger.SetAccounts(accounts)
	if err != nil {
		return err
	}
	output, err := ledger.ViewChart(refPrefix)
	if err != nil {
		return err
	}
	fmt.Println(output)
	return nil
}

func viewJournalInfo(db *sql.DB, ledger *ledger.Ledger) error {
	entries, err := database.GetEntries(db, ledger)
	if err != nil {
		return err
	}
	ledger.SetJournalEntries(entries)
	fmt.Println(ledger.ViewJournalInfo())
	return nil
}
func viewJournalHelp() {
	var msg strings.Builder
	msg.WriteString("Argument options:\n")
	msg.WriteString("<empty>      => view journal info\n")
	msg.WriteString("help         => view this help\n")
	msg.WriteString("all          => view all entries\n")
	msg.WriteString("dates        => view entries from date to date\n")
	msg.WriteString("posted       => view all entries posted\n")
	msg.WriteString("not-posted   => view all entries not posted\n")
	msg.WriteString("entry        => view an entry information\n")
	fmt.Println(msg.String())
}

func viewJournalAll(db *sql.DB, ledger *ledger.Ledger) error {
	entries, err := database.GetEntries(db, ledger)
	if err != nil {
		return err
	}
	ledger.SetJournalEntries(entries)
	fmt.Println(ledger.ViewJournal())
	return nil
}

func viewJournalBetweenDates(db *sql.DB, ledger *ledger.Ledger) error {
	fmt.Println("To view entries from start of journal, press \"return\"")
	fromDate, err := input.InputDate("From")
	if err != nil {
		return nil
	}
	fmt.Println("To view entries until the end of journal, press \"return\"")
	toDate, err := input.InputDate("To")
	if err != nil {
		return nil
	}
	entries, err := database.GetEntriesBetweenDates(db, ledger, fromDate, toDate)
	if err != nil {
		return err
	}
	ledger.SetJournalEntries(entries)
	fmt.Println(ledger.ViewJournal())
	return nil
}

func viewJournalPosted(db *sql.DB, ledger *ledger.Ledger, arePosted bool) error {
	entries, err := database.GetEntriesPosted(db, ledger, arePosted)
	if err != nil {
		return err
	}
	ledger.SetJournalEntries(entries)
	fmt.Println(ledger.ViewJournal())
	return nil
}

func viewJournalEntry(db *sql.DB, ledger *ledger.Ledger) error {
	fromDate, toDate, err := input.InputEntryYearMonth()
	if err != nil {
		return err
	}
	year := fromDate.Year()
	month := int(fromDate.Month())
	entries, err := database.GetEntriesBetweenDates(db, ledger, fromDate, toDate)
	if err != nil {
		return err
	}
	entry, err := input.InputEntryChoice(entries, year, month)
	if err != nil {
		return err
	}
	fmt.Println(entry.String())
	return nil
}