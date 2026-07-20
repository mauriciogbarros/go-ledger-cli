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
	accountTypeId id.Id
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

func (a *Account) GetAccountTypeId() id.Id {
	return a.accountTypeId
}

func (a *Account) SetRef(ref int) {
	a.ref = ref
}


func NewAccount(name string, accountTypeId id.Id) *Account {
	return &Account{
		id: id.GenerateNewId(),
		ref: 0,
		name: name,
		accountTypeId: accountTypeId,
	}
}

func CreateAccountFromDb(sId string, ref int, name string, sAccTypeId string) (Account, error) {
	accountId, err := id.ParseString(sId)
	if err != nil {
		return Account{}, err
	}

	accountTypeId, err := id.ParseString(sAccTypeId)
	if err != nil {
		return Account{}, err
	}

	return Account{
		id: accountId,
		ref: ref,
		name: name,
		accountTypeId: accountTypeId,
	}, nil
}

func (a *Account) String() string {
	var output string = " "
	output += fmt.Sprintf("%-*d", 3, a.ref)
	output += " │ "
	output += fmt.Sprintf("%-*s", MaxNameLength, a.name)
	output += " │"

	return output
}

