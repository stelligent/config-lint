package main

import (
	"encoding/json"
	"github.com/jmespath/go-jmespath"
)

func searchData(expression string, data interface{}) string {
	if len(expression) == 0 {
		return "null"
	}
	result, err := jmespath.Search(expression, data)
	if err != nil {
		panic(err)
	}
	toJSON, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		panic(err)
	}
	return string(toJSON)
}
