package ledger

import (
	"errors"
	"fmt"
	"strings"

	"go.mod/internal/account"
	"go.mod/internal/accountType"
	"go.mod/internal/id"
)

type Ledger struct {
	name string
	accountTypes *[]accountType.AccountType
}

func NewLedger(name string) *Ledger {
	at := make([]accountType.AccountType, 0)
	return &Ledger{
		name: name,
		accountTypes: &at,
	}
}

func (l *Ledger) GetName() string {
	return l.name
}

func (l *Ledger) GetAccountTypes() *[]accountType.AccountType {
	return l.accountTypes
}

func (l *Ledger) GetAccountTypeById(id id.Id) (*accountType.AccountType, error) {
	for _, a := range *l.accountTypes {
		if a.GetId() == id {
			return &a, nil
		}
	}

	return nil, errors.New("account type id not found")
}

func (l *Ledger) GetAccountTypeByName(name string) (*accountType.AccountType, error) {
	for _, a := range *l.accountTypes {
		if a.GetName() == name {
			return &a, nil
		}
	}

	return nil, errors.New("account type name not found")
}

func (l *Ledger) GetAccountTypeByRef(ref int) (*accountType.AccountType, error) {
	r := ref / 1000
	for _, a := range *l.accountTypes {
		aR := a.GetRefCounter() / 1000
		if aR == r {
			return &a, nil
		}
	}

	return nil, errors.New("account type ref not found")
}

func (l *Ledger) AddAccountType(accountType accountType.AccountType) {
	*l.accountTypes = append(*l.accountTypes, accountType)
}

func (l *Ledger) SetAccountTypes(accountTypes *[]accountType.AccountType) {
	l.accountTypes = accountTypes
}

func (l *Ledger) MapAccounts(accounts *[]account.Account) {
	for _, a := range *accounts {
		for _, at := range *l.accountTypes {
			if a.GetAccountTypeId() == at.GetId() {
				at.AddAccount(&a)
				break
			}
		}
	}
}

func (l *Ledger) String() string {
	return "String representation to be implemented"
}

func (l *Ledger) PrintTypes() string {
	width := 1 + 9 + 3 + accountType.MaxNameLength + 3 + 6 + 1
	paddingLeft := (width - len("Account Types")) / 2
	var output string = "\n"
	output += strings.Repeat(" ", paddingLeft)
	output += "Account Types"
	output += "\n"
	output += strings.Repeat("─", width)
	output += "\n"
	output += " "
	output += "Ref Group"
	output += "   "
	output += fmt.Sprintf("%-*s", accountType.MaxNameLength, "Name")
	output += "   "
	output += fmt.Sprintf("%-*s", 6, "Side")
	output += "\n"
	output += "─"
	output += strings.Repeat("─", 9)
	output += "─┬─"
	output += strings.Repeat("─", accountType.MaxNameLength)
	output += "─┬─"
	output += strings.Repeat("─", 6)
	output += "─"
	output += "\n"
	for _, t := range *l.accountTypes {
		output += " "
		output += fmt.Sprintf("%*d", 9, t.GetRefGroup())
		output += " │ "
		output += fmt.Sprintf("%-*s", accountType.MaxNameLength, t.GetName())
		output += " │ "
		output += fmt.Sprintf("%-*s", 6, t.GetSide())
		output += "\n"
	}

	return output
}