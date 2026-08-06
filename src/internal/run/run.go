package run

import (
	"database/sql"
	"errors"
	"fmt"
	"os"
	"strings"

	"go.mod/internal/database"
	"go.mod/internal/ledger"
	"go.mod/internal/run/functions"
)

func Run() error {
	var msg strings.Builder
	var err error

	if len(os.Args) == 1 {
		msg.WriteString("Usage: ledger <command> [args]\n")
		msg.WriteString("Enter \"ledger help\" to view commands.")
		return errors.New(msg.String())
	}

	function := os.Args[1]
	if function == "help" {		
		help()
		return nil
	}

	var db *sql.DB
	db, err = sql.Open("sqlite", "./internal/database/ledger.db")
	if err != nil {
		return err
	}
	defer db.Close()
	err = database.Initialize(db)
	if err != nil {
		return err
	}
	ledger := ledger.NewLedger()
	ledger.CreateChart()
	ledger.CreateJournal()
	accountTypes, err := database.GetAccountTypes(db)
	if err != nil {
		return err
	}
	ledger.SetChartAccountTypes(accountTypes)

	var args []string
	if len(os.Args) > 2 {
		args = os.Args[2:]
	}

	switch function {
	case "view":
		return functions.View(db, ledger, args)

	case "new":
		return functions.New(db, ledger, args)

	case "update":
		return functions.Update(db, ledger, args)

	case "delete":
		return functions.Delete(db, ledger, args)

	default:
		return fmt.Errorf("Unknown function: %s\n", function)
	}
}

func help() {
		var msg strings.Builder
		msg.WriteString("Usage: ledger <function> <command> [arg]\n")
		msg.WriteString("Available functions:\n")
		msg.WriteString("────────────────────\n")
		msg.WriteString("help                 => Show this help\n")
		msg.WriteString("view <command> [arg] => Viewing functions\n")
		msg.WriteString("new <command>        => New functions\n")
		msg.WriteString("update <command>     => Updating functions\n")
		msg.WriteString("delete <command>     => Deleting functions\n")
		fmt.Println(msg.String())
}