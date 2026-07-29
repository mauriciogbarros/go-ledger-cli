package account

import (
	"fmt"
	"strings"
	"time"

	"go.mod/internal/currency"
	"go.mod/internal/id"
)

type AccountEntry struct {
	id id.Id
	date time.Time
	explanation string
	amount currency.Currency
	side bool
	balance currency.Currency
}

func NewAccountEntry(
	id id.Id,
	date time.Time,
	explanation string,
	amount currency.Currency,
	side bool,
) *AccountEntry {
	return &AccountEntry{
		id: id,
		date: date,
		explanation: explanation,
		amount: amount,
		side: side, 
		balance: currency.Currency(0),
	}
}

func NewDbAccountEntry(
	sId string,
	sDate string,
	explanation string,
	cents int,
	intSide int,
) (*AccountEntry, error) {
	aeId, err := id.ParseString(sId)
	if err != nil {
		return nil, err
	}

	date, err := time.Parse(time.DateOnly, sDate)
	if err != nil {
		return nil, err
	}

	return &AccountEntry{
		id: aeId,
		date: date,
		explanation: explanation,
		amount: currency.Currency(cents),
		side: intSide == 1,
		balance: currency.Currency(0),
	}, nil
}

func (ae *AccountEntry) GetAccountEntry() *AccountEntry {
	return ae
}

func (ae *AccountEntry) GetId() id.Id {
	return ae.id
}

func (ae *AccountEntry) GetDate() time.Time {
	return ae.date
}

func (ae *AccountEntry) GetExplanation() string {
	return ae.explanation
}

func (ae *AccountEntry) GetAmount() currency.Currency {
	return ae.amount
}

func (ae *AccountEntry) GetSide() bool {
	return ae.side
}

func (ae *AccountEntry) SetBalance(balance currency.Currency) {
	ae.balance = balance
}

func (ae *AccountEntry) String() string {
	var output string
	output += "               Account Entry Details\n"
	output += strings.Repeat("─", 12)
	output += "─┬─"
	output += strings.Repeat("─", 36)
	output += "\n"
	output += fmt.Sprintf("          Id │ %s\n", ae.id.String())
	output += fmt.Sprintf("        Date │ %s\n", ae.date.Format(time.DateOnly))
	output += fmt.Sprintf(" Explanation │ %s\n", ae.explanation)
	output += "        Side │ "
	if ae.side {
		output += "Credit\n"
	} else {
		output += "Debit\n"
	}
	output += fmt.Sprintf("      Amount │ %s\n", ae.amount.String())
	if ae.balance > 0 {
		output += fmt.Sprintf("   Balance │ %s\n", ae.balance.String())
	}
	output += "\n"

	return output
}