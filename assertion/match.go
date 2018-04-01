package assertion

import (
	"fmt"
	"regexp"
	"strings"
)

type (
	// MatchResult has a true/false result, but also includes a message for better reporting
	MatchResult struct {
		Match   bool
		Message string
	}
)

func matches() (MatchResult, error) {
	return MatchResult{Match: true, Message: ""}, nil
}

func doesNotMatch(format string, args ...interface{}) (MatchResult, error) {
	return MatchResult{
		Match:   false,
		Message: fmt.Sprintf(format, args...),
	}, nil
}

func matchError(err error) (MatchResult, error) {
	return MatchResult{
		Match:   false,
		Message: err.Error(),
	}, err
}

func isMatch(data interface{}, op string, value string, valueType string) (MatchResult, error) {
	// FIXME eliminate searchResult this when all operations converted to use data
	// individual ops can call JSONStringify as needed
	searchResult, err := JSONStringify(data)
	if err != nil {
		return matchError(err)
	}
	searchResult = unquoted(searchResult)
	switch op {
	case "eq":
		if compare(data, value, valueType) == 0 {
			return matches()
		}
		return doesNotMatch("should equal to %v", value)
	case "ne":
		if compare(data, value, valueType) != 0 {
			return matches()
		}
		return doesNotMatch("should not equal to %v", value)
	case "lt":
		if compare(data, value, valueType) < 0 {
			return matches()
		}
		return doesNotMatch("should be less than %v", value)
	case "le":
		if compare(data, value, valueType) <= 0 {
			return matches()
		}
		return doesNotMatch("should be less than or equal to %v", value)
	case "gt":
		if compare(data, value, valueType) > 0 {
			return matches()
		}
		return doesNotMatch("should be greater than %v", value)
	case "ge":
		if compare(data, value, valueType) >= 0 {
			return matches()
		}
		return doesNotMatch("should be greater than or equal to %v", value)
	case "in":
		for _, v := range strings.Split(value, ",") {
			if v == searchResult {
				return matches()
			}
		}
		return doesNotMatch("should be in %v", value)
	case "not-in":
		for _, v := range strings.Split(value, ",") {
			if v == searchResult {
				return doesNotMatch("should not be in %v", value)
			}
		}
		return matches()
	case "absent":
		if isAbsent(searchResult) {
			return matches()
		}
		return doesNotMatch("should be absent")
	case "present":
		if isPresent(searchResult) {
			return matches()
		}
		return doesNotMatch("should be present")
	case "null":
		if data == nil {
			return matches()
		}
		return doesNotMatch("should be null")
	case "not-null":
		if data != nil {
			return matches()
		}
		return doesNotMatch("should not be null")
	case "empty":
		if isEmpty(data) {
			return matches()
		}
		return doesNotMatch("should be empty")
	case "not-empty":
		if !isEmpty(data) {
			return matches()
		}
		return doesNotMatch("should not be empty")
	case "intersect":
		if jsonListsIntersect(searchResult, value) {
			return matches()
		}
		return doesNotMatch("should intersect", value)
	case "contains":
		return contains(data, value)
	case "not-contains":
		return notContains(data, value)
	case "regex":
		re, err := regexp.Compile(value)
		if err != nil {
			return matchError(err)
		}
		if re.MatchString(searchResult) {
			return matches()
		}
		return doesNotMatch("should match %v", value)
	case "has-properties":
		return hasProperties(data, value)
	}
	return doesNotMatch("unknown op %v", op)
}
