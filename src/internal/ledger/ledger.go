package ledger

import (
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
func (l *Ledger) IsBalanced() bool {
	var isBalanced bool = true

	return isBalanced
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
		return entry.Entry{}, nil
	}
	debitAccount, err := accountType.GetAccountByRef(debitAccountRef)
	if err != nil {
		return entry.Entry{}, nil
	}

	refPrefix = creditAccountRef / 1000
	accountType, err = l.chart.GetAccountTypeByRefPrefix(refPrefix)
	if err != nil {
		return entry.Entry{}, nil
	}
	creditAccount, err := accountType.GetAccountByRef(creditAccountRef)
	if err != nil {
		return entry.Entry{}, nil
	}

	newEntry := entry.NewEntry(date, debitAccount, creditAccount, amount, explanation)

	return *newEntry, nil
}

func (l *Ledger) ViewAccount(id id.Id) string {
	var output string

	return output
}

func (l *Ledger) ViewTrialBalance() string {
	var output string

	return output
}

func (l *Ledger) String() string {
	var widthTitle int = 1 + 19 + 3 + 30 + 3 + 12 + 3 + 12 + 3 + 12 + 1
	var paddingTitleLeft int = (widthTitle - len(l.name)) / 2
	var widthSubTitle int = 1 + 19 + 3 + 12 + 3 + 12 + 1
	var paddingSubTitleLeft int = (widthSubTitle - 30 / 2)
	var output string
	output += strings.Repeat(" ", paddingTitleLeft)
	output += fmt.Sprintf("%s\n", l.name)
	output += strings.Repeat("─", widthTitle)
	if l.chart == nil {
		return ""
	}
	for _, at := range *l.chart.GetAccountTypes() {
		if at == nil {
			continue
		}
		accounts := at.GetAccounts()
		if accounts == nil {
			continue
		}
		if len(*accounts) > 0 {
			for _, a := range *at.GetAccounts() {
				if a == nil {
					continue
				}
				output += fmt.Sprintf("%-*s", paddingSubTitleLeft, a.GetName())
				output += fmt.Sprintf("%*d", widthTitle - widthSubTitle, a.GetRef())
				output += "\n"
				output += strings.Repeat("─", 1 + 19)
				output += "─┬─"
				output += strings.Repeat("─", 12)
				output += "─┬─"				
				output += strings.Repeat("─", 12)				
				output += "─┬─"				
				output += strings.Repeat("─", 12)				
				output += "─"				
				output += "\n"
				output += " "
				output += fmt.Sprintf("%-*s", 19, "Date")
				output += " │ "
				output += fmt.Sprintf("%*s", 12, "Debit")
				output += " │ "
				output += fmt.Sprintf("%*s", 12, "Credit")
				output += " │ "
				output += fmt.Sprintf("%*s", 12, "Balance")
				output += "\n"
				output += strings.Repeat("─", 1 + 19)
				output += "─┼─"
				output += strings.Repeat("─", 12)
				output += "─┼─"
				output += strings.Repeat("─", 12)				
				output += "─┼─"
				output += strings.Repeat("─", 12)				
				output += "─"				
				output += "\n"
				entries := a.GetEntries()
				if entries == nil || len(*entries) == 0 {
					output += strings.Repeat(" ", 1 + 19 + 3)
					output += "*No entries posted"
					output += "\n"
				}
				if len(*entries) > 0 {
					for _, e := range *entries{
						if e == nil {
							continue
						}
						output += fmt.Sprintf("%-*s", 1 + 19, e.GetDate().Format(time.DateTime))
						output += " │ "
						output += fmt.Sprintf("%*d", 12, e.GetAmount())
					}
				}
			}
		}
	}

	return output
}
