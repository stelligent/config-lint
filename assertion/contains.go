package assertion

import (
	"strings"
)

func contains(data interface{}, key, value string) (MatchResult, error) {
	switch v := data.(type) {
	case []interface{}:
		for _, element := range v {
			if stringElement, isString := element.(string); isString {
				if stringElement == value {
					return matches()
				}
			}
		}
		return doesNotMatch("%v does not contain %v", key, value)
	case []string:
		for _, stringElement := range v {
			if stringElement == value {
				return matches()
			}
		}
		return doesNotMatch("%v does not contain %v", key, value)
	case string:
		if strings.Contains(v, value) {
			return matches()
		}
		return doesNotMatch("%v does not contain %v", key, value)
	default:
		searchResult, err := JSONStringify(data)
		if err != nil {
			return matches()
		}
		if strings.Contains(searchResult, value) {
			return matches()
		}
		return doesNotMatch("%v does not contain %v", key, value)
	}
}

func notContains(data interface{}, key, value string) (MatchResult, error) {
	m, err := contains(data, key, value)
	if err != nil {
		return matchError(err)
	}
	if m.Match {
		return doesNotMatch("%v should not contain %v", key, value)
	}
	return matches()
}
