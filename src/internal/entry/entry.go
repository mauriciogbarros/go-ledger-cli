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
	debitAccountRef int
	creditAccountRef int
	amount currency.Currency
	explanation string
}

func NewEntry(
	debitAccountRef int,
	creditAccountRef int,
	amount currency.Currency,
	explanation string,
) Entry {
	return Entry{
		id: id.GenerateNewId(),
		date: time.Now(),
		debitAccountRef: debitAccountRef,
		creditAccountRef: creditAccountRef,
		amount: amount,
		explanation: explanation,
	}
}

func NewEntryFromDb(
	id id.Id,
	date time.Time,
	debitAccountRef int,
	creditAccountRef int,
	amount currency.Currency,
	explanation string,
) Entry {
	return Entry{
		id: id,
		date: date,
		debitAccountRef: debitAccountRef,
		creditAccountRef: creditAccountRef,
		amount: amount,
		explanation: explanation,
	}
}

func (e *Entry) GetId() id.Id {
	return e.id
}

func (e *Entry) GetDate() time.Time {
	return e.date
}

func (e *Entry) GetDebitAccountRef() int {
	return e.debitAccountRef
}

func (e *Entry) GetCreditAccountRef() int {
	return e.creditAccountRef
}

func (e *Entry) GetAmount() currency.Currency {
	return e.amount
}

func (e *Entry) GetExplanation() string {
	return e.explanation
}

func (e Entry) Format(debitAccountName, creditAccountName string) string {
	var entry string = " "
	entry += fmt.Sprintf("%-*s", 19, e.date.Format(time.DateTime))
	entry += " │ "
	entry += fmt.Sprintf("%-*s", account.MaxNameLength + 2, debitAccountName)
	entry += " │ "
	entry += strings.Repeat(" ", 3)
	entry += " │ "
	entry += fmt.Sprintf("%*s", 12, e.amount.String())
	entry += " │\n"
	entry += " "
	entry += strings.Repeat(" ", 19)
	entry += " │   "
	entry += fmt.Sprintf("%-*s", account.MaxNameLength, creditAccountName)
	entry += " │ "
	entry += strings.Repeat(" ", 3)
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
	entry += strings.Repeat(" ", 3)
	entry += " │ "
	entry += strings.Repeat(" ", 12)
	entry += " │ "
	entry += strings.Repeat(" ", 12)
	entry += "\n"

	return entry
}