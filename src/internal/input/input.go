package input

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"slices"
	"strconv"
	"strings"
	"time"

	"go.mod/internal/account"
	"go.mod/internal/accountType"
	"go.mod/internal/entry"
	"go.mod/internal/ledger"
)

var reader = bufio.NewReader(os.Stdin)

func InputAccountName() (string, error) {
	fmt.Print("Account name: ")
	name, err := reader.ReadString('\n')
	name = strings.TrimSpace(name)
	if err != nil {
		return "", err
	}
	if len(name) == 0 {
		return "", errors.New("Account name cannot be empty")
	}
	if len(name) > account.MaxNameLength {
		return "", fmt.Errorf("Account name too long (max %d characters)", account.MaxNameLength)
	}

	return name, nil
}

func InputAccountTypeRefPrefix(ledger *ledger.Ledger) (int, error) {
	if ledger == nil {
		return 0, errors.New("ledger is nil")
	}
	accountTypes := ledger.GetAccountTypes()
	if accountTypes == nil {
		return 0, errors.New("no account types available")
	}
	options := make([]int, 0)
	fmt.Println("Choose an account type:")
	fmt.Println(strings.Repeat("─", 3 + accountType.MaxNameLength))
	for _, at := range *accountTypes {
		if at == nil {
			continue
		}
		options = append(options, at.GetRefPrefix())
		fmt.Printf("%d. %s\n", at.GetRefPrefix(), at.GetName())
	}
	fmt.Print("Choice: ")

	input, err := reader.ReadString('\n')
	if err != nil {
		return 0, err
	}
	input = strings.TrimSpace(input)
	choice, err := strconv.Atoi(input)
	if err != nil {
		return 0, errors.New("Invalid input")
	}

	exists := false
	for _, o := range options {
		if choice == o {
			exists = true
			break
		}
	}
	if !exists {
		return 0, errors.New("Invalid choice.")
	}
	return choice, nil
}

func InputAmountF64() (float64, error) {
	fmt.Print("Amount: ")
	amountString, err := reader.ReadString('\n')
	if err != nil {
		return 0, err
	}
	amountString = strings.TrimSpace(amountString)
	amountString = strings.ReplaceAll(amountString, ",", "")
	amountF64, err := strconv.ParseFloat(amountString, 64)
	if err != nil {
		return 0, err
	}
	return amountF64, nil
}

func InputDate(dateField string) (time.Time, error) {
	fmt.Printf("%s date (YYYY-MM-DD): ", dateField)
	dateString, err := reader.ReadString('\n')
	if err != nil {
		return time.Time{}, err
	}
	dateString = strings.TrimSpace(dateString)
	if len(dateString) == 0 {
		return time.Time{}, nil
	}
	parts := strings.Split(dateString, "-")
	if len(parts) != 3 {
		return time.Time{}, fmt.Errorf("invalid date format, use YYYY-MM-DD")
	}
	layout := "2006-"
	if len(parts[1]) == 1 {
		layout += "1-"
	} else {
		layout += "01-"
	}
	if len(parts[2]) == 1 {
		layout += "2"
	} else {
		layout += "02"
	}
	date, err := time.Parse(layout, dateString)
	if err != nil {
		return time.Time{}, fmt.Errorf("invalid date format, use YYYY-MM-DD")
	}
	return date, nil
}

func InputAccountRef(ledger *ledger.Ledger, side string) (int, error) {
	if ledger == nil {
		return 0, errors.New("ledger is nil")
	}
	accounts := ledger.GetAccounts()
	if accounts == nil {
		return 0, errors.New("no accounts available")
	}
	var menu strings.Builder
	menu.WriteString(" Ref Accounts\n")
	menu.WriteString("─")
	menu.WriteString(strings.Repeat("─", 1 + 3))
	menu.WriteString("─┬─")
	menu.WriteString(strings.Repeat("─", account.MaxNameLength))
	menu.WriteString("─")
	menu.WriteString("\n")
	var refs []int
	for _, account := range *accounts {
		if account == nil {
			continue
		}
		ref := account.GetRef()
		refs = append(refs, ref)
		fmt.Fprintf(&menu, " %d │ %s\n", ref, account.GetName())
	}
	menu.WriteString("─")
	menu.WriteString(strings.Repeat("─", 1 + 3))
	menu.WriteString("─┴─")
	menu.WriteString(strings.Repeat("─", account.MaxNameLength))
	menu.WriteString("─")
	fmt.Println(menu.String())
	if len(side) == 0 {
		fmt.Print("Enter ref: ")
	} else {
		fmt.Printf("Enter %s ref: ", side)
	}
	input, err := reader.ReadString('\n')
	fmt.Println()
	if err != nil {
		return 0, err
	}
	input = strings.TrimSpace(input)
	ref, err := strconv.Atoi(input)
	if err != nil {
		return 0, err
	}
	index := slices.Index(refs, ref)
	if index < 0 {
		return 0, errors.New("Invalid account reference")
	}
	return ref, nil
}

