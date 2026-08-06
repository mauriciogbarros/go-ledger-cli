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

func New(
	db *sql.DB,
	ledger *ledger.Ledger,
	args []string,
) error {
	if len(args) == 0 {
		return errors.New("Usage: ledger new <command>")
	}

	switch args[0] {
	case "help":
		newHelp()
		return nil

	case "account":
		return newAccount(db, ledger, args)

	case "entry":
		return newEntry(db, ledger, args)

	default:
		return fmt.Errorf("Invalid command: %s\n", args[0])
	}
}

func newHelp() {
	var msg strings.Builder
	msg.WriteString("Usage: ledger new <command>\n")
	msg.WriteString("Commands:\n")
	msg.WriteString("─────────\n")
	msg.WriteString("account => Create a new account\n")
	msg.WriteString("entry   => Create a new journal entry\n")
	fmt.Println(msg.String())
}

func newAccount(db *sql.DB, ledger *ledger.Ledger, args []string) error {
	if len(args) > 1 {
		return errors.New("Usage: ledger new account")
	}

	accounts, err := database.GetAccounts(db)
	if err != nil {
		return err
	}
	err = ledger.SetAccounts(accounts)
	if err != nil {
		return err
	}
	name, err := input.InputAccountName()
	if err != nil {
		return err
	}
	atRefPrefix, err := input.InputAccountTypeRefPrefix(ledger)
	if err != nil {
		return err
	}
	description, err := input.InputText("Description")
	if err != nil {
		return err
	}
	newAccount, err := ledger.CreateAccount(atRefPrefix, name, description)
	if err != nil {
		return err
	}
	err = database.AddAccount(db, newAccount)
	if err != nil {
		return err
	}
	fmt.Printf("Account created: %s (%d)\n", newAccount.GetName(), newAccount.GetRef())
	return nil
}

func newEntry(db *sql.DB, ledger *ledger.Ledger, args []string) error {
	accounts, err := database.GetAccounts(db)
	if err != nil {
		return err
	}
	err = ledger.SetAccounts(accounts)
	if err != nil {
		return err
	}
	amountF64, err := input.InputAmountF64()
	if err != nil {
		return err
	}
	date, err := input.InputDate("Entry")
	if err != nil {
		return err
	}
	debitAccountRef, err := input.InputAccountRef(ledger, "debit")
	if err != nil {
		return err
	}
	creditAccountRef, err := input.InputAccountRef(ledger, "credit")
	if err != nil {
		return err
	}
	if debitAccountRef == creditAccountRef {
		return errors.New("Debit and Credit accounts must be different")
	}
	explanation, err := input.InputText("Explanation")
	if err != nil {
		return err
	}
	newEntry, err := ledger.CreateJournalEntry(date, debitAccountRef, creditAccountRef, amountF64, explanation)
	if err != nil {
		return err
	}
	err = database.AddEntry(db, &newEntry)
	if err != nil {
		return err
	}
	debitAccountName := newEntry.GetDebitAccount().GetName()
	creditAccountName := newEntry.GetCreditAccount().GetName()
	fmt.Printf("Entry created: Dr = %s; Cr = %s; $%s", debitAccountName, creditAccountName, newEntry.GetAmount().String())
	return nil
}