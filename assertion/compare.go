package assertion

import (
	"strconv"
	"time"
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

func daysOld(data interface{}) int {
	if stringValue, ok := data.(string); ok {
		layout := "2006-01-02T15:04:05Z"
		t, err := time.Parse(layout, stringValue)
		if err != nil {
			return 0
		}
		days := int(time.Since(t).Hours() / 24.0)
		Debugf("Date: %v Days ago: %d\n", data, days)
		return days
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
		switch v := data.(type) {
		case float64:
			n1 := int(v)
			n2, _ := strconv.Atoi(value)
			return intCompare(n1, n2)
		case int:
			n2, _ := strconv.Atoi(value)
			return intCompare(v, n2)
		case string:
			n1, _ := strconv.Atoi(v)
			n2, _ := strconv.Atoi(value)
			return intCompare(n1, n2)
		}
		return 0
	case "age":
		n, _ := strconv.Atoi(value)
		return intCompare(daysOld(data), n)
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
