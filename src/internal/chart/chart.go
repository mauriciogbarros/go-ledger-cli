package chart

import (
	"errors"
	"fmt"
	"strings"

	"go.mod/db"
	"go.mod/internal/account"
)

var refCounter int

type ChartOfAccounts struct {
	name string
	accounts []account.Account
}

func (c *ChartOfAccounts) NewChartOfAccounts(name string) error {
	c.name = name
	c.accounts = make([]account.Account, 0)
	accounts, err := db.GetAccounts()
	if err != nil {
		return err
	}
	c.accounts = accounts
	refCounter = len(accounts)
	return nil
}

func (c *ChartOfAccounts) GetName() string {
	return c.name
}

func (c *ChartOfAccounts) GetAccounts() []account.Account {
	return c.accounts
}

func (c *ChartOfAccounts) AddAccount(name string, accountType int) (account.Account, error) {
	for _, a := range c.accounts {
		if strings.Compare(a.GetName(), name) == 0 {
			return account.Account{}, errors.New("Account name must be unique.")
		}
	}
	refCounter++
	newAccount := account.NewAccount(refCounter, name, accountType)
	c.accounts = append(c.accounts, newAccount)
	if err := db.CreateAccount(newAccount); err != nil {
		return account.Account{}, err
	}

	return newAccount, nil
}

func (c ChartOfAccounts) String() string {
	width := 1 + 3 + 3 + account.MaxNameLength + 3 + 9 + 1
	paddingLeft := (width - len(c.name)) / 2
	var output string = "\n"
	output += strings.Repeat(" ", paddingLeft)
	output += c.name
	output += "\n"
	output += strings.Repeat("─", width)
	output += "\n"
	output += " Ref   "
	output += fmt.Sprintf("%-*s", account.MaxNameLength, "Name")
	output += "   "
	output += fmt.Sprintf("%-*s", 9, "Type")
	output += " \n"
	output += strings.Repeat("─", 4)
	output += "─┬─"
	output += strings.Repeat("─", account.MaxNameLength)
	output += "─┬─"
	output += strings.Repeat("─", 9)
	output += "─\n"
	if len(c.accounts) == 0 {
		output += " No acccounts\n"
	} else {
		for _, account := range c.accounts {
			output += account.String()
		}
	}

	return output
}

func (c *ChartOfAccounts) GetAccountByRef(ref int) (account.Account, error) {
	for _, a := range c.accounts {
		if a.GetRef() == ref {
			return a, nil
		}
	}
	return account.Account{}, errors.New("Account not found")
}