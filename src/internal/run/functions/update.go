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

func Update(
	db *sql.DB,
	ledger *ledger.Ledger,
	args []string,
) error {
	if len(args) == 0 {
		return errors.New("Usage: ledger update <command>")
	}

	accounts, err := database.GetAccounts(db)
	if err != nil {
		return err
	}
	err = ledger.SetAccounts(accounts)
	if err != nil {
		return err
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
	switch len(args) {
	case 1:
		ref, err := input.InputAccountRef(ledger, "update")
		if err != nil {
			return err
		}
		account, err := ledger.GetChart().GetAccountByRef(ref)
		if err != nil {
			return nil
		}
		fmt.Println(account.String())
		choice, err := input.InputAccountFieldChoice()
		if err != nil {
			return err
		}
		switch choice {
		case 1:
			name, err := input.InputAccountName()
			if err != nil {
				return err
			}
			err = account.SetName(name)
			if err != nil {
				return err
			}

		case 2:
			description, err := input.InputText("Description")
			if err != nil {
				return err
			}
			account.SetDescription(description)

		default:
			return errors.New("Invalid field")
		}
		err = database.UpdateAccount(db, account)
		if err != nil {
			return err
		}
		fmt.Println("Account updated successfully")
		fmt.Println(account.String())
		return nil
	
	default:
		return errors.New("Usage: ledger update account")
	}
}

func updateEntry(db *sql.DB, ledger *ledger.Ledger, args []string) error {
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
			isPosted, err := input.InputEntryIsPostedStatus(entry.IsPosted())
			if err != nil {
				return err
			}
			entry.SetIsPosted(isPosted)
		} else {
		choice, err := input.InputEntryFieldChoice()
		if err != nil {
			return err
		}
			switch choice {
			case 1:
				date, err := input.InputDate("Entry")
				if err != nil {
					return err
				}
				entry.SetDate(date)

			case 2:
				amountF64, err := input.InputAmountF64()
				if err != nil {
					return err
				}
				entry.SetAmount(amountF64)

			case 3:
				debitAccountRef, err := input.InputAccountRef(ledger, "debit")
				if err != nil {
					return err
				}
				if debitAccountRef == entry.GetCreditAccount().GetRef() {
					return errors.New("Debit and Credit accounts must be different")
				}
				debitAccount, err := ledger.GetChart().GetAccountByRef(debitAccountRef)
				if err != nil {
					return err
				}
				entry.SetDebitAccount(debitAccount)

			case 4:
				creditAccountRef, err := input.InputAccountRef(ledger, "credit")
				if err != nil {
					return err
				}
				if entry.GetDebitAccount().GetRef() == creditAccountRef {
					return errors.New("Debit and Credit accounts must be different")
				}
				creditAccount, err := ledger.GetChart().GetAccountByRef(creditAccountRef)
				if err != nil {
					return err
				}
				entry.SetCreditAccount(creditAccount)

			case 5:
				explanation, err := input.InputText("Explanation")
				if err != nil {
					return err
				}
				entry.SetExplanation(explanation)

			case 6:
				isPosted, err := input.InputEntryIsPostedStatus(entry.IsPosted())
				if err != nil {
					return err
				}
				entry.SetIsPosted(isPosted)

			default:
				return errors.New("Invalid field")
			}
		}
		err = database.UpdateEntry(db, entry)
		if err != nil {
			return err
		}
		fmt.Println("Entry update successfuly")
		fmt.Println(entry.String())
		return nil

	default:
		return errors.New("Usage: ledger update entry")
	}
}