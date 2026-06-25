package utils

import "regexp"

// EmailRegex is a simple email format validator
var EmailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)

// IsValidEmail checks if an email address has a valid format
func IsValidEmail(email string) bool {
	return EmailRegex.MatchString(email)
}
