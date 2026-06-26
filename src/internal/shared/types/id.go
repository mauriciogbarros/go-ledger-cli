package types

import (
	"github.com/google/uuid"
)

type Id string

func NewID() Id {
	return Id(uuid.New().String())
}