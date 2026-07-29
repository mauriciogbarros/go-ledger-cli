package journal

import (
	"errors"
	"fmt"
	"sort"
	"strings"
	"time"

	"go.mod/internal/account"
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
	if j.entries == nil {
		return errors.New("journal entries map is nil")
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
	var width int = 1 + 10 + 3 + account.MaxNameLength + 4 + 3 + 4 + 3 + 12 + 3 + 12 + 1
	var paddingLeft = (width - len(j.name))/2
	var output strings.Builder
	output.WriteString(strings.Repeat(" ", paddingLeft))
	output.WriteString(j.name)
	output.WriteString("\n")
	output.WriteString(strings.Repeat("─", 1 + 10))
	output.WriteString("─┬─")
	output.WriteString(strings.Repeat("─", account.MaxNameLength + 4))
	output.WriteString("─┬─")
	output.WriteString(strings.Repeat("─", 4))
	output.WriteString("─┬─")
	output.WriteString(strings.Repeat("─", 12))
	output.WriteString("─┬─")
	output.WriteString(strings.Repeat("─", 12 + 1))
	output.WriteString("\n")
	output.WriteString(" ")
	fmt.Fprintf(&output, "%-*s", 10, "Date")
	output.WriteString(" │ ")
	fmt.Fprintf(&output, "%-*s", account.MaxNameLength + 4, "Accounts & Explanation")
	output.WriteString(" │ ")
	output.WriteString("Ref ")
	output.WriteString(" │ ")
	fmt.Fprintf(&output, "%*s", 12, "Debit")
	output.WriteString(" │ ")
	fmt.Fprintf(&output, "%*s", 12, "Credit")
	output.WriteString("\n")
	output.WriteString(strings.Repeat("─", 1 + 10))
	output.WriteString("─┼─")
	output.WriteString(strings.Repeat("─", account.MaxNameLength + 4))
	output.WriteString("─┼─")
	output.WriteString(strings.Repeat("─", 4))
	output.WriteString("─┼─")
	output.WriteString(strings.Repeat("─", 12))
	output.WriteString("─┼─")
	output.WriteString(strings.Repeat("─", 12 + 1))
	output.WriteString("\n")

	if j.entries == nil || len(*j.entries) == 0 {
		output.WriteString(strings.Repeat(" ", 1 + 10 + 3))
		output.WriteString("*No entires\n")
	} else {
		entries := j.GetEntries()
		for _, e := range *entries {
			if e == nil || e.GetDebitAccount() == nil || e.GetCreditAccount() == nil {
				continue
			}
			output.WriteString(" ")
			fmt.Fprintf(&output, "%-*s", 10, e.GetDate().Format(time.DateOnly))
			output.WriteString(" │ ")
			fmt.Fprintf(&output, "%-*s", account.MaxNameLength + 4, e.GetDebitAccount().GetName())
			output.WriteString(" │ ")
			if e.IsPosted() {
				fmt.Fprintf(&output, "%*d", 4, e.GetDebitAccount().GetRef())
			} else {
				output.WriteString(strings.Repeat(" ", 4))
			}
			output.WriteString(" │ ")
			fmt.Fprintf(&output, "%*s", 12, e.GetAmount().String())
			output.WriteString(" │\n")
			output.WriteString(" ")
			output.WriteString(strings.Repeat(" ", 10))
			output.WriteString(" │   ")
			fmt.Fprintf(&output, "%-*s", account.MaxNameLength + 2, e.GetCreditAccount().GetName())
			output.WriteString(" │ ")
			if e.IsPosted() {
				fmt.Fprintf(&output, "%*d", 4, e.GetCreditAccount().GetRef())
			} else {
				output.WriteString(strings.Repeat(" ", 4))
			}
			output.WriteString(" │ ")
			output.WriteString(strings.Repeat(" ", 12))
			output.WriteString(" │ ")
			fmt.Fprintf(&output,"%*s", 12, e.GetAmount().String())
			output.WriteString("\n")
			words := strings.Split(e.GetExplanation(), " ")
			var explanation strings.Builder
			for i := 0; i < len(words); {
				for exLen := 0; exLen <= account.MaxNameLength && i < len(words); {
					explanation.WriteString(words[i])
					explanation.WriteString(" ")
					if i < len(words) {
						i++
					}
					if i < len(words) {
						exLen = explanation.Len() + len(words[i])
					}
				}

				output.WriteString(" ")
				output.WriteString(strings.Repeat(" ", 10))
				output.WriteString(" │     ")
				fmt.Fprintf(&output, "%-*s", account.MaxNameLength, &explanation)
				output.WriteString(" │ ")
				output.WriteString(strings.Repeat(" ", 4))
				output.WriteString(" │ ")
				output.WriteString(strings.Repeat(" ", 12))
				output.WriteString(" │ ")
				output.WriteString(strings.Repeat(" ", 12))
				output.WriteString("\n")

				explanation.Reset()
			}
		}
	}
	output.WriteString("\n")
	
	return output.String()
}
