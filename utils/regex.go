package utils

import (
	"regexp"
)

// REGetNamedGroupsResults takes a Regexp with named groups, a search string,
// and returns a map of the found search groups, with their values. The values
// may be empty.
func REGetNamedGroupsResults(re *regexp.Regexp, searchStr string) map[string]string {
	match := re.FindStringSubmatch(searchStr)
	result := make(map[string]string)
	for i, name := range re.SubexpNames() {
		if i != 0 {
			result[name] = match[i]
		}
	}
	return result
}
