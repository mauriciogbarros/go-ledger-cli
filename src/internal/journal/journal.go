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
	var width int = 1 + 19 + 3 + account.MaxNameLength + 2 + 9 + 12 + 3 + 12 + 1
	var paddingLeft = (width - len(j.name))/2

	var output string
	output += strings.Repeat(" ", paddingLeft)
	output += j.name
	output += "\n"
	output += strings.Repeat("─", width)
	output += "\n"
	output += " "
	output += fmt.Sprintf("%-*s", 19, "Date")
	output += "   "
	output += fmt.Sprintf("%-*s", account.MaxNameLength + 2, "Accounts & Explanation")
	output += "   "
	output += "Ref"
	output += "   "
	output += fmt.Sprintf("%*s", 12, "Debit")
	output += "   "
	output += fmt.Sprintf("%*s", 12, "Credit")
	output += "\n"
	output += strings.Repeat("─", 20)
	output += "─┬─"
	output += strings.Repeat("─", account.MaxNameLength + 2)
	output += "─┬─"
	output += strings.Repeat("─", 3)
	output += "─┬─"
	output += strings.Repeat("─", 12)
	output += "─┬─"
	output += strings.Repeat("─", 13)
	output += "\n"

	if len(j.entries) == 0 {
		output += " No entires\n"
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
			output += e.Format(debitAccount.GetName(), creditAccount.GetName())
		}
	}
	return output
}
