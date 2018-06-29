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