func InputText(fieldName string) (string, error) {
	fmt.Printf("%s: ", fieldName)
	explanation, err := reader.ReadString('\n')
	explanation = strings.TrimSpace(explanation)
	if err != nil {
		return "", err
	}
	if len(explanation) == 0 {
		explanation = ""
	}
	return explanation, nil
}

func InputEntryYearMonth() (time.Time, time.Time, error) {

	fmt.Print("Provide year and month (YYYY-MM): ")
	yearMonthString, err := reader.ReadString('\n')
	if err != nil {
		return time.Time{}, time.Time{}, err
	}
	yearMonthString = strings.TrimSpace(yearMonthString)
	yearMonthSplitString := strings.Split(yearMonthString, "-")
	if len(yearMonthSplitString) != 2 {
		return time.Time{}, time.Time{}, errors.New("Invalid year-month format, use YYYY-MM")
	}
	year, err := strconv.Atoi(yearMonthSplitString[0])
	month, err := strconv.Atoi(yearMonthSplitString[1])
	if err != nil {
		return time.Time{}, time.Time{}, err
	}
	if year < 1900 {
		return time.Time{}, time.Time{}, errors.New("Year must be greater than or equal to 1900")
	}
	if month < 1 || month > 12 {
		return time.Time{}, time.Time{}, errors.New("Month must be between 1 and 12")
	}
	fromDate := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)
	toDate := time.Date(year, time.Month(month) + 1, 0, 0, 0, 0, 0, time.UTC)

	return fromDate, toDate, nil
}

func InputEntryChoice(entries *[]*entry.Entry, year int, month int) (*entry.Entry, error) {
	if entries == nil {
		return nil, errors.New("Entries - nil pointer dereference")
	}
	choices := make(map[int]*entry.Entry, 0)
	width := 6 + 3 + 10 + 3 + account.MaxNameLength + 3 + account.MaxNameLength + 3 + account.MaxNameLength
	title := fmt.Sprintf("Entries for %d-%d", year, month)
	var menu strings.Builder
	menu.WriteString("\n")
	menu.WriteString(strings.Repeat(" ", (width - len(title)) / 2))
	menu.WriteString(title)
	menu.WriteString("\n")
	menu.WriteString(strings.Repeat("─", 6))
	menu.WriteString("─┬─")
	menu.WriteString(strings.Repeat("─", account.MaxNameLength))
	menu.WriteString("─┬─")
	menu.WriteString(strings.Repeat("─", account.MaxNameLength))
	menu.WriteString("─┬─")
	menu.WriteString(strings.Repeat("─", account.MaxNameLength))
	menu.WriteString("\n")
	menu.WriteString("Option")
	menu.WriteString(" │ ")
	fmt.Fprintf(&menu, "%-*s", account.MaxNameLength, "Debit Account")
	menu.WriteString(" │ ")
	fmt.Fprintf(&menu, "%-*s", account.MaxNameLength, "Credit Account")
	menu.WriteString(" │ ")
	fmt.Fprintf(&menu, "%-*s", account.MaxNameLength, "Explanation")
	menu.WriteString("\n")
	menu.WriteString(strings.Repeat("─", 6))
	menu.WriteString("─┼─")
	menu.WriteString(strings.Repeat("─", account.MaxNameLength))
	menu.WriteString("─┼─")
	menu.WriteString(strings.Repeat("─", account.MaxNameLength))
	menu.WriteString("─┼─")
	menu.WriteString(strings.Repeat("─", account.MaxNameLength))
	menu.WriteString("\n")

	for c, e := range *entries {
		if e == nil {
			return nil, errors.New("Entry - nil pointer dereferencing")
		}
		choices[c + 1] = e
		menu.WriteString(strings.Repeat(" ", 2))
		fmt.Fprintf(&menu, "%-*d", 4, c + 1)
		menu.WriteString(" │ ")
		fmt.Fprintf(&menu,"%-*s", account.MaxNameLength, e.GetDebitAccount().GetName())
		menu.WriteString(" │ ")
		fmt.Fprintf(&menu, "%-*s", account.MaxNameLength, e.GetCreditAccount().GetName())
		menu.WriteString(" │ ")

		words := strings.Split(e.GetExplanation(), " ")
		var explanation strings.Builder
		for i := 0; i < len(words); {
			for exLen := 0; exLen <= account.MaxNameLength && i < len(words); {
				explanation.WriteString(words[i])
				if i < len(words) {
					i++
				}
				if i < len(words) {
					explanation.WriteString(" ")
					exLen = explanation.Len() + len(words[i])
				}
			}
			fmt.Fprintf(&menu, "%-*s\n", account.MaxNameLength, &explanation)
			if i < len(words) {
				explanation.Reset()
				menu.WriteString(strings.Repeat(" ", 6))
				menu.WriteString(" │ ")
				menu.WriteString(strings.Repeat(" ", account.MaxNameLength))
				menu.WriteString(" │ ")
				menu.WriteString(strings.Repeat(" ", account.MaxNameLength))
				menu.WriteString(" │ ")
			}
		}
	}

	menu.WriteString(strings.Repeat("─", 6))
	menu.WriteString("─┴─")
	menu.WriteString(strings.Repeat("─", account.MaxNameLength))
	menu.WriteString("─┴─")
	menu.WriteString(strings.Repeat("─", account.MaxNameLength))
	menu.WriteString("─┴─")
	menu.WriteString(strings.Repeat("─", account.MaxNameLength))
	menu.WriteString("\n")
	fmt.Print(menu.String())
	fmt.Print("Choice: ")
	input, err := reader.ReadString('\n')
	if err != nil {
		return nil, err
	}
	input = strings.TrimSpace(input)
	choice, err := strconv.Atoi(input)
	if err != nil {
		return nil, errors.New("Invalid input")
	}

	entry, exists := choices[choice]
	if exists {
		return entry, nil
	} else {
		return nil, errors.New("Invalid choice")
	}
}

