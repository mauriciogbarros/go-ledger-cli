package journal

import (
	"fmt"
	"strings"

	"go.mod/internal/account"
	"go.mod/internal/entry"
)

type Journal struct {
	Name string
	Entries []entry.Entry
}

func (j *Journal) NewJournal() {
	j.Entries = make([]entry.Entry, 0)
}

func (j *Journal) AddEntry(entry entry.Entry) {
	j.Entries = append(j.Entries, entry)
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
	journal += strings.Repeat("─", 34)
	journal += "─┬─"
	journal += strings.Repeat("─", 5)
	journal += "─┬─"
	journal += strings.Repeat("─", 12)
	journal += "─┬─"
	journal += strings.Repeat("─", 13)
	journal += "\n"

	if len(j.Entries) == 0 {
		journal += " No entires\n"
	} else {
		for _, entry := range j.Entries {
			journal += entry.String()
		}
	}
	return journal
}
