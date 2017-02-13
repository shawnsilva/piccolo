package utils

import (
	"fmt"
	"testing"
)

type (
	stringTestConcat struct {
		strings []string
		answer  string
	}
)

var (
	strConcatTestData = []stringTestConcat{
		{[]string{"hello", "bob"}, "hellobob"},
		{[]string{"a12a", "bg0s"}, "a12abg0s"},
		{[]string{"What the", " space"}, "What the space"},
		{[]string{"1.", "2."}, "1.2."},
		{[]string{"how", "many", "strings", "to", "use"}, "howmanystringstouse"},
		{[]string{"What the ", "!$#", " punk"}, "What the !$# punk"},
	}
)

func TestStrConcat(t *testing.T) {
	for _, strLists := range strConcatTestData {
		concatResult := StrConcat(strLists.strings)
		if concatResult != strLists.answer {
			t.Error(
				"StrConcat of: ", strLists.strings,
				"Got: ", fmt.Sprintf(`"%s"`, concatResult),
			)
		}
	}
}
