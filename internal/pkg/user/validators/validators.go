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
	return !strings.Contains(username, "admin")
}

func ValidatePassword(password string) bool {
	if len(password) < 8 {
		return false
	}
	var digit, upper, lower bool
	for _, s := range password {
		switch {
		case unicode.IsDigit(s):
			digit = true
		case unicode.IsUpper(s):
			upper = true
		case unicode.IsLower(s):
			lower = true
		}
	}
	return digit && upper && lower
}
