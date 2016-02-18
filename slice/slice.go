// Package slice provides functions for working with slices
package slice

// ContainsString returns true if an identical string is found inside a slice of strings
func ContainsString(list []string, str string) bool {
	for _, v := range list {
		if v == str {
			return true
		}
	}
	return false
}
