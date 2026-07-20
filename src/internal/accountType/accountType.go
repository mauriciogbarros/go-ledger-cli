package accountType

import (
	"errors"
	"fmt"

	"go.mod/internal/account"
	"go.mod/internal/id"
)

var MaxNameLength int = 20

type AccountType struct {
	id id.Id
	name string
	side int
	refCounter int
	accounts []*account.Account
}

func (at *AccountType) GetId() id.Id {
	return at.id
}

func (at *AccountType) GetName() string {
	return at.name
}

func (at *AccountType) GetSide() string {
	if at.side == 0 {
		return "Debit"
	} else {
		return "Credit"
	}
}

func (at *AccountType) GetSideNum() int {
	return at.side
}

func (at *AccountType) GetRefCounter() int {
	return at.refCounter
}

func (at *AccountType) GetRefGroup() int {
	return (at.refCounter / 1000) * 1000
}

func (at *AccountType) addToCounter() (int, error) {
	var count int = at.refCounter % 1000
	if count < 999 {
		at.refCounter++
		return at.refCounter, nil
	} else {
		return -1, errors.New("maximum reference number reached for account type" + at.name)
	}
}

func (at *AccountType) AddAccount(account *account.Account) error {
	ref, err := at.addToCounter()
	if err != nil {
		return err
	}

	account.SetRef(ref)
	at.accounts = append(at.accounts, account)
	return nil
}

func (at *AccountType) GetAccounts() []*account.Account {
	return at.accounts
}

func (at *AccountType) GetAccountByRef(ref int) (*account.Account, error) {
	for _, a := range at.accounts {
		if (*a).GetRef() == ref {
			return a, nil
		}
	}

	return nil, errors.New("reference not found for this account type.")
}

func (at *AccountType) GetAccountByName(name string) (*account.Account, error) {
	for _, a := range at.accounts {
		if (*a).GetName() == name {
			return a, nil
		}
	}

	return nil, errors.New("account type with name " + name + " not found")
}

func CreateAccountTypeFromDb(sId string, name string, side int, refCounter int) (AccountType, error) {
	id, err := id.ParseString(sId)
	if err != nil {
		return AccountType{}, err
	}

	var at = AccountType{
		id: id,
		name: name,
		side: side,
		refCounter: refCounter,
		accounts: make([]*account.Account, 0),
	}

	return at, nil
}

func CreateDefaultAccountTypes() *[]AccountType {
	var accountTypes = make([]AccountType, 0)

	var defaultAccountTypes = []AccountType{
		{
			id: id.GenerateNewId(),
			name: "Asset",
			side: 0,
			refCounter: 1000,
			accounts: make([]*account.Account, 0),
		},
		{
			id: id.GenerateNewId(),
			name: "Liability",
			side: 1,
			refCounter: 2000,
			accounts: make([]*account.Account, 0),
		},
		{
			id: id.GenerateNewId(),
			name: "Equity",
			side: 1,
			refCounter: 3000,
			accounts: make([]*account.Account, 0),
		},
		{
			id: id.GenerateNewId(),
			name: "Revenue",
			side: 1,
			refCounter: 4000,
			accounts: make([]*account.Account, 0),
		},
		{
			id: id.GenerateNewId(),
			name: "Dividends",
			side: 0,
			refCounter: 5000,
			accounts: make([]*account.Account, 0),
		},
		{
			id: id.GenerateNewId(),
			name: "Expense",
			side: 0,
			refCounter: 9000,
			accounts: make([]*account.Account, 0),
		},
	}
	
	for _, at := range defaultAccountTypes {
		accountTypes = append(accountTypes, at)
	}

	return &accountTypes
}

func (at *AccountType) String() string {
	var output string = " "
	output += fmt.Sprintf("%-*s", MaxNameLength, at.GetName())
	output += " │ "
	output += fmt.Sprintf("%-*s", 6, at.GetSide())

	return output
}