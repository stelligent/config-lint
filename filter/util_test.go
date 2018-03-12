package filter

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
