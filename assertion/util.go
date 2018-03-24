package assertion

import (
	"encoding/json"
	"fmt"
	"path/filepath"
)

func convertToSlice(data interface{}) ([]interface{}, bool) {
	s, ok := data.([]interface{})
	return s, ok
}

func unquoted(s string) string {
	if s[0] == '"' {
		return s[1 : len(s)-1]
	}
	return s
}

func quoted(s string) string {
	return fmt.Sprintf("\"%s\"", s)
}

func isAbsent(s string) bool {
	if s == "" || s == "null" || s == "[]" {
		return true
	}
	return false
}

func isPresent(s string) bool {
	return !isAbsent(s)
}

func isNotNull(s string) bool {
	return s != "null"
}

func isEmpty(data interface{}) bool {
	if data == nil {
		return true
	}
	if s, isString := convertToString(data); isString {
		return len(s) == 0
	}
	if c, isSlice := convertToSlice(data); isSlice {
		return len(c) == 0
	}
	return false
}

func listsIntersect(list1 []string, list2 []string) bool {
	for _, a := range list1 {
		for _, b := range list2 {
			if a == b {
				return true
			}
		}
	}
	return false
}

func jsonListsIntersect(s1 string, s2 string) bool {
	var a1 []string
	var a2 []string
	err := json.Unmarshal([]byte(s1), &a1)
	if err != nil {
		return false
	}
	err = json.Unmarshal([]byte(s2), &a2)
	if err != nil {
		return false
	}
	return listsIntersect(a1, a2)
}

// ShouldIncludeFile return true if a filename matches one of a list of patterns
func ShouldIncludeFile(patterns []string, filename string) bool {
	for _, pattern := range patterns {
		_, file := filepath.Split(filename)
		matched, err := filepath.Match(pattern, file)
		if err != nil {
			panic(err)
		}
		if matched {
			return true
		}
	}
	return false
}

// FilterResourcesByType filters a list of resources that match a single resource type
func FilterResourcesByType(resources []Resource, resourceType string) []Resource {
	if resourceType == "*" {
		return resources
	}
	filtered := make([]Resource, 0)
	for _, resource := range resources {
		if resource.Type == resourceType {
			filtered = append(filtered, resource)
		}
	}
	return filtered
}

// JSONStringify converts a JSON object into an indented string suitable for printing
func JSONStringify(data interface{}) (string, error) {
	b, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return "", err
	}
	return string(b), nil
}
