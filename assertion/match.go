package assertion

import (
	"fmt"
	"regexp"
	"strings"
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

func isMatch(data interface{}, expression Expression) (MatchResult, error) {
	// FIXME eliminate searchResult this when all operations converted to use data
	// individual ops can call JSONStringify as needed
	searchResult, err := JSONStringify(data)
	if err != nil {
		return matchError(err)
	}
	searchResult = unquoted(searchResult)
	key := expression.Key
	op := expression.Op
	value := expression.Value
	valueType := expression.ValueType

	switch op {
	case "eq":
		if compare(data, value, valueType) == 0 {
			return matches()
		}
		return doesNotMatch("%v(%v) should equal to %v", key, searchResult, value)
	case "ne":
		if compare(data, value, valueType) != 0 {
			return matches()
		}
		return doesNotMatch("%v(%v) should not be equal to %v", key, searchResult, value)
	case "lt":
		if compare(data, value, valueType) < 0 {
			return matches()
		}
		return doesNotMatch("%v(%v) should be less than %v", key, searchResult, value)
	case "le":
		if compare(data, value, valueType) <= 0 {
			return matches()
		}
		return doesNotMatch("%v(%v) should be less than or equal to %v", key, searchResult, value)
	case "gt":
		if compare(data, value, valueType) > 0 {
			return matches()
		}
		return doesNotMatch("%v(%v) should be greater than %v", key, searchResult, value)
	case "ge":
		if compare(data, value, valueType) >= 0 {
			return matches()
		}
		return doesNotMatch("%v(%v) should be greater than or equal to %v", key, searchResult, value)
	case "in":
		for _, v := range strings.Split(value, ",") {
			if v == searchResult {
				return matches()
			}
		}
		return doesNotMatch("%v(%v) should be in %v", key, searchResult, value)
	case "not-in":
		for _, v := range strings.Split(value, ",") {
			if v == searchResult {
				return doesNotMatch("%v(%v) should not be in %v", key, searchResult, value)
			}
		}
		return matches()
	case "absent":
		if isAbsent(searchResult) {
			return matches()
		}
		return doesNotMatch("%v should be absent", key)
	case "present":
		if isPresent(searchResult) {
			return matches()
		}
		return doesNotMatch("%v should be present", key)
	case "null":
		if data == nil {
			return matches()
		}
		return doesNotMatch("%v should be null", key)
	case "not-null":
		if data != nil {
			return matches()
		}
		return doesNotMatch("%v should not be null", key)
	case "empty":
		if isEmpty(data) {
			return matches()
		}
		return doesNotMatch("%v should be empty", key)
	case "not-empty":
		if !isEmpty(data) {
			return matches()
		}
		return doesNotMatch("%v should not be empty", key)
	case "intersect":
		if jsonListsIntersect(searchResult, value) {
			return matches()
		}
		return doesNotMatch("%v should intersect with %v", key, value)
	case "contains":
		return contains(data, key, value)
	case "not-contains":
		return notContains(data, key, value)
	case "regex":
		re, err := regexp.Compile(value)
		if err != nil {
			return matchError(err)
		}
		if re.MatchString(searchResult) {
			return matches()
		}
		return doesNotMatch("%v should match %v", key, value)
	case "has-properties":
		return hasProperties(data, value)
	case "is-true":
		if searchResult == "true" {
			return matches()
		}
		return doesNotMatch("%v should be 'true', not '%v'", key, value)
	case "is-false":
		if searchResult == "false" {
			return matches()
		}
		return doesNotMatch("%v should be 'false', not '%v'", key, value)
	}
	return doesNotMatch("unknown op %v", op)
}
