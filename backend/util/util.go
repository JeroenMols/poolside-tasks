package util

import "github.com/google/uuid"

type GenerateUuid func() string

func GenerateRandomUuid() string {
	return uuid.New().String()
}
