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

	testCases := map[string]TestCase{
		"eqTrue":                        {"Foo", "eq", "Foo", true},
		"eqFalse":                       {"Foo", "eq", "Bar", false},
		"neFalse":                       {"Foo", "ne", "Foo", false},
		"neTrue":                        {"Foo", "ne", "Bar", true},
		"inTrue":                        {"Foo", "in", "Foo,Bar,Baz", true},
		"inFalse":                       {"Foo", "in", "Bar,Baz", false},
		"notInFalse":                    {"Foo", "not-in", "Foo,Bar,Baz", false},
		"notInTrue":                     {"Foo", "not-in", "Bar,Baz", true},
		"absentFalse":                   {"Foo", "absent", "", false},
		"absentTrueForEmptyString":      {"", "absent", "", true},
		"absentTrueForNull":             {"null", "absent", "", true},
		"absentTrueForEmptyArray":       {"[]", "absent", "", true},
		"presentTrue":                   {sliceOfTags, "present", "", true},
		"presentStringTrue":             {"Foo", "present", "", true},
		"presentFalseForNil":            {nil, "present", "", false},
		"presentFalseForEmptyString":    {"", "present", "", false},
		"presentFalseForNull":           {"null", "present", "", false},
		"presentFalseForEmptyArray":     {"[]", "present", "", false},
		"containsTrueForString":         {"Foo", "contains", "oo", true},
		"containsFalseForString":        {"Foo", "contains", "aa", false},
		"containsTrueForSlice":          {sliceOfTags, "contains", "Bar", true},
		"containsFalseForSubstring":     {sliceOfTags, "contains", "oo", false},
		"regexTrueForEndOfString":       {"Foo", "regex", "o$", true},
		"regExTrueForBeginningOfString": {"Foo", "regex", "^F", true},
		"reqExFalseForEntireString":     {"Foo", "regex", "^Bar$", false},
		"ltTrue":                        {"a", "lt", "b", true},
		"ltFalse":                       {"a", "lt", "a", false},
		"leTrue":                        {"a", "le", "a", true},
		"leFalse":                       {"b", "le", "a", false},
		"gtTrue":                        {"b", "gt", "a", true},
		"gtFalse":                       {"b", "gt", "b", false},
		"geTrue":                        {"b", "ge", "b", true},
		"geFalse":                       {"b", "ge", "c", false},
		"nullTrue":                      {"", "null", "", true},
		"nullFalse":                     {"1", "null", "", false},
		"notNullFalse":                  {"", "not-null", "", false},
		"notNullTrue":                   {"1", "not-null", "", true},
		"emptyTrueForEmptyString":       {"", "empty", "", true},
		"emptyFalseForString":           {"Foo", "empty", "", false},
		"emptyTrueForEmptySlice":        {emptySlice, "empty", "", true},
		"emptyFalseForSlice":            {sliceOfTags, "empty", "", false},
		"intersectTrue":                 {"[\"one\",\"two\"]", "intersect", "[\"two\",\"three\"]", true},
		"intersectFalse":                {"[\"one\",\"two\"]", "intersect", "[\"three\",\"four\"]", false},
	}
	for k, tc := range testCases {
		var b bool
		var err error
		if s, isString := tc.SearchResult.(string); isString {
			searchResult, err := unmarshal(s)
			if err != nil {
				fmt.Println(err)
				t.Errorf("Unable to parse %s\n", tc.SearchResult)
			}
			b, err = isMatch(searchResult, tc.Op, tc.Value)
		} else {
			b, err = isMatch(tc.SearchResult, tc.Op, tc.Value)
		}
		if err != nil {
			t.Errorf("%s Failed with error: %s", k, err.Error())
		}
		if b != tc.ExpectedResult {
			t.Errorf("%s Failed Expected '%s' %s '%s' to be %t", k, tc.SearchResult, tc.Op, tc.Value, tc.ExpectedResult)
		}
	}
}
