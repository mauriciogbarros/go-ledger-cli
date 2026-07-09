package main

import (
	"go.mod/db"
)

// "go.mod/internal/cli"
// "go.mod/internal/journal"
// "go.mod/internal/storage"

func main() {
	// var j journal.Journal
	// j.NewJournal()
	// storage.SetJournal(&j)

	// cli.Run()
	db.Initialize()
}
