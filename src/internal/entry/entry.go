package entry

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"go.mod/internal/account"
	"go.mod/internal/currency"
	"go.mod/internal/id"
)

var MaxExplanationLength = 36

type Entry struct {
	id id.Id
	date time.Time
	debitAccount *account.Account
	creditAccount *account.Account
	amount currency.Currency
	explanation string
	isPosted bool
}

func NewEntry(
	date time.Time,
	debitAccount *account.Account,
	creditAccount *account.Account,
	amount currency.Currency,
	explanation string,
) *Entry {
	return &Entry{
		id: id.GenerateNewId(),
		date: date,
		debitAccount: debitAccount,
		creditAccount: creditAccount,
		amount: amount,
		explanation: explanation,
		isPosted: false,
	}
}

func NewDbEntry(
	sId string,
	sDate string,
	debitAccount *account.Account,
	creditAccount *account.Account,
	cents int,
	explanation string,
	intIsPosted int,
) (*Entry, error) {
	newId, err := id.ParseString(sId)
	if err != nil {
		return nil, err
	}

	date, err := time.Parse(time.DateTime, sDate)
	if err != nil {
		return nil, err
	}

	amount := currency.Currency(cents)

	isPosted := intIsPosted == 1

	return &Entry{
		id: newId,
		date: date,
		debitAccount: debitAccount,
		creditAccount: creditAccount,
		amount: amount,
		explanation: explanation,
		isPosted: isPosted,
	}, nil
}

func (e *Entry) GetEntry() *Entry {
	return e
}

func (e *Entry) GetId() id.Id {
	return e.id
}

func (e *Entry) GetDate() time.Time {
	return e.date
}

func (e *Entry) GetDebitAccount() *account.Account {
	return e.debitAccount
}

func (e *Entry) GetCreditAccount() *account.Account {
	return e.creditAccount
}

func (e *Entry) GetAmount() currency.Currency {
	return e.amount
}

func (e *Entry) GetExplanation() string {
	return e.explanation
}

func (e *Entry) IsPosted() bool {
	return e.isPosted
}

func (e *Entry) UpdateEntry(entry Entry) (*Entry, error) {
	if e.isPosted {
		return nil, errors.New("journal entry posted.")
	}

	e.date = entry.date
	e.debitAccount = entry.debitAccount
	e.creditAccount = entry.creditAccount
	e.amount = entry.amount
	e.explanation = entry.explanation

	return e, nil
}

func (e *Entry) Post() error {
	if e.isPosted {
		return errors.New("journal entry already posted.")
	}

	e.isPosted = true

	return nil
}

func (e *Entry) String() string {
	var output string
	output += "          Journal Entry Details\n"
	output += strings.Repeat("─", 12)
	output += "─┬─"
	output += strings.Repeat("─", 36)
	output += "\n"
	output += fmt.Sprintf("          Id │ %s\n", e.id.String())	
	output += fmt.Sprintf("        Date │ %s\n", e.date.Format(time.DateTime))
	output += fmt.Sprintf("      Amount │ %s\n", e.amount.String())
	output += fmt.Sprintf("       Debit │ %s", e.debitAccount.GetName())
	if e.isPosted {
		output += fmt.Sprintf(" (%d)\n", e.debitAccount.GetRef())
	} else {
		output += "\n"
	}
	output += fmt.Sprintf("      Credit │ %s", e.creditAccount.GetName())
	if e.isPosted {
		output += fmt.Sprintf(" (%d)\n", e.creditAccount.GetRef())
	} else {
		output += "\n"
	}
	output += fmt.Sprintf(" Explanation │ %s\n", e.explanation)
	output += "\n"

	return output
}

