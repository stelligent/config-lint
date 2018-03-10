package main

import (
	"testing"
)

type TestCase struct {
	SearchResult   string
	Op             string
	Value          string
	ExpectedResult bool
	Message        string
}

func TestIsMatch(t *testing.T) {

	testCases := []TestCase{
		{SearchResult: "Foo", Op: "eq", Value: "Foo", ExpectedResult: true},
		{SearchResult: "Foo", Op: "eq", Value: "Bar", ExpectedResult: false},
		{SearchResult: "Foo", Op: "ne", Value: "Foo", ExpectedResult: false},
		{SearchResult: "Foo", Op: "ne", Value: "Bar", ExpectedResult: true},
		{SearchResult: "Foo", Op: "in", Value: "Foo,Bar,Baz", ExpectedResult: true},
		{SearchResult: "Foo", Op: "in", Value: "Bar,Baz", ExpectedResult: false},
		{SearchResult: "Foo", Op: "notin", Value: "Foo,Bar,Baz", ExpectedResult: false},
		{SearchResult: "Foo", Op: "notin", Value: "Bar,Baz", ExpectedResult: true},
		{SearchResult: "Foo", Op: "absent", Value: "", ExpectedResult: false},
		{SearchResult: "", Op: "absent", Value: "", ExpectedResult: true},
		{SearchResult: "null", Op: "absent", Value: "", ExpectedResult: true},
		{SearchResult: "[]", Op: "absent", Value: "", ExpectedResult: true},
		{SearchResult: "Foo", Op: "present", Value: "", ExpectedResult: true},
		{SearchResult: "", Op: "present", Value: "", ExpectedResult: false},
		{SearchResult: "null", Op: "present", Value: "", ExpectedResult: false},
		{SearchResult: "[]", Op: "present", Value: "", ExpectedResult: false},
		{SearchResult: "Foo", Op: "contains", Value: "oo", ExpectedResult: true},
		{SearchResult: "Foo", Op: "contains", Value: "aa", ExpectedResult: false},
		{SearchResult: "Foo", Op: "regex", Value: "o$", ExpectedResult: true},
		{SearchResult: "Foo", Op: "regex", Value: "^F", ExpectedResult: true},
		{SearchResult: "Foo", Op: "regex", Value: "^Bar$", ExpectedResult: false},
	}
	for _, tc := range testCases {
		b := isMatch(tc.SearchResult, tc.Op, tc.Value)
		if b != tc.ExpectedResult {
			t.Errorf("Expected '%s' %s '%s' to be %t", tc.SearchResult, tc.Op, tc.Value, tc.ExpectedResult)
		}
	}
}
