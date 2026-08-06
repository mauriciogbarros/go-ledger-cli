package ledger

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"go.mod/internal/account"
	"go.mod/internal/accountType"
	"go.mod/internal/chart"
	"go.mod/internal/currency"
	"go.mod/internal/entry"
	"go.mod/internal/id"
	"go.mod/internal/journal"
)

type Ledger struct {
	name string
	chart *chart.Chart
	journal *journal.Journal
}

func NewLedger() *Ledger {
	return &Ledger{
		name: "General Ledger",
		chart: nil,
		journal: nil,
	}
}

func (l *Ledger) GetLedger() *Ledger {
	return l
}

func (l *Ledger) GetName() string {
	return l.name
}

func (l *Ledger) GetChart() *chart.Chart {
	return l.chart
}

func (l *Ledger) GetJournal() *journal.Journal {
	return l.journal
}

func (l *Ledger) GetAccountTypes() *[]*accountType.AccountType {
	return l.chart.GetAccountTypes()
}

func (l *Ledger) GetAccountByStringId(stringId string) (*account.Account, error) {
	return l.chart.GetAccountByStringId(stringId)
}

func (l *Ledger) GetAccounts() *[]*account.Account {
	return l.chart.GetAccounts()
}

func (l *Ledger) SetChartAccountTypes(accountTypes *map[id.Id]*accountType.AccountType) {
	l.chart.SetAccountTypes(accountTypes)
}

func (l *Ledger) SetAccounts(accounts *[]*account.Account) error {
	return l.chart.MapAccountsToTypes(accounts)
}
func (l *Ledger) SetJournalEntries(entries *[]*entry.Entry) {
	l.journal.SetEntries(entries)
}

func (l *Ledger) ViewChart(refPrefix int) (string, error) {
	return l.chart.String(refPrefix)
}

func (l *Ledger) ViewAccountTypes() string {
	return l.chart.ViewAccountTypes()
}

func (l *Ledger) ViewAccountType(refPrefix int) (string, error) {
	return l.chart.ViewAccountType(refPrefix)
}

func (l *Ledger) ViewJournal() (string) {
	return l.journal.ViewJournal()
}

func (l *Ledger) ViewJournalInfo() string {
	return l.journal.String()
}

func (l *Ledger) CreateChart() {
	l.chart = chart.NewChart()
}

func (l *Ledger) CreateJournal() {
	l.journal = journal.NewJournal()
}

func (l *Ledger) CreateAccount(accountTypeRefPrefix int, name string, description string) (*account.Account, error) {
	return l.chart.CreateAccount(accountTypeRefPrefix, name, description)
}

func (l *Ledger) CreateJournalEntry(date time.Time, debitAccountRef int, creditAccountRef int, amountF64 float64, explanation string) (entry.Entry, error) {
	var accountType *accountType.AccountType
	var err error
	var refPrefix int

	amount := currency.ConvertF64(amountF64)

	refPrefix = debitAccountRef / 1000
	accountType, err = l.chart.GetAccountTypeByRefPrefix(refPrefix)
	if err != nil {
		return entry.Entry{}, err
	}
	debitAccount, err := accountType.GetAccountByRef(debitAccountRef)
	if err != nil {
		return entry.Entry{}, err
	}

	refPrefix = creditAccountRef / 1000
	accountType, err = l.chart.GetAccountTypeByRefPrefix(refPrefix)
	if err != nil {
		return entry.Entry{}, err
	}
	creditAccount, err := accountType.GetAccountByRef(creditAccountRef)
	if err != nil {
		return entry.Entry{}, err
	}

	newEntry := entry.NewEntry(date, debitAccount, creditAccount, amount, explanation)

	return *newEntry, nil
}

func (l *Ledger) ViewAccounts() (string, error) {
	return l.chart.ViewAccounts()
}

func (l *Ledger) ViewAccountByRef(ref int) (string, error) {
	return l.chart.ViewAccountByRef(ref)
}

func (l *Ledger) ViewTrialBalance() string {
	// TODO: View trial balance
	var output string

	return output
}

