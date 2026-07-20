package id

import "github.com/google/uuid"

type Id uuid.UUID

func (i Id) String() string {
	return uuid.UUID(i).String()
}

func GenerateNewId() Id {
	return Id(uuid.New())
}

func ParseString(s string) (Id, error) {
	u, err := uuid.Parse(s)
	if err != nil {
		return Id{}, err
	}
	return Id(u), nil
}