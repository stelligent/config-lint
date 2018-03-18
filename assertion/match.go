package assertion

import (
	"regexp"
	"strings"
)

func isNil(data interface{}) bool {
	return data == nil
}

func isString(data interface{}) bool {
	_, ok := data.(string)
	return ok
}

func convertToString(data interface{}) (string, bool) {
	s, ok := data.(string)
	return s, ok
}

func convertToSliceOfStrings(data interface{}) ([]string, bool) {
	s, ok := data.([]string)
	return s, ok
}

func isObject(data interface{}) bool {
	_, ok := data.(map[string]interface{})
	return ok
}

func isMatch(data interface{}, op string, value string) bool {
	searchResult, err := JSONStringify(data)
	if err != nil {
		panic(err)
	}
	searchResult = unquoted(searchResult)
	switch op {
	case "eq":
		if searchResult == value {
			return true
		}
	case "ne":
		if searchResult != value {
			return true
		}
	case "lt":
		if searchResult < value {
			return true
		}
	case "le":
		if searchResult <= value {
			return true
		}
	case "gt":
		if searchResult > value {
			return true
		}
	case "ge":
		if searchResult >= value {
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
	case "null":
		return isNil(data)
	case "not-null":
		return !isNil(data)
	case "empty":
		if isEmpty(searchResult) {
			return true
		}
	case "intersect":
		if jsonListsIntersect(searchResult, value) {
			return true
		}
	case "contains":
		if s, isString := convertToString(data); isString {
			if strings.Contains(s, value) {
				return true
			}
		}
		if c, isSlice := convertToSliceOfStrings(data); isSlice {
			for _, element := range c {
				if element == value {
					return true
				}
			}
			return false
		}
		return strings.Contains(searchResult, value)
	case "regex":
		if regexp.MustCompile(value).MatchString(searchResult) {
			return true
		}
		return false
	}
	return false
}
