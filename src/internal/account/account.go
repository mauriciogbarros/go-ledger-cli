package account

import (
	"fmt"

	"go.mod/internal/id"
)

var MaxNameLength int = 38

type Account struct {
	id id.Id
	ref int
	name string
	accountType AccountType
}

func (a *Account) GetId() id.Id {
	return a.id
}

func (a *Account) GetRef() int {
	return a.ref
}

func (a *Account) GetName() string {
	return a.name
}

func (a *Account) GetAccountType() AccountType {
	return a.accountType
}

func (a *Account) GetAccountTypeInt() int {
	return int(a.accountType)
}

func NewAccount(ref int, name string, accountType int) Account {
	return Account{
		id: id.GenerateNewId(),
		ref: ref,
		name: name,
		accountType: AccountType(accountType),
	}
}

func NewAccountFromDb(id id.Id, ref int, name string, accountType int) Account {
	return Account{
		id: id,
		ref: ref,
		name: name,
		accountType: AccountType(accountType),
	}
}

func (a *Account) GetSide() AccountSide {
	switch a.accountType {
	case Asset, Expense:
		return Debit
	default:
		return Credit
	}
}

func (a *Account) String() string {
	var output string = " "
	output += fmt.Sprintf("%-*d", 3, a.ref)
	output += " │ "
	output += fmt.Sprintf("%-*s", MaxNameLength, a.name)
	output += " │ "
	output += a.accountType.String()
	output += "\n"

	return output
}

type AccountType int

const (
	Asset     AccountType = 1
	Liability AccountType = 2
	Equity    AccountType = 3
	Revenue   AccountType = 4
	Expense   AccountType = 5
)

func (at AccountType) String() string {
	switch at {
	case Asset:
		return "Asset"
	case Liability:
		return "Liability"
	case Equity:
		return "Equity"
	case Revenue:
		return "Revenue"
	case Expense:
		return "Expense"
	}
	return ""
}

type AccountSide int

const (
	Debit  AccountSide = 0
	Credit AccountSide = 1
)

func GetSideName(side int) string {
	switch side {
	case 0:
		return "Debit"
	case 1:
		return "Credit"
	}
	return ""
}
