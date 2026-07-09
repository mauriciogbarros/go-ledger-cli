package post

import (
	"time"

	"go.mod/internal/currency"
	"go.mod/internal/id"
)

type Post struct {
	id id.Id
	date time.Time
	explanation string
	debitAmount currency.Currency
	creditAmount currency.Currency
	balance currency.Currency
}