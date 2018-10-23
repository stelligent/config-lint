package assertion

import (
	"strings"
)

func interfaceListContains(v []interface{}, key, value string) (MatchResult, error) {
	for _, element := range v {
		if stringElement, isString := element.(string); isString {
			if stringElement == value {
				return matches()
			}
			if strings.Contains(stringElement, value) {
				return matches()
			}
		}
	}
	return doesNotMatch("%v does not contain %v", key, value)
}

func stringListContains(v []string, key, value string) (MatchResult, error) {
	for _, stringElement := range v {
		if stringElement == value {
			return matches()
		}
		if strings.Contains(stringElement, value) {
			return matches()
		}
	}
	return doesNotMatch("%v does not contain %v", key, value)
}

func stringContains(v string, key, value string) (MatchResult, error) {
	if strings.Contains(v, value) {
		return matches()
	}
	return doesNotMatch("%v does not contain %v", key, value)
}

func defaultContains(data interface{}, key, value string) (MatchResult, error) {
	searchResult, err := JSONStringify(data)
	if err != nil {
		return matchError(err)
	}
	if strings.Contains(searchResult, value) {
		return matches()
	}
	return doesNotMatch("%v does not contain %v", key, value)
}

func contains(data interface{}, key, value string) (MatchResult, error) {
	switch v := data.(type) {
	case []interface{}:
		return interfaceListContains(v, key, value)
	case []string:
		return stringListContains(v, key, value)
	case string:
		return stringContains(v, key, value)
	default:
		return defaultContains(v, key, value)
	}
}

func doesNotContain(data interface{}, key, value string) (MatchResult, error) {
	m, err := contains(data, key, value)
	if err != nil {
		return matchError(err)
	}
	if m.Match {
		return doesNotMatch("%v should not contain %v", key, value)
	}
	return matches()
}

func startsWith(data interface{}, key, prefix string) (MatchResult, error) {
	switch v := data.(type) {
	case string:
		if strings.HasPrefix(v, prefix) {
			return matches()
		}
		return doesNotMatch("%v does not start with %v", key, prefix)
	default:
		return doesNotMatch("%v is not a string %v", key, prefix)
	}
}

func endsWith(data interface{}, key, suffix string) (MatchResult, error) {
	switch v := data.(type) {
	case string:
		if strings.HasSuffix(v, suffix) {
			return matches()
		}
		return doesNotMatch("%v does not end with %v", key, suffix)
	default:
		return doesNotMatch("%v is not a string %v", key, suffix)
	}
}
