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
		return nil
	}
	for _, entry := range *j.entries {
		entries = append(entries, entry)
	}

	sort.Slice(entries, func(i, j int) bool {
		return entries[i].GetDate().Before(entries[j].GetDate())
	})

	return &entries
}

func (j *Journal) GetEntriesPosted(arePosted bool) *[]*entry.Entry {
	entries := make([]*entry.Entry, 0)
	if j.entries == nil {
		return nil
	}
	if arePosted {
		for _, entry := range *j.entries {
			if entry.IsPosted() {
				entries = append(entries, entry)
			}
		}
	} else {
		for _, entry := range *j.entries {
			if !entry.IsPosted() {
				entries = append(entries, entry)
			}
		}
	}

	sort.Slice(entries, func(i, j int) bool {
		return entries[i].GetDate().Before(entries[j].GetDate())
	})

	return &entries
}

func (j *Journal) GetEntryById(id id.Id) (*entry.Entry, error) {
	if j.entries == nil {
		return nil, errors.New("Entries - nil pointer dereference")
	}
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

func (j *Journal) ViewJournal() string {
	var width int = 10 + 3 + account.MaxNameLength + 4 + 3 + 4 + 3 + 12 + 3 + 12
	var padding = (width - len(j.name))/2
	var output strings.Builder
	output.WriteString(strings.Repeat(" ", padding))
	output.WriteString(j.name)
	output.WriteString("\n")
	output.WriteString(strings.Repeat("─", 10))
	output.WriteString("─┬─")
	output.WriteString(strings.Repeat("─", account.MaxNameLength + 4))
	output.WriteString("─┬─")
	output.WriteString(strings.Repeat("─", 4))
	output.WriteString("─┬─")
	output.WriteString(strings.Repeat("─", 12))
	output.WriteString("─┬─")
	output.WriteString(strings.Repeat("─", 12))
	output.WriteString("\n")
	fmt.Fprintf(&output, "%-*s", 10, "Date")
	output.WriteString(" │ ")
	fmt.Fprintf(&output, "%-*s", account.MaxNameLength + 4, "Accounts & Explanation")
	output.WriteString(" │ ")
	output.WriteString("Ref ")
	output.WriteString(" │ ")
	fmt.Fprintf(&output, "%*s", 12, "Debit")
	output.WriteString(" │ ")
	fmt.Fprintf(&output, "%*s\n", 12, "Credit")
	output.WriteString(strings.Repeat("─", 10))
	output.WriteString("─┼─")
	output.WriteString(strings.Repeat("─", account.MaxNameLength + 4))
	output.WriteString("─┼─")
	output.WriteString(strings.Repeat("─", 4))
	output.WriteString("─┼─")
	output.WriteString(strings.Repeat("─", 12))
	output.WriteString("─┼─")
	output.WriteString(strings.Repeat("─", 12))
	output.WriteString("\n")

	if j.entries == nil || len(*j.entries) == 0 {
		output.WriteString(strings.Repeat(" ", 10 + 3))
		output.WriteString("*No entires\n")
	} else {
		entries := j.GetEntries()
		for _, e := range *entries {
			debitAccount := e.GetDebitAccount()
			creditAccount := e.GetCreditAccount()
			if e == nil || debitAccount == nil || creditAccount == nil {
				continue
			}
			fmt.Fprintf(&output, "%-*s", 10, e.GetDate().Format(time.DateOnly))
			output.WriteString(" │ ")
			fmt.Fprintf(&output, "%-*s", account.MaxNameLength + 4, debitAccount.GetName())
			output.WriteString(" │ ")
			if e.IsPosted() {
				fmt.Fprintf(&output, "%*d", 4, debitAccount.GetRef())
			} else {
				output.WriteString(strings.Repeat(" ", 4))
			}
			output.WriteString(" │ ")
			fmt.Fprintf(&output, "%*s", 12, e.GetAmount().String())
			output.WriteString(" │\n")
			output.WriteString(strings.Repeat(" ", 10))
			output.WriteString(" │   ")
			fmt.Fprintf(&output, "%-*s", account.MaxNameLength + 2, creditAccount.GetName())
			output.WriteString(" │ ")
			if e.IsPosted() {
				fmt.Fprintf(&output, "%*d", 4, creditAccount.GetRef())
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
	
	return output.String()
}

func (j *Journal) String() string {
	nEntriesPosted := len(*j.GetEntriesPosted(true))
	nEntriesNotPosted := len(*j.GetEntriesPosted(false))
	title := "Journal Information"
	width := 18 + 3 + len(j.name)
	padding := (width - len(title)) / 2
	var output strings.Builder
	output.WriteString(strings.Repeat(" ", padding))
	fmt.Fprintf(&output, "%s\n", title)
	output.WriteString(strings.Repeat("─", 18))
	output.WriteString("─┬─")
	output.WriteString(strings.Repeat("─", len(j.name)))
	output.WriteString("\n")
	fmt.Fprintf(&output, "              Name │ %s\n", j.name)
	fmt.Fprintf(&output, "    Entries Posted │ %d\n", nEntriesPosted)
	fmt.Fprintf(&output, "Entries Not-Posted │ %d\n", nEntriesNotPosted)
	return output.String()
}