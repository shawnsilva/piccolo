package utils

import (
	"bytes"
)

func StrConcat(strings []string) string {
	var buffer bytes.Buffer
	for _, v := range strings {
		buffer.WriteString(v)
	}
	return buffer.String()
}
