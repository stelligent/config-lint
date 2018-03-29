package assertion

import (
	"strconv"
)

func intCompare(n1 int, n2 int) int {
	if n1 < n2 {
		return -1
	}
	if n1 > n2 {
		return 1
	}
	return 0
}

func compare(data interface{}, value string, valueType string) int {
	switch valueType {
	case "size":
		n, _ := strconv.Atoi(value)
		l := 0
		switch v := data.(type) {
		case []interface{}:
			l = len(v)
		case map[string]interface{}:
			l = len(v)
		}
		return intCompare(l, n)
	case "integer":
		n1, _ := strconv.Atoi(data.(string))
		n2, _ := strconv.Atoi(value)
		return intCompare(n1, n2)
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
