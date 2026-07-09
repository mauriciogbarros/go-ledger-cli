package storage

import "go.mod/internal/journal"

var Journal *journal.Journal

func GetJournal() journal.Journal {
	return *Journal
}

func SetJournal(j *journal.Journal) {
	Journal = j
}