package assertion

import (
	"testing"
)

func TestUnquotedWithoutQuotes(t *testing.T) {
	if unquoted("Foo") != "Foo" {
		t.Errorf("Unquoted for not quoted string fails")
	}
}

func TestUnquotedWithQuotes(t *testing.T) {
	if unquoted("\"Foo\"") != "Foo" {
		t.Errorf("Unquoted for quoted string fails")
	}
}

func TestIsAbsentEmptyString(t *testing.T) {
	if isAbsent("") != true {
		t.Errorf("isAbsent for empty string fails")
	}
}

func TestIsAbsentEmptyArray(t *testing.T) {
	if isAbsent("[]") != true {
		t.Errorf("isAbsent for empty array fails")
	}
}

func TestIsAbsentNull(t *testing.T) {
	if isAbsent("null") != true {
		t.Errorf("isAbsent for null fails")
	}
}

func TestIsAbsentFalse(t *testing.T) {
	if isAbsent("something") != false {
		t.Errorf("isAbsent for value fails")
	}
}

func TestIntersectTrue(t *testing.T) {
	a := []string{"foo", "bar"}
	b := []string{"bar", "baz"}
	if listsIntersect(a, b) != true {
		t.Errorf("listsIntersect should return true fails")
	}
}

func TestIntersectFalse(t *testing.T) {
	a := []string{"foo", "bar"}
	b := []string{"baz"}
	if listsIntersect(a, b) != false {
		t.Errorf("listsIntersect should return false fails")
	}
}

func TestJSONListsIntersectTrue(t *testing.T) {
	s1 := "[ \"foo\", \"bar\" ]"
	s2 := "[ \"baz\", \"bar\" ]"
	if jsonListsIntersect(s1, s2) != true {
		t.Errorf("JSONIntersect should return true")
	}
}

func TestShouldIncludeFile(t *testing.T) {
	patterns := []string{"*.tf", "*.yml"}
	include, err := ShouldIncludeFile(patterns, "instance.tf")
	if err != nil {
		t.Errorf("ShouldIncludeFile generated an unexpected error: %v", err)
	}
	if !include {
		t.Errorf("ShouldIncludeFile failed to include file with matching pattern")
	}
}

func TestShouldNotIncludeFile(t *testing.T) {
	patterns := []string{"*.tf", "*.yml"}
	include, err := ShouldIncludeFile(patterns, "instance.config")
	if err != nil {
		t.Errorf("ShouldIncludeFile generated an unexpected error: %v", err)
	}
	if include {
		t.Errorf("ShouldIncludeFile failed to exclude file with no matching pattern")
	}
}

func TestFilterShouldIncludeResources(t *testing.T) {
	resources := []Resource{
		Resource{Type: "instance"},
		Resource{Type: "volume"},
	}
	filtered := FilterResourcesByType(resources, "instance")
	if len(filtered) != 1 {
		t.Errorf("FilterResourcesByType expected to match one resource")
	}
}

func TestFilterShouldExcludeResources(t *testing.T) {
	resources := []Resource{
		Resource{Type: "instance"},
		Resource{Type: "volume"},
	}
	filtered := FilterResourcesByType(resources, "database")
	if len(filtered) != 0 {
		t.Errorf("FilterResourcesByType expected to match no resources")
	}
}

func TestFilterShouldIncludeAllResources(t *testing.T) {
	resources := []Resource{
		Resource{Type: "instance"},
		Resource{Type: "volume"},
	}
	filtered := FilterResourcesByType(resources, "*")
	if len(filtered) != len(resources) {
		t.Errorf("FilterResourcesByType expected to include all resources")
	}
}
