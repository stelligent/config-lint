package assertion

import (
	"github.com/jmespath/go-jmespath"
)

func SearchData(expression string, data interface{}) (interface{}, error) {
	if len(expression) == 0 {
		return "null", nil
	}
	return jmespath.Search(expression, data)
}
