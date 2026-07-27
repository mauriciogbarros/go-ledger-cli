package accountType

import (
	"errors"
	"fmt"
	"sort"
	"strings"

	"go.mod/internal/account"
	"go.mod/internal/id"
)

var MaxNameLength int = 20

type AccountType struct {
	id id.Id
	name string
	refPrefix int
	accounts *map[id.Id]*account.Account
}

func NewDbAccountType(
	sId string,
	name string,
	refPrefix int,
) (*AccountType, error) {
	atId, err := id.ParseString(sId)
	if err != nil {
		return nil, err
	}

	accounts := make(map[id.Id]*account.Account, 0)

	return &AccountType{
		id: atId,
		name: name,
		refPrefix: refPrefix,
		accounts: &accounts,
	}, nil
}

func (at *AccountType) GetAccountType() *AccountType {
	return at
}

func (at *AccountType) GetId() id.Id {
	return at.id
}

func (at *AccountType) GetName() string {
	return at.name
}

func (at *AccountType) GetRefPrefix() int {
	return at.refPrefix
}

func (at *AccountType) GetAccounts() *[]*account.Account {
	accounts := make([]*account.Account, 0)
	for _, account := range *at.accounts {
		accounts = append(accounts, account)
	}

	sort.Slice(accounts, func(i, j int) bool {
		return accounts[i].GetRef() < accounts[j].GetRef()
	})

	return &accounts
}

func (at *AccountType) GetAccountsMap() *map[id.Id]*account.Account {
	return at.accounts
}

func (at *AccountType) GetAccountById(id id.Id) (*account.Account, error) {
	account, exist := (*at.accounts)[id]
	if !exist {
		return nil, errors.New("Account id not found.")
	}

	return account, nil
}

func (at *AccountType) GetAccountByRef(ref int) (*account.Account, error) {
	for _, a := range *at.accounts {
		if (*a).GetRef() == ref {
			return a, nil
		}
	}

	return nil, errors.New("reference not found for this account type.")
}

func (at *AccountType) GetAccountByName(name string) (*account.Account, error) {
	for _, a := range *at.accounts {
		if (*a).GetName() == name {
			return a, nil
		}
	}

	return nil, errors.New("account type with name " + name + " not found")
}

func (at *AccountType) AddAccount(account *account.Account) (int, error) {
	if len(*at.accounts) == 999 {
		return 0, errors.New("account type full.")
	}

	_, exists := (*at.accounts)[account.GetId()]
	if exists {
		return 0, errors.New("account already in account type.")
	}

	(*at.accounts)[account.GetId()] = account

	return (at.refPrefix * 1000) + len(*at.accounts), nil
}

func (at *AccountType) String() string {
	var output strings.Builder
	output.WriteString("              Account Type\n")
	output.WriteString(strings.Repeat("─", 11))
	output.WriteString("─┬─")
	output.WriteString(strings.Repeat("─", 38))
	output.WriteString("\n")
	output.WriteString(fmt.Sprintf("       Name │ %s\n", at.name))
	output.WriteString(fmt.Sprintf("         Id │ %s\n", at.id))
	output.WriteString(fmt.Sprintf(" Ref Prefix │ %d\n", at.refPrefix))
	output.WriteString(fmt.Sprintf("   Accounts │ %d\n", len(*at.accounts)))

	return output.String()
}

func GetDefaultAccountTypes() map[id.Id]*AccountType {
	var accountTypes = make(map[id.Id]*AccountType, 0)
	var atId id.Id

	atId = id.GenerateNewId()
	accounts1 := make(map[id.Id]*account.Account)
	accountTypes[atId] = &AccountType{
		id: atId,
		name: "Assets",
		refPrefix: 1,
		accounts: &accounts1,
	}

	atId = id.GenerateNewId()
	accounts2 := make(map[id.Id]*account.Account)
	accountTypes[atId] = &AccountType{
		id: atId,
		name: "Liabilities",
		refPrefix: 2,
		accounts: &accounts2,
	}

	atId = id.GenerateNewId()
	accounts3 := make(map[id.Id]*account.Account)
	accountTypes[atId] = &AccountType{
		id: atId,
		name: "Equities",
		refPrefix: 3,
		accounts: &accounts3,
	}

	atId = id.GenerateNewId()
	accounts4 := make(map[id.Id]*account.Account)
	accountTypes[atId] = &AccountType{
		id: atId,
		name: "Revenues",
		refPrefix: 4,
		accounts: &accounts4,
	}

	atId = id.GenerateNewId()
	accounts5 := make(map[id.Id]*account.Account)
	accountTypes[atId] = &AccountType{
		id: atId,
		name: "Dividends",
		refPrefix: 5,
		accounts: &accounts5,
	}

	atId = id.GenerateNewId()
	accounts9 := make(map[id.Id]*account.Account)
	accountTypes[atId] = &AccountType{
		id: atId,
		name: "Expenses",
		refPrefix: 9,
		accounts: &accounts9,
	}

	return accountTypes
}