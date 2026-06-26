package account

import types "go.mod/internal/shared/types"

type Account struct {
	ID types.Id
	Name string
	balance types.Currency
}