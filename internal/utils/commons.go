package utils

import (
	"strings"
)

func IsBlank(value string) bool {
	return len(strings.TrimSpace(value)) == 0
}
