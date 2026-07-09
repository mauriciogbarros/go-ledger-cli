package id

import (
	"github.com/google/uuid"
)

type Id uuid.UUID

func NewId() Id {
	return Id(uuid.New())
}

func (i Id) String() string {
	return uuid.UUID(i).String()
}

func FromString(stringId string) (uuid.UUID, error) {
	parsedId, err := uuid.Parse(stringId)
	if err != nil {
		return uuid.Nil, err
	}
	return parsedId, nil
}