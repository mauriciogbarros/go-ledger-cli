package journal

import (
	"errors"
	"fmt"
	"sort"
	"strings"
	"time"

	"go.mod/internal/entry"
	"go.mod/internal/id"
)

type Journal struct {
	name string
	entries *map[id.Id]*entry.Entry
}

func NewJournal() *Journal {
	entries := make(map[id.Id]*entry.Entry, 0)
	return &Journal{
		name: "General Journal",
		entries: &entries,
	}
}

func (j *Journal) GetName() string {
	return j.name
}

func (j *Journal) GetEntries() *[]*entry.Entry {
	entries := make([]*entry.Entry, 0)
	if j.entries == nil {
		return &entries
	}
	for _, entry := range *j.entries {
		entries = append(entries, entry)
	}

	sort.Slice(entries, func(i, j int) bool {
		return entries[i].GetDate().Before(entries[j].GetDate())
	})

	return &entries
}

func (j *Journal) GetEntryById(id id.Id) (*entry.Entry, error) {
	entry, exist := (*j.entries)[id]
	if !exist {
		return nil, errors.New("Entry id not found.")
	}
	
	return entry, nil
}

func (j *Journal) SetEntries(entries *[]*entry.Entry) error {
	if entries == nil {
		return errors.New("No entries")
	}
	clear(*j.entries)
	for _, e := range *entries {
		if e == nil {
			continue
		}
		(*j.entries)[e.GetId()] = e
	}
	return nil
}

func (j Journal) String() string {
	var width int = 1 + 19 + 3 + 34 + 3 + 4 + 3 + 12 + 3 + 12 + 1
	var paddingLeft = (width - len(j.name))/2
	var output strings.Builder
	output.WriteString(strings.Repeat(" ", paddingLeft))
	output.WriteString(j.name)
	output.WriteString("\n")
	output.WriteString(strings.Repeat("─", 1 + 19))
	output.WriteString("─┬─")
	output.WriteString(strings.Repeat("─", 34))
	output.WriteString("─┬─")
	output.WriteString(strings.Repeat("─", 4))
	output.WriteString("─┬─")
	output.WriteString(strings.Repeat("─", 12))
	output.WriteString("─┬─")
	output.WriteString(strings.Repeat("─", 12 + 1))
	output.WriteString("\n")
	output.WriteString(" ")
	output.WriteString(fmt.Sprintf("%-*s", 19, "Date"))
	output.WriteString(" │ ")
	output.WriteString(fmt.Sprintf("%-*s", 34, "Accounts & Explanation"))
	output.WriteString(" │ ")
	output.WriteString("Ref ")
	output.WriteString(" │ ")
	output.WriteString(fmt.Sprintf("%*s", 12, "Debit"))
	output.WriteString(" │ ")
	output.WriteString(fmt.Sprintf("%*s", 12, "Credit"))
	output.WriteString("\n")
	output.WriteString(strings.Repeat("─", 1 + 19))
	output.WriteString("─┼─")
	output.WriteString(strings.Repeat("─", 34))
	output.WriteString("─┼─")
	output.WriteString(strings.Repeat("─", 4))
	output.WriteString("─┼─")
	output.WriteString(strings.Repeat("─", 12))
	output.WriteString("─┼─")
	output.WriteString(strings.Repeat("─", 12 + 1))
	output.WriteString("\n")

	if j.entries == nil || len(*j.entries) == 0 {
		output.WriteString(strings.Repeat(" ", 1 + 19 + 3))
		output.WriteString("*No entires\n")
	} else {
		entries := j.GetEntries()
		for _, e := range *entries {
			if e == nil || e.GetDebitAccount() == nil || e.GetCreditAccount() == nil {
				continue
			}
			output.WriteString(" ")
			output.WriteString(fmt.Sprintf("%-*s", 19, e.GetDate().Format(time.DateTime)))
			output.WriteString(" │ ")
			output.WriteString(fmt.Sprintf("%-*s", 38 + 2, e.GetDebitAccount().GetName()))
			output.WriteString(" │ ")
			if e.IsPosted() {
				output.WriteString(fmt.Sprintf("%*d", 4, e.GetDebitAccount().GetRef()))
			} else {
				output.WriteString(strings.Repeat(" ", 4))
			}
			output.WriteString(" │ ")
			output.WriteString(fmt.Sprintf("%*d", 12, e.GetAmount()))
			output.WriteString(" │\n")
			output.WriteString(" ")
			output.WriteString(strings.Repeat(" ", 19))
			output.WriteString(" │   ")
			output.WriteString(fmt.Sprintf("%-*s", 38, e.GetCreditAccount().GetName()))
			output.WriteString(" │ ")
			if e.IsPosted() {
				output.WriteString(fmt.Sprintf("%*d", 4, e.GetCreditAccount().GetRef()))
			} else {
				output.WriteString(strings.Repeat(" ", 4))
			}
			output.WriteString(" │ ")
			output.WriteString(strings.Repeat(" ", 12))
			output.WriteString(" │ ")
			output.WriteString(fmt.Sprintf("%*s", 12, e.GetAmount().String()))
			output.WriteString("\n")
			output.WriteString(" ")
			output.WriteString(strings.Repeat(" ", 19))
			output.WriteString(" │     ")
			output.WriteString(fmt.Sprintf("%-*s", 36, e.GetExplanation()))
			output.WriteString(" │ ")
			output.WriteString(strings.Repeat(" ", 4))
			output.WriteString(" │ ")
			output.WriteString(strings.Repeat(" ", 12))
			output.WriteString(" │ ")
			output.WriteString(strings.Repeat(" ", 12))
			output.WriteString("\n")
		}
	}
	output.WriteString("\n")
	
	return output.String()
}
