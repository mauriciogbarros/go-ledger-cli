package entry

import (
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"go.mod/internal/account"
	"go.mod/internal/currency"
	"go.mod/internal/id"
)

var MaxExplanationLength = 36

type Entry struct {
	id id.Id
	date time.Time
	debitAccount string
	debitRef id.Id
	debitAmount currency.Currency
	creditAccount string
	creditRef id.Id
	creditAmount currency.Currency
	explanation string
}

func NewEntry(
	debitAccount string,
	creditAccount string,
	amount float64,
	explanation string,
) Entry {
	return Entry{
		id: id.NewId(),
		date: time.Now(),
		debitAccount: debitAccount,
		debitRef: id.Id(uuid.Nil),
		debitAmount: currency.Convert64(amount),
		creditAccount: creditAccount,
		creditRef: id.Id(uuid.Nil),
		creditAmount: currency.Convert64(amount),
		explanation: explanation,
	}
}

func (e Entry) String() string {
	var entry string
	entry += " "
	entry += fmt.Sprintf("%-*s", 19, e.date.Format(time.DateTime))
	entry += " │ "
	entry += fmt.Sprintf("%-*s", account.MaxNameLength + 2, e.debitAccount)
	entry += " │ "
	entry += strings.Repeat(" ", 5)
	entry += " │ "
	entry += fmt.Sprintf("%*s", 12, e.debitAmount.String())
	entry += " │\n"
	entry += " "
	entry += strings.Repeat(" ", 19)
	entry += " │   "
	entry += fmt.Sprintf("%-*s", account.MaxNameLength, e.creditAccount)
	entry += " │ "
	entry += strings.Repeat(" ", 5)
	entry += " │ "
	entry += strings.Repeat(" ", 12)
	entry += " │ "
	entry += fmt.Sprintf("%*s", 12, e.creditAmount.String())
	entry += "\n"
	entry += " "
	entry += strings.Repeat(" ", 19)
	entry += " │     "
	entry += fmt.Sprintf("%-*s", MaxExplanationLength, e.explanation)
	entry += " │ "
	entry += strings.Repeat(" ", 5)
	entry += " │ "
	entry += strings.Repeat(" ", 12)
	entry += " │ "
	entry += strings.Repeat(" ", 12)
	entry += "\n"

	return entry
}
