package validator

import "regexp"

var (
	EmailRX = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
)

// In() checks if a string in a string slice and returns bool
func In(s string, list ...string) bool {
	for _, v := range list {
		if v == s {
			return true
		}
	}
	return false
}

// IsUniqueSS() checks if a string slice is a set and returns bool
func IsUniqueSS(ss []string) bool {
	cache := make(map[string]bool)
	for _, v := range ss {
		cache[v] = true
	}
	return len(cache) == len(ss)
}

// IsUniqueIS() checks if a int64 slice is a set and returns bool
func IsUniqueIS(is []int64) bool {
	cache := make(map[int64]bool)
	for _, v := range is {
		cache[v] = true
	}
	return len(cache) == len(is)
}

// Matches() checks for regex pattern and returns bool
func Matches(s string, rx *regexp.Regexp) bool {
	return rx.MatchString(s)
}

func ValidateEmail(v *Validator, email string) {
	v.Check(email != "", "email", "must be provided")
	v.Check(Matches(email, EmailRX), "email", "must be a valid email address")
}

func ValidateTokenPlaintext(v *Validator, tokenPlaintext string) {
	v.Check(tokenPlaintext != "", "token", "must be provided")
	v.Check(len(tokenPlaintext) == 26, "token", "must be 26 bytes long")
}
