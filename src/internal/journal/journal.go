package journal

import (
	"fmt"
	"strings"

	"go.mod/internal/account"
	"go.mod/internal/chart"
	"go.mod/internal/entry"
)

type Journal struct {
	name string
	chart *chart.ChartOfAccounts
	entries *[]*entry.Entry
}

func NewJournal(name string, chart *chart.ChartOfAccounts) *Journal {
	entries := make([]*entry.Entry, 0)
	return &Journal{
		name: name,
		chart: chart,
		entries: &entries,
	}
}

func (j *Journal) GetName() string {
	return j.name
}

func (j *Journal) GetChart() *chart.ChartOfAccounts {
	return j.chart
}

func (j *Journal) GetEntries() *[]*entry.Entry {
	return j.entries
}

func (j *Journal) SetEntries(entries *[]*entry.Entry) {
	j.entries = entries
}

func (j Journal) String() string {
	var width int = 1 + 19 + 3 + account.MaxNameLength + 2 + 9 + 12 + 4 + 12 + 1
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
	output += "Ref "
	output += "   "
	output += fmt.Sprintf("%*s", 12, "Debit")
	output += "   "
	output += fmt.Sprintf("%*s", 12, "Credit")
	output += "\n"
	output += strings.Repeat("─", 20)
	output += "─┬─"
	output += strings.Repeat("─", account.MaxNameLength + 2)
	output += "─┬─"
	output += strings.Repeat("─", 4)
	output += "─┬─"
	output += strings.Repeat("─", 12)
	output += "─┬─"
	output += strings.Repeat("─", 13)
	output += "\n"

	if len(*j.entries) == 0 {
		output += strings.Repeat(" ", 1 + 19 + 3)
		output += "*No entires\n"
	} else {
		for _, e := range *j.entries {
			output += e.String()
		}
	}
	return output
}
