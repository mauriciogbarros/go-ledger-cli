package transaction

import (
	"time"

	"go.mod/internal/shared/types"
)

type TransactionType string

const (
	Deposit TransactionType = "deposit"
	Withdrawal TransactionType = "withdrawal"
	Transfer TransactionType = "transfer"
)

type Transaction struct {
	ID types.Id
	Type TransactionType
	Amount types.Currency
	From types.Id
	To types.Id
	Timestamp time.Time
}