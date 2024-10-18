package utils

import (
	"github.com/google/uuid"
	"strings"
)

var (
	UuidSeparator = "-"
)

func UUID2STR() string {
	return strings.ReplaceAll(uuid.New().String(), UuidSeparator, "")
}

func UUID2INT() int64 {
	return int64(uuid.New().ID())
}
