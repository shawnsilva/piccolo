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
