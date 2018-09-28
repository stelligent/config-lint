package linter

import (
	"testing"
)

type interpolationTestCase struct {
	Input    string
	Expected string
}

func TestInterpolation(t *testing.T) {
	testCases := []interpolationTestCase{
		{"${2+6}", "8"},
		{"bucket_${var.environment}", "bucket_development"},
		{"${var.environment == \"development\" ? \"YES\" : \"NO\"}", "YES"},
		{"${missing_func(1)}", ""},
		{"${local.count + local.count}", "202"},
	}
	vars := []Variable{
		{Name: "var.environment", Value: "development"},
		{Name: "local.count", Value: "101"},
	}
	for _, tc := range testCases {
		result := interpolate(tc.Input, vars)
		if result != tc.Expected {
			t.Errorf("Expected %s returned %s instead of %s", tc.Input, result, tc.Expected)
		}
	}
}
