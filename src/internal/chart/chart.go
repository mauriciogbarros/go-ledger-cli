package chart

import (
	"errors"
	"fmt"
	"strings"

	"go.mod/internal/account"
	"go.mod/internal/accountType"
	"go.mod/internal/id"
	"go.mod/internal/ledger"
)

type ChartOfAccounts struct {
	name string
	ledger *ledger.Ledger
	accounts *[]account.Account
}

func NewChartOfAccounts(name string, ledger *ledger.Ledger) *ChartOfAccounts {
	accounts := make([]account.Account, 0)
	return &ChartOfAccounts{
		name: name,
		ledger: ledger,
		accounts: &accounts,
	}
}

func (c *ChartOfAccounts) GetName() string {
	return c.name
}

func (c *ChartOfAccounts) GetLedger() *ledger.Ledger {
	return c.ledger
}

func (c *ChartOfAccounts) GetAccounts() *[]account.Account {
	return c.accounts
}

func (c *ChartOfAccounts) GetAccountByName(name string) (*account.Account, error) {
	for _, a := range *c.accounts {
		if strings.Compare(a.GetName(), name) == 0 {
			return &a, nil
		}
	}

	return nil, errors.New("account name not found.")
}

func (c *ChartOfAccounts) GetAccountById(id id.Id) (*account.Account, error) {
	for _, a := range *c.accounts {
		if a.GetId() == id {
			return &a, nil
		}
	}

	return nil, errors.New("account id not found.")
}

func (c *ChartOfAccounts) GetAccountByStringId(strId string) (*account.Account, error) {
	aId, err := id.ParseString(strId)
	if err != nil {
		return nil, err
	}

	return c.GetAccountById(aId)
}

func (c *ChartOfAccounts) GetAccountByRef(ref int) (*account.Account, error) {
	for _, a := range *c.accounts {
		if a.GetRef() == ref {
			return &a, nil
		}
	}

	return nil, errors.New("account ref not found.")
}

func (c *ChartOfAccounts) AddAccount(account *account.Account) error {
	for _, a := range *c.accounts {
		if strings.Compare(a.GetName(), account.GetName()) == 0 {
			return errors.New("Account name must be unique.")
		}
	}

	*c.accounts = append(*c.accounts, *account)
	
	return nil
}

func (c *ChartOfAccounts) RemoveAccount(id id.Id) error {
	for i, a := range *c.accounts {
		if a.GetId() == id {
			*c.accounts = append((*c.accounts)[:i], (*c.accounts)[i+1:]...)
			return nil
		}
	}

	return errors.New("account id not found.")
}

func (c *ChartOfAccounts) SetAccounts(accounts *[]account.Account) {
	c.accounts = accounts
}
func (c *ChartOfAccounts) String() string {
	width := 1 + 4 + 3 + account.MaxNameLength + 3 + accountType.MaxNameLength + 3 + 6 + 1
	paddingLeft := (width - len(c.name)) / 2
	var output string = "\n"
	output += strings.Repeat(" ", paddingLeft)
	output += c.name
	output += "\n"
	output += strings.Repeat("─", width)
	output += "\n"
	output += " "
	output += fmt.Sprintf("%*s", 4, "Ref")
	output += "   "
	output += fmt.Sprintf("%-*s", account.MaxNameLength, "Name")
	output += "   "
	output += fmt.Sprintf("%-*s", accountType.MaxNameLength, "Type")
	output += "   "
	output += fmt.Sprintf("%-*s", 6, "Side")
	output += "\n"
	output += "─"
	output += strings.Repeat("─", 4)
	output += "─┬─"
	output += strings.Repeat("─", account.MaxNameLength)
	output += "─┬─"
	output += strings.Repeat("─", accountType.MaxNameLength)
	output += "─┬─"
	output += strings.Repeat("─", 6)
	output += "─\n"
	if len(*c.accounts) == 0 {
		output += "        *No acccounts\n"
	} else {
		for _, account := range *c.accounts {
			output += account.String()
			accountTypeId := account.GetAccountTypeId()
			accountType, err := c.ledger.GetAccountTypeById(accountTypeId)
			if err != nil {
				panic(err)
			}
			output += accountType.String()
			output += "\n"
		}
	}
	output += "\n"

	return output
}