func InputAccountFieldChoice() (int, error) {
	options := []string {
		"Name",
		"Description",
	}
	var menu strings.Builder
	menu.WriteString(" Choose account field\n")
	menu.WriteString(strings.Repeat("─", 2))
	menu.WriteString("─┬─")
	menu.WriteString(strings.Repeat("─", 16))
	menu.WriteString("\n")
	for option, field := range options {
		fmt.Fprintf(&menu, " %d │ %s\n", option + 1, field)
	}
	menu.WriteString(strings.Repeat("─", 2))
	menu.WriteString("─┴─")
	menu.WriteString(strings.Repeat("─", 16))
	fmt.Println(menu.String())
	fmt.Print("Choice: ")
	input, err := reader.ReadString('\n')
	fmt.Println()
	if err != nil {
		return 0, err
	}
	input = strings.TrimSpace(input)
	choice, err := strconv.Atoi(input)
	if err != nil {
		return 0, err
	}
	if choice < 1 || choice > len(options) {
		return 0, errors.New("Invalid choice")
	}
	return choice, nil
}

func InputEntryFieldChoice() (int, error) {
	options := []string {
		"Date",
		"Amount",
		"Debit Account",
		"Credit Account",
		"Explanation",
		"Posted",
	}
	var menu strings.Builder
	menu.WriteString(" Choose entry field\n")
	menu.WriteString(strings.Repeat("─", 2))
	menu.WriteString("─┬─")
	menu.WriteString(strings.Repeat("─", 14))
	menu.WriteString("\n")
	for option, field := range options {
		fmt.Fprintf(&menu, " %d │ %s\n", option + 1, field)
	}
	menu.WriteString(strings.Repeat("─", 2))
	menu.WriteString("─┴─")
	menu.WriteString(strings.Repeat("─", 14))
	fmt.Println(menu.String())
	fmt.Print("Choice: ")
	input, err := reader.ReadString('\n')
	fmt.Println()
	if err != nil {
		return 0, err
	}
	input = strings.TrimSpace(input)
	choice, err := strconv.Atoi(input)
	if err != nil {
		return 0, err
	}
	if choice < 1 || choice > len(options) {
		return 0, errors.New("Invalid choice")
	}
	return choice, nil
}

func InputEntryIsPostedStatus(isPosted bool) (bool, error) {
	if isPosted {
		fmt.Print("Change entry to not-posted? (y/n) ")
		input, err := reader.ReadString('\n')
		if err != nil {
			return false, err
		}
		input = strings.TrimSpace(input)
		if input == "y" || input == "Y" {
			return false, nil
		} else if input == "n" || input == "N" {
			return true, nil
		} else {
			return false, errors.New("Invalid input")
		}
	} else {
		fmt.Print("Post entry? (y/n) ")
		input, err := reader.ReadString('\n')
		if err != nil {
			return false, err
		}
		input = strings.TrimSpace(input)
		if input == "y" || input == "Y" {
			return true, nil
		} else if input == "n" || input == "N" {
			return false, nil
		} else {
			return false, errors.New("Invalid input")
		}
	}
}