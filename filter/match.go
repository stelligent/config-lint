package filter

import (
	"regexp"
	"strings"
)

func isMatch(searchResult string, op string, value string) bool {
	// TODO see Cloud Custodian for ideas
	// ADD gt, ge, lt, le, not-null, empty, intersect, glob
	switch op {
	case "eq":
		if searchResult == value {
			return true
		}
	case "ne":
		if searchResult != value {
			return true
		}
	case "in":
		for _, v := range strings.Split(value, ",") {
			if v == searchResult {
				return true
			}
		}
	case "notin":
		for _, v := range strings.Split(value, ",") {
			if v == searchResult {
				return false
			}
		}
		return true
	case "absent":
		if isAbsent(searchResult) {
			return true
		}
	case "present":
		if isPresent(searchResult) {
			return true
		}
	case "contains":
		if strings.Contains(searchResult, value) {
			return true
		}
		return false
	case "regex":
		if regexp.MustCompile(value).MatchString(searchResult) {
			return true
		}
		return false
	}
	return false
}
