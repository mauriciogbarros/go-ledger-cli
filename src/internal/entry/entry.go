package entry

import (
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
	posted bool
}

func NewEntry(
	debitAccount *account.Account,
	creditAccount *account.Account,
	amount currency.Currency,
	explanation string,
) Entry {
	return Entry{
		id: id.GenerateNewId(),
		date: time.Now(),
		debitAccount: debitAccount,
		creditAccount: creditAccount,
		amount: amount,
		explanation: explanation,
		posted: false,
	}
}

func NewEntryFromDb(
	sId string,
	sDate string,
	debitAccount *account.Account,
	creditAccount *account.Account,
	cents int,
	explanation string,
	intPosted int,
) (*Entry, error) {
	newId, err := id.ParseString(sId)
	if err != nil {
		return nil, err
	}

	date, err := time.Parse(time.RFC3339, sDate)
	if err != nil {
		return nil, err
	}

	amount := currency.Currency(cents)

	posted := false
	if intPosted == 1 {
		posted = true
	}

	return &Entry{
		id: newId,
		date: date,
		debitAccount: debitAccount,
		creditAccount: creditAccount,
		amount: amount,
		explanation: explanation,
		posted: posted,
	}, nil
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

func (e *Entry) GetPostedInt() int {
	if e.posted {
		return 1
	}
	return 0
}

func (e *Entry) IsPosted() bool {
	return e.posted
}

func (e *Entry) Post() {
	e.posted = true
}

func (e *Entry) String() string {
	var entry string = " "
	entry += fmt.Sprintf("%-*s", 19, e.date.Format(time.DateTime))
	entry += " │ "
	entry += fmt.Sprintf("%-*s", account.MaxNameLength + 2, e.debitAccount.GetName())
	entry += " │ "
	if e.posted {
		entry += fmt.Sprintf("%*d", 4, e.debitAccount.GetRef())
	} else {
		entry += strings.Repeat(" ", 4)
	}
	entry += " │ "
	entry += fmt.Sprintf("%*s", 12, e.amount.String())
	entry += " │\n"
	entry += " "
	entry += strings.Repeat(" ", 19)
	entry += " │   "
	entry += fmt.Sprintf("%-*s", account.MaxNameLength, e.creditAccount.GetName())
	entry += " │ "
	if e.posted {
		entry += fmt.Sprintf("%*d", 4, e.creditAccount.GetRef())
	} else {
		entry += strings.Repeat(" ", 4)
	}
	entry += " │ "
	entry += strings.Repeat(" ", 12)
	entry += " │ "
	entry += fmt.Sprintf("%*s", 12, e.amount.String())
	entry += "\n"
	entry += " "
	entry += strings.Repeat(" ", 19)
	entry += " │     "
	entry += fmt.Sprintf("%-*s", MaxExplanationLength, e.explanation)
	entry += " │ "
	entry += strings.Repeat(" ", 4)
	entry += " │ "
	entry += strings.Repeat(" ", 12)
	entry += " │ "
	entry += strings.Repeat(" ", 12)
	entry += "\n"

	return entry
}