func (l *Ledger) ViewLedger() (string, error) {
	var widthTitle int = 10 + 3 + 12 + 3 + 12 + 3 + 12
	var paddingTitle int = (widthTitle - len(l.name)) / 2
	var widthSubtitle int = widthTitle - 4 - 3
	var output strings.Builder
	output.WriteString(strings.Repeat(" ", paddingTitle))
	fmt.Fprintf(&output, "%s\n", l.name)
	output.WriteString(strings.Repeat("─", widthTitle))
	output.WriteString("\n")
	if l.chart == nil {
		return "", errors.New("Chart - nil pointer dereference")
	}
	for _, at := range *l.chart.GetAccountTypes() {
		if at == nil {
			return "", errors.New("Account Types - nil pointer dereference")
		}
		fmt.Fprintf(&output, "%-*d", 4, at.GetRefPrefix())
		output.WriteString(strings.Repeat(" ", 3))
		output.WriteString(strings.Repeat(" ", (widthSubtitle - len(at.GetName())) / 2))
		fmt.Fprintf(&output, "%s\n", at.GetName())
		output.WriteString(strings.Repeat("─", widthTitle))
		output.WriteString("\n")
		accounts := at.GetAccounts()
		if accounts == nil {
			return "", errors.New("Accounts - nil pointer dereference")
		}
		if len(*accounts) > 0 {
			for _, a := range *at.GetAccounts() {
				if a == nil {
					return "", errors.New("Account - nil pointer dereference")
				}
				fmt.Fprintf(&output, "%-*d", 4, a.GetRef())
				output.WriteString(strings.Repeat(" ", 3))
				output.WriteString(strings.Repeat(" ", (widthSubtitle - len(a.GetName())) / 2))
				fmt.Fprintf(&output, "%s\n", a.GetName())
				output.WriteString(strings.Repeat("─", 10))
				output.WriteString("─┬─")
				output.WriteString(strings.Repeat("─", 12))
				output.WriteString("─┬─")
				output.WriteString(strings.Repeat("─", 12))
				output.WriteString("─┬─")
				output.WriteString(strings.Repeat("─", 12))
				output.WriteString("\n")
				fmt.Fprintf(&output, "%-*s", 10, "Date")
				output.WriteString(" │ ")
				fmt.Fprintf(&output, "%*s", 12, "Debit")
				output.WriteString(" │ ")
				fmt.Fprintf(&output, "%*s", 12, "Credit")
				output.WriteString(" │ ")
				fmt.Fprintf(&output, "%*s", 12, "Balance")
				output.WriteString("\n")
				output.WriteString(strings.Repeat("─", 10))
				output.WriteString("─┼─")
				output.WriteString(strings.Repeat("─", 12))
				output.WriteString("─┼─")
				output.WriteString(strings.Repeat("─", 12))
				output.WriteString("─┼─")
				output.WriteString(strings.Repeat("─", 12))
				output.WriteString("\n")
				entries := a.GetEntries()
				if entries == nil || len(*entries) == 0 {
					output.WriteString(strings.Repeat(" ", 10 + 3))
					output.WriteString("*No entries posted")
					output.WriteString("\n")
				}
				if len(*entries) > 0 {
					for _, e := range *entries{
						if e == nil {
							return "", errors.New("Entries - nil pointer dereference")
						}
						fmt.Fprintf(&output, "%-*s", 10, e.GetDate().Format(time.DateOnly))
						output.WriteString(" │ ")
						fmt.Fprintf(&output, "%*d", 12, e.GetAmount())
					}
				}
				output.WriteString("\n")
			}
		}
	}

	return output.String(), nil
}

func (l *Ledger) String() string {
	title := "Ledger Information"
	width := 13 + 3 + len(l.name)
	padding := (width - len(title)) / 2
	nEntries := len(*l.journal.GetEntriesPosted(true))
	var output strings.Builder
	output.WriteString(strings.Repeat(" ", padding))
	fmt.Fprintf(&output, "%s\n", title)
	output.WriteString(strings.Repeat("─", 13))
	output.WriteString("─┬─")
	output.WriteString(strings.Repeat("─", len(l.name)))
	output.WriteString("\n")
	fmt.Fprintf(&output, "         Name │ %s\n", l.name)
	fmt.Fprintf(&output, "Account Types │ %d\n", len(*l.GetAccountTypes()))
	fmt.Fprintf(&output, "     Accounts │ %d\n", len(*l.GetAccounts()))
	fmt.Fprintf(&output, "      Entries │ %d\n", nEntries)
	fmt.Fprintf(&output, "      Balance │ %s\n", l.CalculateBalance().String())
	return output.String()
}

func (l *Ledger) CalculateBalance() currency.Currency {
	// TODO: Calculate balance
	balance := currency.Currency(0)
	return balance
}
