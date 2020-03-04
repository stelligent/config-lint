package assertion

import (
	"github.com/jmespath/go-jmespath"
)

// SearchData applies a JMESPath to a JSON object
func SearchData(expression string, data interface{}) (interface{}, error) {
	if len(expression) == 0 {
		return "null", nil
	}

	Debugf("Search Expression: %s\n", expression)
	// Replace invalid chars for jmespath search expression
	//jmespath_expr := strings.ReplaceAll(expression, ":", strconv.Quote(":"))
	//jmespath_expr := expression
	//if strings.Contains(expression, ":") {
	//	jmespath_expr = strconv.Quote(expression)
	//	Debugf("JMESPath Search Expression: %s\n", jmespath_expr)
	//}

	//return jmespath.Search(jmespath_expr, data)
	return jmespath.Search(expression, data)
}
