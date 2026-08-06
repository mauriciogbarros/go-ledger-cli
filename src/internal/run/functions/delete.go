package functions

import (
	"database/sql"
	"errors"
	"fmt"

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

	case "entry":
		return deleteEntry(db, ledger, args)

	default:
		return fmt.Errorf("Invalid command: %s\n", args[0])
	}
}

func deleteHelp() {
	// TODO: Delete help
	fmt.Println("Delete help - TBI")
}

func deleteEntry(db *sql.DB, ledger *ledger.Ledger, args []string) error {
	// TODO: Delete entry
	fmt.Println("Delete entry - TBI")
	return nil
}