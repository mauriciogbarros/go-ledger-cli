package chart

import (
	"errors"
	"fmt"
	"sort"
	"strings"

	"go.mod/internal/account"
	"go.mod/internal/accountType"
	"go.mod/internal/id"
)

type Chart struct {
	name string
	accountTypes *map[id.Id]*accountType.AccountType
}

func NewChart() *Chart {
	accountTypes := make(map[id.Id]*accountType.AccountType, 0)
	return &Chart{
		name: "Chart of Accounts",
		accountTypes: &accountTypes,
	}
}

func (c *Chart) GetChart() *Chart {
	return c
}

func (c *Chart) GetName() string {
	return c.name
}

func (c *Chart) GetAccountTypes() *[]*accountType.AccountType {
	accountTypes := make([]*accountType.AccountType, 0)
	for _, at := range *c.accountTypes {
		accountTypes = append(accountTypes, at)
	}

	sort.Slice(accountTypes, func(i, j int) bool {
		return accountTypes[i].GetRefPrefix() < accountTypes[j].GetRefPrefix()
	})

	return &accountTypes
}

func (c *Chart) GetAccountTypeById(id id.Id) (*accountType.AccountType, error) {
	at, exists := (*c.accountTypes)[id]
	if exists {
		return at, nil
	}

	return nil, errors.New("account type id not found.")
}

func (c *Chart) GetAccountTypeByName(name string) (*accountType.AccountType, error) {
	for _, at := range *c.accountTypes {
		if strings.Compare(at.GetName(), name) == 0 {
			return at, nil
		}
	}

	return nil, errors.New("account type name not found.")
}

func (c *Chart) GetAccountTypeByRefPrefix(refPrefix int) (*accountType.AccountType, error) {
	for _, at := range *c.accountTypes {
		if at.GetRefPrefix() == refPrefix {
			return at, nil
		}
	}

	return nil, errors.New("account type ref counter not found.")
}

func (c *Chart) GetAccountByStringId(stringId string) (*account.Account, error) {
	aId, err := id.ParseString(stringId)
	if err != nil {
		return nil, err
	}

	for _, at := range *c.accountTypes {
		account, err := at.GetAccountById(aId)
		if err == nil {
			return account, nil
		}
	}

	return nil, errors.New("account id not found.")
}

func (c *Chart) GetAccounts() *[]*account.Account {
	accounts := make([]*account.Account, 0)
	for _, at := range *c.accountTypes {
		accounts = append(accounts, *at.GetAccounts()...)
	}

	return &accounts
}

func (c *Chart) SetAccountTypes(accountTypes *map[id.Id]*accountType.AccountType) {
	c.accountTypes = accountTypes
}

func (c *Chart) MapAccountsToTypes(accounts *[]*account.Account) error {
	if accounts == nil {
		return errors.New("No accounts")
	}
	for _, at := range *c.accountTypes {
		clear(*at.GetAccountsMap())
	}

	for _, a := range *accounts {
		accountRefPrefix := a.GetRef() / 1000
		for _, at := range *c.accountTypes {
			if accountRefPrefix == at.GetRefPrefix() {
				ref, err := at.AddAccount(a)
				if err != nil {
					return err
				}
				a.SetRef(ref)
			}
		}
	}

	return nil
}

func (c *Chart) CreateAccount(name string, accountTypeRefPrefix int) (*account.Account, error) {
	account := account.NewAccount(name)
	accountType, err := c.GetAccountTypeByRefPrefix(accountTypeRefPrefix)
	if err != nil {
		return nil, err
	}
	ref, err := accountType.AddAccount(account)
	if err != nil {
		return nil, err
	}
	account.SetRef(ref)
	return account, nil
}

