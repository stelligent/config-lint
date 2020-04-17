package main

import (
	"testing"
)

func emptyCommandLineOptions() CommandLineOptions {
	emptyString := ""
	verbose := false
	return CommandLineOptions{
		Tags:             &emptyString,
		Ids:              &emptyString,
		IgnoreIds:        &emptyString,
		QueryExpression:  &emptyString,
		SearchExpression: &emptyString,
		VerboseReport:    &verbose,
		TerraformParser:  &emptyString,
	}
}

func TestCommandLineOnlyOptions(t *testing.T) {
	tags := "1,2,3"
	o := emptyCommandLineOptions()
	o.Tags = &tags
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
	o := emptyCommandLineOptions()
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
	o := emptyCommandLineOptions()
	o.Tags = &tags
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
	o := emptyCommandLineOptions()
	o.Variables = []string{"namespace=web"}
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
	o := emptyCommandLineOptions()
	o.Variables = []string{"namespace=web"}
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

func TestLoadProfile(t *testing.T) {
	p, err := loadProfile("./testdata/profile.yml")
	if err != nil {
		t.Errorf("Expecting loadProfile to run without error: %v\n", err.Error())
	}
	if len(p.Tags) != 1 || p.Tags[0] != "iam" {
		t.Errorf("Expecting single tag in profile: %v\n", p.Tags)
	}
}

func TestProfileExclude(t *testing.T) {
	p, err := loadProfile("./testdata/profile.yml")
	if err != nil {
		t.Errorf("Expecting loadProfile to run without error: %v\n", err.Error())
	}
	o := emptyCommandLineOptions()
	l, err := getLinterOptions(o, p)
	if err != nil {
		t.Errorf("Expecting getLinterOptions to run without error: %v\n", err.Error())
	}
	if len(l.ExcludePatterns) != 3 {
		t.Errorf("Expecting 3 excludes in total using 'exclude' and 'exclude_from' in profile: %v\n", l.ExcludePatterns)
	}
	if l.ExcludePatterns[0] != "this_file_will_be_excluded.tf" {
		t.Errorf("Expecting 1st pattern using 'exclude' in profile to be 'this_file_will_be_excluded.tf', not '%s'", l.ExcludePatterns[0])
	}
	if l.ExcludePatterns[1] != "*1.tf" {
		t.Errorf("Expecting 2nd pattern using 'exclude_from' in profile to be '*1.tf', not '%s'", l.ExcludePatterns[1])
	}
	if l.ExcludePatterns[2] != "*2.tf" {
		t.Errorf("Expecting 3rd pattern using 'exclude_from' in profile to be '*2.tf', not '%s'", l.ExcludePatterns[2])
	}
}

func TestValidateParser(t *testing.T) {
	parser, err := validateParser("")
	if err != nil {
		t.Errorf("Expected %s, got %v", parser, err)
	}
	parser, err = validateParser("tf11")
	if err != nil {
		t.Errorf("Expected %s, got %v", parser, err)
	}
	parser, err = validateParser("tf12")
	if err != nil {
		t.Errorf("Expected %s, got %v", parser, err)
	}
	parser, err = validateParser("tf13")
	if err == nil {
		t.Errorf("Expected %v, got nil", err)
	}
}
