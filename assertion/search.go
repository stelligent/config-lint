package assertion

import (
	"github.com/jmespath/go-jmespath"
)

func SearchData(expression string, data interface{}) (string, error) {
	if len(expression) == 0 {
		return "null", nil
	}
	result, err := jmespath.Search(expression, data)
	if err != nil {
		return "", err
	}
	return JSONStringify(result)
}
