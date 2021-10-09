package utils

import (
	valid "github.com/asaskevich/govalidator"
	"strings"
)

func IsBlank(value string) bool {
	return len(strings.TrimSpace(value)) == 0
}

func IsNumeric(value string) bool {
	return valid.IsNumeric(value)
}

func IsInt(value string) bool {
	return valid.IsInt(value)
}