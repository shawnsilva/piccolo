package utils

import (
	"bytes"
)

// StrConcat takes an array of strings, and concatenates them together, returning
// one string.
func StrConcat(strings []string) string {
	var buffer bytes.Buffer
	for _, v := range strings {
		buffer.WriteString(v)
	}
	return buffer.String()
}

// StringInSlice takes a string and a slice and returns true if the string is
// in the slice, false otherwise
func StringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}
