package assertion

import (
	"strings"
)

func contains(data interface{}, value string) (bool, error) {
	if c, isSlice := convertToSlice(data); isSlice {
		for _, element := range c {
			if stringElement, isString := element.(string); isString {
				if stringElement == value {
					return true, nil
				}
			}
		}
		return false, nil
	}
	if s, isString := convertToString(data); isString {
		if strings.Contains(s, value) {
			return true, nil
		}
		return false, nil
	}
	searchResult, err := JSONStringify(data)
	if err != nil {
		return false, err
	}
	return strings.Contains(searchResult, value), nil
}

func notContains(data interface{}, value string) (bool, error) {
	b, err := contains(data, value)
	if err != nil {
		return false, err
	}
	return !b, nil
}
