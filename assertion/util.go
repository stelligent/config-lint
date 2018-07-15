package assertion

import (
	"encoding/json"
	"fmt"
	"path/filepath"
	"time"
)

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

func isEmpty(data interface{}) bool {
	switch v := data.(type) {
	case nil:
		return true
	case string:
		return len(v) == 0
	case []interface{}:
		return len(v) == 0
	case []map[string]interface{}:
		return len(v) == 0
	default:
		Debugf("isEmpty default: %v %T\n", data, data)
		return false
	}
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
func ShouldIncludeFile(patterns []string, filename string) (bool, error) {
	for _, pattern := range patterns {
		_, file := filepath.Split(filename)
		matched, err := filepath.Match(pattern, file)
		if err != nil {
			return false, err
		}
		if matched {
			return true, nil
		}
	}
	return false, nil
}

// FilterResourcesByType filters a list of resources that match a single resource type
func FilterResourcesByType(resources []Resource, resourceType string, resourceCategory string) []Resource {
	if resourceType == "*" {
		return resources
	}
	filtered := make([]Resource, 0)
	for _, resource := range resources {
		if resource.Type == resourceType && categoryMatches(resourceCategory, resource.Category) {
			filtered = append(filtered, resource)
		}
	}
	return filtered
}

// FilterResourcesByTypes filters a list of resources that match a slice of resource types
func FilterResourcesByTypes(resources []Resource, resourceTypes []string, resourceCategory string) []Resource {
	filtered := make([]Resource, 0)
	for _, resource := range resources {
		if SliceContains(resourceTypes, resource.Type) && categoryMatches(resourceCategory, resource.Category) {
			filtered = append(filtered, resource)
		}
	}
	return filtered
}

func categoryMatches(c1, c2 string) bool {
	if c1 == "" || c1 == "*" {
		return true
	}
	return c1 == c2
}

// JSONStringify converts a JSON object into an indented string suitable for printing
func JSONStringify(data interface{}) (string, error) {
	b, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return "", err
	}
	return string(b), nil
}

func currentTime() string {
	return time.Now().UTC().Format(time.RFC3339)
}

func SliceContains(list []string, value string) bool {
	for _, item := range list {
		if item == value {
			return true
		}
	}
	return false
}

// FilterResourcesForRule returns resources applicable to the given rule
func FilterResourcesForRule(resources []Resource, rule Rule) []Resource {
	var filteredResources []Resource
	if rule.Resource != "" {
		Debugf("filtering rule resources on Resource string")
		filteredResources = FilterResourcesByType(resources, rule.Resource, rule.Category)
	} else {
		Debugf("filtering rule resources on Resources slice")
		filteredResources = FilterResourcesByTypes(resources, rule.Resources, rule.Category)
	}
	return filteredResources
}
