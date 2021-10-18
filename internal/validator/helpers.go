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

// Unique() checks if a string slice is a set and returns bool
func Unique(ss []string) bool {
	tmp := make(map[string]bool)
	for _, v := range ss {
		tmp[v] = true
	}
	return len(tmp) == len(ss)
}

// Matches() checks for regex pattern and returns bool
func Matches(s string, rx *regexp.Regexp) bool {
	return rx.MatchString(s)
}
