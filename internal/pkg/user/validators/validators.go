package validators

import (
	"regexp"
	"strings"
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
