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

	date, err := time.Parse(time.DateOnly, sDate)
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
	var output strings.Builder
	output.WriteString(strings.Repeat(" ", 18))
	output.WriteString("Journal Entry Details\n")
	output.WriteString(strings.Repeat("─", 15))
	output.WriteString("─┬─")
	output.WriteString(strings.Repeat("─", account.MaxNameLength))
	output.WriteString("\n")
	fmt.Fprintf(&output, "             Id │ %s\n", e.id.String())
	fmt.Fprintf(&output, "           Date │ %s\n", e.date.Format(time.DateOnly))
	fmt.Fprintf(&output, "         Amount │ %s\n", e.amount.String())
	fmt.Fprintf(&output, "  Debit Account │ %s", e.debitAccount.GetName())
	if e.isPosted {
		fmt.Fprintf(&output, " (%d)\n", e.debitAccount.GetRef())
	} else {
		output.WriteString("\n")
	}
	fmt.Fprintf(&output, " Credit Account │ %s", e.creditAccount.GetName())
	if e.isPosted {
		fmt.Fprintf(&output, " (%d)\n", e.creditAccount.GetRef())
	} else {
		output.WriteString("\n")
	}

	output.WriteString("    Explanation │ ")
	words := strings.Split(e.GetExplanation(), " ")
	var explanation strings.Builder
	for i := 0; i < len(words); {
		for exLen := 0; exLen <= account.MaxNameLength && i < len(words); {
			explanation.WriteString(words[i])
			if i < len(words) {
				i++
			}
			if i < len(words) {
				explanation.WriteString(" ")
				exLen = explanation.Len() + len(words[i])
			}
		}
		output.WriteString(explanation.String())
		output.WriteString("\n")
		if i < len(words) {
			explanation.Reset()
			output.WriteString(strings.Repeat(" ", 15))
			output.WriteString(" │ ")
		}
	}
	output.WriteString("         Posted │ ")
	if e.isPosted {
		output.WriteString("Yes\n")
	} else {
		output.WriteString("No\n")
	}
	output.WriteString("\n")

	return output.String()
}

