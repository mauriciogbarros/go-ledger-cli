package cli

import (
	"fmt"
	"os"
	"strconv"

	"go.mod/internal/entry"
	"go.mod/internal/storage"
)

func Run() {
	args := os.Args[1:]
	if len(args) == 0 {
		fmt.Println("Usage: ledger <command> [args]")
		return
	}

	switch args[0] {
	case "create-account":
		if err := storage.AddAccount(); err != nil {
			fmt.Println("Error:", err)
			return
		}
		accounts := storage.GetAccounts()
		fmt.Printf("✅ Account created: acc_%d\n", len(accounts))

	case "list-accounts":
		accounts := storage.GetAccounts()
		if len(accounts) == 0 {
			fmt.Println("No accounts found.")
			return
		}
		for i, a := range accounts {
			fmt.Printf("acc_%d — %s (%s)\n", i+1, a.Name, a.Type)
		}

	case "deposit":
		if len(args) < 3 {
			fmt.Println("Usage: ledger deposit <account_id> <amount>")
			return
		}
		amount, err := strconv.ParseFloat(args[2], 64)
		if err != nil || amount <= 0 {
			fmt.Println("Invalid amount.")
			return
		}
		storage.Journal.AddEntry(entry.NewEntry(args[1], "Cash", amount, "Deposit"))
		fmt.Printf("✅ Deposited %.2f to %s\n", amount, args[1])

	case "withdraw":
		if len(args) < 3 {
			fmt.Println("Usage: ledger withdraw <account_id> <amount>")
			return
		}
		amount, err := strconv.ParseFloat(args[2], 64)
		if err != nil || amount <= 0 {
			fmt.Println("Invalid amount.")
			return
		}
		storage.Journal.AddEntry(entry.NewEntry("Cash", args[1], amount, "Withdrawal"))
		fmt.Printf("✅ Withdrew %.2f from %s\n", amount, args[1])

	case "transfer":
		if len(args) < 4 {
			fmt.Println("Usage: ledger transfer <from> <to> <amount>")
			return
		}
		amount, err := strconv.ParseFloat(args[3], 64)
		if err != nil || amount <= 0 {
			fmt.Println("Invalid amount.")
			return
		}
		storage.Journal.AddEntry(entry.NewEntry(args[2], args[1], amount, "Transfer"))
		fmt.Printf("✅ Transferred %.2f from %s to %s\n", amount, args[1], args[2])

	case "balance":
		if len(args) < 2 {
			fmt.Println("Usage: ledger balance <account_id>")
			return
		}
		fmt.Printf("💰 Balance: (not yet implemented — requires ledger posting)\n")

	case "history":
		if len(args) < 2 {
			fmt.Println("Usage: ledger history <account_id>")
			return
		}
		fmt.Println("📜 Transactions:")
		fmt.Println(storage.GetJournal())

	default:
		fmt.Printf("Unknown command: %s\n", args[0])
	}
}