func (c *Chart) String(refPrefix int) (string, error) {
	width := 1 + 4 + 3 + 38 + 1
	var output strings.Builder
	output.WriteString(strings.Repeat(" ", (width - len(c.name))/ 2))
	output.WriteString(fmt.Sprintf("%s\n", c.name))
	output.WriteString(strings.Repeat("─", width))
	output.WriteString("\n")
	if refPrefix == 0 {
		for _, at := range *c.GetAccountTypes() {
			output.WriteString(" ")
			output.WriteString(fmt.Sprintf("%-*d", 4, at.GetRefPrefix()))
			output.WriteString(strings.Repeat(" ", (width - len(at.GetName()) - 1 - 4) / 2))
			output.WriteString(at.GetName())
			output.WriteString("\n")
			output.WriteString(strings.Repeat("─", width))
			output.WriteString("\n")
			accounts := *at.GetAccounts()
			if len(accounts) == 0 {
				output.WriteString(strings.Repeat(" ", 1 + 4 + 3))
				output.WriteString("* No accounts\n")
			} else {
				for _, a := range accounts {
					output.WriteString(fmt.Sprintf(" %d | %s\n", a.GetRef(), a.GetName()))
				}
			}
			output.WriteString("\n")
		}
	} else {
		at, err := c.GetAccountTypeByRefPrefix(refPrefix)
		if err != nil {
			return "", err
		}
		output.WriteString(" ")
		output.WriteString(fmt.Sprintf("%-*d", 4, at.GetRefPrefix()))
		output.WriteString(strings.Repeat(" ", (width - len(at.GetName()) - 1 - 4) / 2))
		output.WriteString(at.GetName())
		output.WriteString("\n")
		output.WriteString(strings.Repeat("─", width))
		output.WriteString("\n")
		accounts := *at.GetAccounts()
		if len(accounts) == 0 {
			output.WriteString(strings.Repeat(" ", 1 + 4 + 3))
			output.WriteString("* No accounts\n")
		} else {
			for _, a := range accounts {
				output.WriteString(fmt.Sprintf(" %d | %s\n", a.GetRef(), a.GetName()))
			}
		}
	}
	output.WriteString("\n")

	return output.String(), nil
}

func (c *Chart) ViewAccountTypes() string {
	width := 1 + 10 + 3 + accountType.MaxNameLength + 3 + 8 + 1
	title := "Account Types"
	paddingTitle := (width - len(title)) / 2
	
	var output strings.Builder
	output.WriteString("\n")
	output.WriteString(strings.Repeat(" ", paddingTitle))
	output.WriteString(title)
	output.WriteString("\n")
	output.WriteString(strings.Repeat("─", 1 + 10))
	output.WriteString("─┬─")
	output.WriteString(strings.Repeat("─", accountType.MaxNameLength))
	output.WriteString("─┬─")
	output.WriteString(strings.Repeat("─", 8 + 1))
	output.WriteString("\n")
	output.WriteString(" Ref Prefix │ ")
	output.WriteString(fmt.Sprintf("%-*s", accountType.MaxNameLength, "Name"))
	output.WriteString(" │ Accounts\n")
	output.WriteString(strings.Repeat("─", 1 + 10))
	output.WriteString("─┼─")
	output.WriteString(strings.Repeat("─", accountType.MaxNameLength))
	output.WriteString("─┼─")
	output.WriteString(strings.Repeat("─", 8 + 1))
	output.WriteString("\n")
	for _, at := range *c.GetAccountTypes() {
		output.WriteString(fmt.Sprintf(" %*d", 10, at.GetRefPrefix()))
		output.WriteString(" │ ")
		output.WriteString(fmt.Sprintf("%-*s", accountType.MaxNameLength, at.GetName()))
		output.WriteString(" │ ")
		output.WriteString(fmt.Sprintf("%*d", 8, len(*at.GetAccountsMap())))
		output.WriteString("\n")
	}

	return output.String()
}

func (c *Chart) ViewAccountType(refPrefix int) (string, error) {
	accountType, err := c.GetAccountTypeByRefPrefix(refPrefix)
	if err != nil {
		return "", err
	}

	return accountType.String(), nil
}