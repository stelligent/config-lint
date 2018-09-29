package linter

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

type interpolationTestCase struct {
	Input    string
	Expected interface{}
}

func TestInterpolation(t *testing.T) {
	testCases := []interpolationTestCase{
		{"${2+6}", "8"},
		{"bucket_${var.environment}", "bucket_development"},
		{"${var.environment == \"development\" ? \"YES\" : \"NO\"}", "YES"},
		{"${missing_func(1)}", ""},
		{"${local.count + local.count}", "202"},
		{"${replace(var.template,var.user_pattern,var.user)}", "https://users/adam"},
		{"${list(var.a, var.b, var.c)}", []interface{}{"one", "two", "three"}},
		{"${element(list(var.a, var.b, var.c),2)}", "three"},
		{"${join(var.pipe, list(var.a, var.b))}", "one|two"},
		{"${concat(list(var.a,var.b), list(var.c))}", []interface{}{"one", "two", "three"}},
		{"${format(\"id-%s\",var.a)}", "id-one"},
		{"${map(var.k1,var.v1,var.k2,var.v2)}", map[string]interface{}{"key1": "value1", "key2": "value2"}},
	}
	vars := []Variable{
		{Name: "var.environment", Value: "development"},
		{Name: "local.count", Value: "101"},
		{Name: "var.template", Value: "https://users/USER_ID"},
		{Name: "var.user_pattern", Value: "USER_ID"},
		{Name: "var.user", Value: "adam"},
		{Name: "var.a", Value: "one"},
		{Name: "var.b", Value: "two"},
		{Name: "var.c", Value: "three"},
		{Name: "var.pipe", Value: "|"},
		{Name: "var.k1", Value: "key1"},
		{Name: "var.k2", Value: "key2"},
		{Name: "var.v1", Value: "value1"},
		{Name: "var.v2", Value: "value2"},
	}
	for _, tc := range testCases {
		result := interpolate(tc.Input, vars)
		assert.Equal(t, tc.Expected, result)
	}
}
