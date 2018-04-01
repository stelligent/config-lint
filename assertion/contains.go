package assertion

import (
	"strings"
)

func contains(data interface{}, value string) (MatchResult, error) {
	switch v := data.(type) {
	case []interface{}:
		for _, element := range v {
			if stringElement, isString := element.(string); isString {
				if stringElement == value {
					return matches()
				}
			}
		}
		return doesNotMatch("does not contain %v", value)
	case []string:
		for _, stringElement := range v {
			if stringElement == value {
				return matches()
			}
		}
		return doesNotMatch("does not contain %v", value)
	case string:
		if strings.Contains(v, value) {
			return matches()
		}
		return doesNotMatch("does not contain %v", value)
	default:
		searchResult, err := JSONStringify(data)
		if err != nil {
			return matches()
		}
		if strings.Contains(searchResult, value) {
			return matches()
		}
		return doesNotMatch("does not contain %v", value)
	}
}

func notContains(data interface{}, value string) (MatchResult, error) {
	m, err := contains(data, value)
	if err != nil {
		return matchError(err)
	}
	if m.Match {
		return doesNotMatch("should not contain %v", value)
	}
	return matches()
}
