package main

import (
	"testing"
)

func TestCommandLineOnlyOptions(t *testing.T) {
	tags := "1,2,3"
	emptyString := ""
	verbose := false
	o := CommandLineOptions{
		Tags:             &tags,
		Ids:              &emptyString,
		IgnoreIds:        &emptyString,
		QueryExpression:  &emptyString,
		SearchExpression: &emptyString,
		VerboseReport:    &verbose,
	}
	p := ProfileOptions{}
	l, err := getLinterOptions(o, p)

	if err != nil {
		t.Errorf("getLinterOptions should not return error: %s\n", err.Error())
	}
	if len(l.Tags) != 3 {
		t.Errorf("getLinterOptions should find 3 tags: %v\n", l.Tags)
	}
}

func TestProfileOnlyOptions(t *testing.T) {
	emptyString := ""
	verbose := false
	o := CommandLineOptions{
		Tags:             &emptyString,
		Ids:              &emptyString,
		IgnoreIds:        &emptyString,
		QueryExpression:  &emptyString,
		SearchExpression: &emptyString,
		VerboseReport:    &verbose,
	}
	p := ProfileOptions{
		Tags: []string{"1", "2", "3"},
	}
	l, err := getLinterOptions(o, p)

	if err != nil {
		t.Errorf("getLinterOptions should not return error: %s\n", err.Error())
	}
	if len(l.Tags) != 3 {
		t.Errorf("getLinterOptions should find 3 tags: %v\n", l.Tags)
	}
}

func TestCommandLineOverridesProfile(t *testing.T) {
	tags := "1,2,3,4"
	emptyString := ""
	verbose := false
	o := CommandLineOptions{
		Tags:             &tags,
		Ids:              &emptyString,
		IgnoreIds:        &emptyString,
		QueryExpression:  &emptyString,
		SearchExpression: &emptyString,
		VerboseReport:    &verbose,
	}
	p := ProfileOptions{
		Tags: []string{"1", "2", "3"},
	}
	l, err := getLinterOptions(o, p)

	if err != nil {
		t.Errorf("getLinterOptions should not return error: %s\n", err.Error())
	}
	if len(l.Tags) != 4 {
		t.Errorf("getLinterOptions should find 4 tags: %v\n", l.Tags)
	}
}

func TestCommandLineVariables(t *testing.T) {
	emptyString := ""
	verbose := false
	variables := []string{"namespace=web"}
	o := CommandLineOptions{
		Tags:             &emptyString,
		Ids:              &emptyString,
		IgnoreIds:        &emptyString,
		QueryExpression:  &emptyString,
		SearchExpression: &emptyString,
		VerboseReport:    &verbose,
		Variables:        variables,
	}
	p := ProfileOptions{}
	l, err := getLinterOptions(o, p)

	if err != nil {
		t.Errorf("getLinterOptions should not return error: %s\n", err.Error())
	}
	v, ok := l.Variables["namespace"]
	if !ok {
		t.Errorf("Expecting namespace variable to have a value\n")
	} else {
		if v != "web" {
			t.Errorf("Expecting namespace variable to be 'web', not '%s'\n", v)
		}
	}
}

func TestMergeVariables(t *testing.T) {
	emptyString := ""
	verbose := false
	variables := []string{"namespace=web"}
	o := CommandLineOptions{
		Tags:             &emptyString,
		Ids:              &emptyString,
		IgnoreIds:        &emptyString,
		QueryExpression:  &emptyString,
		SearchExpression: &emptyString,
		VerboseReport:    &verbose,
		Variables:        variables,
	}
	p := ProfileOptions{
		Variables: map[string]string{"kind": "Pod"},
	}
	l, err := getLinterOptions(o, p)

	if err != nil {
		t.Errorf("getLinterOptions should not return error: %s\n", err.Error())
	}
	namespace, ok := l.Variables["namespace"]
	if !ok {
		t.Errorf("Expecting namespace variable to have a value\n")
	} else {
		if namespace != "web" {
			t.Errorf("Expecting namespace variable to be 'web', not '%s'\n", namespace)
		}
	}
	kind, ok := l.Variables["kind"]
	if !ok {
		t.Errorf("Expecting kind variable to have a value\n")
	} else {
		if kind != "Pod" {
			t.Errorf("Expecting kind variable to be 'Pod', not '%s'\n", kind)
		}
	}
}
