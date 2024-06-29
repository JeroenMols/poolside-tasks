package util

import "github.com/google/uuid"

// TODO: support short human readible uuids e.g. list_abh135asdfjkl
type GenerateUuid func() string

func GenerateRandomUuid() string {
	return uuid.New().String()
}
