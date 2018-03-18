package assertion

import (
	"encoding/json"
	"fmt"
	"testing"
)

type TestCase struct {
	SearchResult   interface{}
	Op             string
	Value          string
	ExpectedResult bool
	Message        string
}

func getQuotesRight(jsonString string) string {
	if len(jsonString) == 0 {
		return jsonString
	}
	if jsonString[0] != '[' {
		jsonString = quoted(jsonString)
	}
	return jsonString
}

func unmarshal(s string) (interface{}, error) {
	var searchResult interface{}
	jsonString := getQuotesRight(s)
	if len(jsonString) > 0 {
		err := json.Unmarshal([]byte(jsonString), &searchResult)
		if err != nil {
			return "", err
		}
	}
	return searchResult, nil
}

func TestIsMatch(t *testing.T) {

	sliceOfTags := []string{"Foo", "Bar"}
	emptySlice := []interface{}{}

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
		{SearchResult: sliceOfTags, Op: "contains", Value: "Foo", ExpectedResult: true},
		{SearchResult: sliceOfTags, Op: "contains", Value: "Bar", ExpectedResult: true},
		{SearchResult: sliceOfTags, Op: "contains", Value: "oo", ExpectedResult: false},
		{SearchResult: "Foo", Op: "regex", Value: "o$", ExpectedResult: true},
		{SearchResult: "Foo", Op: "regex", Value: "^F", ExpectedResult: true},
		{SearchResult: "Foo", Op: "regex", Value: "^Bar$", ExpectedResult: false},
		{SearchResult: "a", Op: "lt", Value: "b", ExpectedResult: true},
		{SearchResult: "a", Op: "lt", Value: "a", ExpectedResult: false},
		{SearchResult: "a", Op: "le", Value: "a", ExpectedResult: true},
		{SearchResult: "b", Op: "le", Value: "a", ExpectedResult: false},
		{SearchResult: "b", Op: "gt", Value: "a", ExpectedResult: true},
		{SearchResult: "b", Op: "gt", Value: "b", ExpectedResult: false},
		{SearchResult: "b", Op: "ge", Value: "b", ExpectedResult: true},
		{SearchResult: "b", Op: "ge", Value: "c", ExpectedResult: false},
		{SearchResult: "", Op: "null", Value: "", ExpectedResult: true},
		{SearchResult: "1", Op: "null", Value: "", ExpectedResult: false},
		{SearchResult: "", Op: "not-null", Value: "", ExpectedResult: false},
		{SearchResult: "1", Op: "not-null", Value: "", ExpectedResult: true},
		{SearchResult: "", Op: "empty", Value: "", ExpectedResult: true},
		{SearchResult: "Foo", Op: "empty", Value: "", ExpectedResult: false},
		{SearchResult: emptySlice, Op: "empty", Value: "", ExpectedResult: true},
		{SearchResult: sliceOfTags, Op: "empty", Value: "", ExpectedResult: false},
		{SearchResult: "[\"one\",\"two\"]", Op: "intersect", Value: "[\"two\",\"three\"]", ExpectedResult: true},
		{SearchResult: "[\"one\",\"two\"]", Op: "intersect", Value: "[\"three\",\"four\"]", ExpectedResult: false},
	}
	for _, tc := range testCases {
		var b bool
		if s, isString := tc.SearchResult.(string); isString {
			searchResult, err := unmarshal(s)
			if err != nil {
				fmt.Println(err)
				t.Errorf("Unable to parse %s\n", tc.SearchResult)
			}
			b = isMatch(searchResult, tc.Op, tc.Value)
		} else {
			b = isMatch(tc.SearchResult, tc.Op, tc.Value)
		}
		if b != tc.ExpectedResult {
			t.Errorf("Expected '%s' %s '%s' to be %t", tc.SearchResult, tc.Op, tc.Value, tc.ExpectedResult)
		}
	}
}