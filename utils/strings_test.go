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

	strJoinTest struct {
		strings   []string
		separator string
		answer    string
	}

	stringInSliceTest struct {
		strToFind     string
		listOfStrings []string
		answer        bool
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

	strJoinTestData = []strJoinTest{
		{[]string{"hello", "bob"}, " ", "hello bob"},
		{[]string{"a12a", "bg0s"}, " ", "a12a bg0s"},
		{[]string{"What the", " space"}, ",", "What the, space"},
		{[]string{"1.", "2."}, " ", "1. 2."},
		{[]string{"how", "many", "strings", "to", "use"}, "+", "how+many+strings+to+use"},
		{[]string{"What the ", "!$#", " punk"}, "^", "What the ^!$#^ punk"},
	}

	stringInSliceTestData = []stringInSliceTest{
		{"bill", []string{"jim", "george", "alice", "bill"}, true},
		{"Aklc(*)&", []string{"jim", "george", "alice", "bill"}, false},
		{"Aklc(*)&", []string{"jim", "george", "alice", "bill", "Aklc(*)&"}, true},
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

func TestStrJoin(t *testing.T) {
	for _, testData := range strJoinTestData {
		joinResult := StrJoin(testData.strings, testData.separator)
		if joinResult != testData.answer {
			t.Error(
				"StrJoin of: ", testData.strings,
				"Got: ", fmt.Sprintf(`"%s"`, joinResult),
			)
		}
	}
}

func TestStringInSlice(t *testing.T) {
	for _, testData := range stringInSliceTestData {
		findResult := StringInSlice(testData.strToFind, testData.listOfStrings)
		if findResult != testData.answer {
			t.Error(
				"StringInSlice of string: ", testData.strToFind,
				"In Slice: ", testData.listOfStrings,
				"Got: ", findResult,
			)
		}
	}
}
