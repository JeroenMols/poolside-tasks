package util

import (
	"fmt"
	"github.com/lithammer/shortuuid/v4"
	"time"
)

// GenerateUuid TODO: support short human readable uuids e.g. list_abh135asdfjkl
type GenerateUuid func(string) string

func GenerateRandomUuid(prefix string) string {
	return fmt.Sprintf("%s_%s", prefix, shortuuid.New())
}

type CurrentTime func() time.Time

func GetCurrentTime() time.Time {
	return time.Now()
}
