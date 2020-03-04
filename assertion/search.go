package assertion

import (
	"github.com/jmespath/go-jmespath"
)

// SearchData applies a JMESPath to a JSON object
func SearchData(expression string, data interface{}) (interface{}, error) {
	if len(expression) == 0 {
		return "null", nil
	}

	return jmespath.Search(expression, data)
}
