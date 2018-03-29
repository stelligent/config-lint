package assertion

import (
	"strconv"
)

func compare(data interface{}, value string, valueType string) int {
	switch valueType {
	case "size":
		n, _ := strconv.Atoi(value)
		l := 0
		switch v := data.(type) {
		case []string:
			l = len(v)
		case map[string]string:
			l = len(v)
		}
		if l < n {
			return -1
		}
		if l > n {
			return 1
		}
		return 0
	default:
		tmp, _ := JSONStringify(data)
		s := unquoted(tmp)
		if s > value {
			return 1
		}
		if s < value {
			return -1
		}
		return 0
	}
}
