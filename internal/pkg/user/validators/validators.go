package validators

import (
	"regexp"
	"strings"
	"unicode"
)

func ValidateUsername(username string) bool {
	username = strings.ToLower(username)

	checker, _ := regexp.Compile("^[A-Za-z0-9_]{3,}")
	matchedLen := len(checker.FindString(username))

	if matchedLen == 0 || matchedLen < len(username) {
		return false
	}
	if strings.Contains(username, "admin") || strings.Contains(username, "postmaster") {
		return false
	}
	return true
}

func ValidatePassword(password string) bool {
	if len(password) < 6 {
		return false
	}
	var digit, letter bool
	for _, s := range password {
		switch {
		case unicode.IsDigit(s):
			digit = true
		case unicode.IsLetter(s):
			letter = true
		}
	}
	return digit && letter
}
