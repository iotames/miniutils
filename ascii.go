package miniutils

import (
	"strings"
)

// Copy from: "net/http/internal/ascii", print.go

func asciiIsPrint(s string) bool {
	for i := 0; i < len(s); i++ {
		if s[i] < ' ' || s[i] > '~' {
			return false
		}
	}
	return true
}

// ToLower returns the lowercase version of s if s is ASCII and printable.
func asciiToLower(s string) (lower string, ok bool) {
	if !asciiIsPrint(s) {
		return "", false
	}
	return strings.ToLower(s), true
}
