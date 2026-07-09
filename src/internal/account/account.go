package account

import (
	"go.mod/internal/id"
	"go.mod/internal/post"
)

var MaxNameLength int = 38

type Account struct {
	ID   id.Id
	Name string
	Type AccountType
	Ledger []post.Post
}

func NewAccount(name string, tAccount int) Account {
	id := id.NewId()
	return Account{
		ID:   id,
		Name: name,
		Type: AccountType(tAccount),
		Ledger: make([]post.Post, 0),
	}
}

func (a Account) GetSide() AccountSide {
	switch a.Type {
	case Asset, Expense:
		return Debit
	default:
		return Credit
	}
}

type AccountType int

const (
	Asset     AccountType = iota
	Liability AccountType = iota
	Equity    AccountType = iota
	Revenue   AccountType = iota
	Expense   AccountType = iota
)

type AccountSide int

const (
	Debit  AccountSide = iota
	Credit AccountSide = iota
)