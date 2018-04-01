package assertion

import (
	"strings"
)

func contains(data interface{}, value string) (bool, error) {
	switch v := data.(type) {
	case []interface{}:
		for _, element := range v {
			if stringElement, isString := element.(string); isString {
				if stringElement == value {
					return true, nil
				}
			}
		}
		return false, nil
	case []string:
		for _, stringElement := range v {
			if stringElement == value {
				return true, nil
			}
		}
		return false, nil
	case string:
		if strings.Contains(v, value) {
			return true, nil
		}
		return false, nil
	default:
		searchResult, err := JSONStringify(data)
		if err != nil {
			return false, err
		}
		return strings.Contains(searchResult, value), nil
	}
}

func notContains(data interface{}, value string) (bool, error) {
	b, err := contains(data, value)
	if err != nil {
		return false, err
	}
	return !b, nil
}
