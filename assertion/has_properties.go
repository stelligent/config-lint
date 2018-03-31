package assertion

import (
	"strings"
)

func hasProperties(data interface{}, list string) (bool, error) {
	for _, key := range strings.Split(list, ",") {
		if m, ok := data.(map[string]interface{}); ok {
			if _, ok := m[key]; !ok {
				return false, nil
			}
		}
	}
	return true, nil
}
