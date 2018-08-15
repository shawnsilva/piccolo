package utils

import (
	"reflect"
	"regexp"
	"testing"
)

type (
	REGetNamedGroupsResultsTest struct {
		regex        string
		searchString string
		answer       map[string]string
	}
)

var (
	REGetNamedGroupsResultsTestData = []REGetNamedGroupsResultsTest{
		{
			regex:        `^(?:v(?P<version>\d\.\d\.\d)(?:rc(?P<rc>\d+))?)?(?:-(?P<numcommits>[0-9]*))?(?:(?:-?g)?(?P<commit>[0-9a-f]{6,40}))?(?P<dirty>-dirty)?$`,
			searchString: "v1.0.0",
			answer:       map[string]string{"version": "1.0.0", "rc": "", "numcommits": "", "commit": "", "dirty": ""},
		},
		{
			regex:        `^(?:v(?P<version>\d\.\d\.\d)(?:rc(?P<rc>\d+))?)?(?:-(?P<numcommits>[0-9]*))?(?:(?:-?g)?(?P<commit>[0-9a-f]{6,40}))?(?P<dirty>-dirty)?$`,
			searchString: "v1.0.0rc1-2-gaef132-dirty",
			answer:       map[string]string{"version": "1.0.0", "rc": "1", "numcommits": "2", "commit": "aef132", "dirty": "-dirty"},
		},
	}
)

func TestREGetNamedGroupsResults(t *testing.T) {
	for _, retest := range REGetNamedGroupsResultsTestData {
		re := regexp.MustCompile(retest.regex)
		reResult := REGetNamedGroupsResults(re, retest.searchString)
		eq := reflect.DeepEqual(retest.answer, reResult)
		if !eq {
			t.Error(
				"Results Didnt match",
				"Answer:", retest.answer,
				"Result:", reResult,
			)
		}
	}
}
