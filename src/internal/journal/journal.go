package journal

import (
	"fmt"
	"strings"

	"go.mod/db"
	"go.mod/internal/account"
	"go.mod/internal/chart"
	"go.mod/internal/entry"
)

type Journal struct {
	name string
	chart *chart.ChartOfAccounts
	entries []entry.Entry
}

func (j *Journal) NewJournal(name string, chart *chart.ChartOfAccounts) error {
	j.name = name
	j.chart = chart
	j.entries = make([]entry.Entry, 0)
	entries, err := db.GetEntries()
	if err != nil {
		return err
	}
	j.entries = append(j.entries, entries...)
	return nil
}

func (j *Journal) AddEntry(entry entry.Entry) error {
	fmt.Println("Adding entry")
	j.entries = append(j.entries, entry)
	err := db.CreateEntry(entry)
	if err != nil {
		return err
	}
	return nil
}

func (j Journal) String() string {
	var journal string
	journal += " "
	journal += fmt.Sprintf("%-*s", 19, "Date")
	journal += "   "
	journal += fmt.Sprintf("%-*s", account.MaxNameLength + 2, "Accounts & Explanation")
	journal += "   "
	journal += fmt.Sprintf("%-*s", 5, "Ref")
	journal += "   "
	journal += fmt.Sprintf("%*s", 12, "Debit")
	journal += "   "
	journal += fmt.Sprintf("%*s", 12, "Credit")
	journal += "\n"
	journal += strings.Repeat("─", 20)
	journal += "─┬─"
	journal += strings.Repeat("─", account.MaxNameLength + 2)
	journal += "─┬─"
	journal += strings.Repeat("─", 3)
	journal += "─┬─"
	journal += strings.Repeat("─", 12)
	journal += "─┬─"
	journal += strings.Repeat("─", 13)
	journal += "\n"

	if len(j.entries) == 0 {
		journal += " No entires\n"
	} else {
		for _, e := range j.entries {
			debitAccount, err := j.chart.GetAccountByRef(e.GetDebitAccountRef())
			if err != nil {
				panic(err)
			}
			creditAccount, err := j.chart.GetAccountByRef(e.GetCreditAccountRef())
			if err != nil {
				panic(err)
			}
			journal += e.Format(debitAccount.GetName(), creditAccount.GetName())
		}
	}
	return journal
}
