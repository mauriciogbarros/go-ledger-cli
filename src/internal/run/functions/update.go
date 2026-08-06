package functions

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"go.mod/internal/database"
	"go.mod/internal/ledger"
)

func Update(
	db *sql.DB,
	ledger *ledger.Ledger,
	args []string,
) error {
	if len(args) == 0 {
		return errors.New("Usage: ledger update <command>")
	}

	switch args[0] {
	case "help":
		updateHelp()
		return nil

	case "account":
		return updateAccount(db, ledger, args)

	case "entry":
		return updateEntry(db, ledger, args)

	default:
		return fmt.Errorf("Invalid command: %s\n", args[0])
	}
}

func updateHelp() {
	var msg strings.Builder
	msg.WriteString("Usage: ledger update <command>\n")
	msg.WriteString("Commands:\n")
	msg.WriteString("─────────\n")
	msg.WriteString("account => Update account information\n")
	msg.WriteString("entry   => Update journal entry information\n")
	fmt.Println(msg.String())
}

func updateAccount(db *sql.DB, ledger *ledger.Ledger, args []string) error {
	// TODO: Update account
	fmt.Println("Update account - TBI")
	return nil
}

func updateEntry(db *sql.DB, ledger *ledger.Ledger, args []string) error {
	switch len(args) {
	case 1:
		entries, err := database.GetEntriesPosted(db, ledger, false)
		if err != nil {
			return err
		}
		if entries == nil {
			fmt.Println("All entries posted")
			return nil
		}
		
		var msg strings.Builder
		msg.WriteString("Usage: ledger update entry\n")
		msg.WriteString("Only accounts not-posted may be updated.")
		msg.WriteString("Argument options:\n")
		msg.WriteString("─────────────────\n")
		msg.WriteString("date => udpate the entry's date\n")
		msg.WriteString("debit => update the entry's debit account\n")
		msg.WriteString("credit => update the entry's credit account\n")
		msg.WriteString("explanation => update the entry's explanation\n")
		msg.WriteString("amount => update the entry's amount\n")
		msg.WriteString("post => change the entry's status")

	default:
		var msg strings.Builder
		msg.WriteString("Usage: ledger udpate entry\n")
		msg.WriteString("=> Only entries not-posted may be updated\n")
		return errors.New(msg.String())
	}
	return nil
}