package account

import (
	"errors"
	"fmt"
	"sort"
	"strings"

	"go.mod/internal/currency"
	"go.mod/internal/id"
)

var MaxNameLength int = 38

type Account struct {
	id id.Id
	ref int
	name string
	description string
	entries *map[id.Id]*AccountEntry
	balance currency.Currency
}

func NewAccount(name string) *Account {
	entries := make(map[id.Id]*AccountEntry)
	return &Account{
		id: id.GenerateNewId(),
		ref: 0,
		name: name,
		entries: &entries,
		balance: currency.Currency(0),
	}
}

func NewDbAccount(
	sId string,
	ref int,
	name string,
) (*Account, error) {
	aId, err := id.ParseString(sId)
	if err != nil {
		return nil, err
	}

	entries := make(map[id.Id]*AccountEntry, 0)

	return &Account{
		id: aId,
		ref: ref,
		name: name,
		entries: &entries,
		balance: currency.Currency(0),
	}, nil
}

func (a *Account) GetAccount() *Account {
	return a
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

func (a *Account) GetEntries() *[]*AccountEntry {
	entries := make([]*AccountEntry, 0)
	for _, entry := range *a.entries {
		entries = append(entries, entry)
	}

	sort.Slice(entries, func(i, j int) bool {
		return entries[i].GetDate().Before(entries[j].GetDate())
	})

	return &entries
}

func (a *Account) GetEntryById(id id.Id) (*AccountEntry, error) {
	entry, exist := (*a.entries)[id]
	if !exist {
		return nil, errors.New("Entry id not found")
	}

	return entry, nil
}

func (a *Account) GetBalance() currency.Currency {
	return a.balance
}

func (a *Account) SetRef(ref int) {
	a.ref = ref
}

func (a *Account) AddEntry(entry *AccountEntry) error {
	_, exists := (*a.entries)[entry.GetId()]
	if exists {
		return errors.New("account entry already exists")
	}

	(*a.entries)[entry.GetId()] = entry
	return nil
}

func (a *Account) CalculateBalance() {
	a.balance = currency.Currency(0)

	for _, e := range *a.GetEntries() {
		if !e.GetSide() {
			a.balance += e.GetAmount()
			e.SetBalance(a.balance)
		} else {
			a.balance -= e.GetAmount()
			e.SetBalance(a.balance)
		}
	}
}

func (a *Account) String() string {
	var output string
	output += "           Account Details\n"
	output += strings.Repeat("─", 8)
	output += "─┬─"
	output += strings.Repeat("─\n", 36)
	output += fmt.Sprintf("      Id │ %s\n", a.id.String())
	output += fmt.Sprintf("     Ref │ %d\n", a.ref)
	output += fmt.Sprintf("    Name │ %s\n", a.name)
	output += fmt.Sprintf(" Entries | %d\n", len(*a.entries))
	output += fmt.Sprintf(" Balance │ %s\n", a.balance.String())
	output += "\n"

	return output
}

