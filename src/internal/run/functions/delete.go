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

func Delete(
	db *sql.DB,
	ledger *ledger.Ledger,
	args []string,
) error {
	if len(args) == 0 {
		return errors.New("Usage: ledger view <command> [args]")
	}

	switch args[0] {
	case "help":
		deleteHelp()
		return nil

	case "account":
		return deleteAccount(db, ledger, args)

	case "entry":
		return deleteEntry(db, ledger, args)

	default:
		return fmt.Errorf("Invalid command: %s\n", args[0])
	}
}

func deleteHelp() {
	var msg strings.Builder
	msg.WriteString("Usage: ledger delete <command>\n")
	msg.WriteString("Commands:\n")
	msg.WriteString("─────────\n")
	msg.WriteString("account => delete an empty account\n")
	msg.WriteString("entry   => delete an entry not-posted\n")
	fmt.Println(msg.String())
}

func deleteAccount(db *sql.DB, ledger *ledger.Ledger, args []string) error {
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
		ref, err := input.InputAccountRef(ledger, "delete")
		if err != nil {
			return err
		}
		account, err := ledger.GetChart().GetAccountByRef(ref)
		if err != nil {
			return err
		}
		if len(*account.GetEntries()) > 0 {
			return errors.New("Cannot delete account with entries")
		}
		err = database.DeleteAccount(db, account.GetId())
		if err != nil {
			return err
		}
		fmt.Printf("Deleted account %s\n", account.GetName())
		return nil

	default:
		return errors.New("Usage: ledger delete account")
	}
}

func deleteEntry(db *sql.DB, ledger *ledger.Ledger, args []string) error {
	switch len(args) {
	case 1:
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
		if entry.IsPosted() {
			return errors.New("Cannot delete a posted entry")
		}
		err = database.DeleteEntry(db, entry.GetId())
		if err != nil {
			return err
		}
		fmt.Printf("Deleted entry %s\n", entry.GetExplanation())
		return nil
	default:
		return errors.New("Usage: ledger delete entry")
	}
}