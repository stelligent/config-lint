package linter

import (
	"reflect"
	"testing"
)

type NewLinterTestCase struct {
	Filename string
	TypeName string
}

func TestNewLinter(t *testing.T) {

	testCases := []NewLinterTestCase{
		{"./testdata/rules/terraform_instance.yml", "FileLinter"},
		{"./testdata/rules/generic.yml", "FileLinter"},
		{"./testdata/rules/aws_sg_resource.yml", "AWSResourceLinter"},
		{"./testdata/rules/aws_iam_resource.yml", "AWSResourceLinter"},
		{"./testdata/rules/kubernetes.yml", "FileLinter"},
		{"./testdata/rules/rules.yml", "FileLinter"},
	}

	for _, tc := range testCases {
		ruleSet := loadRulesForTest(tc.Filename, t)
		l, err := NewLinter(ruleSet, []string{})
		if err != nil {
			t.Errorf("Expecting TestNewLinter to not return an error: %s", err.Error())
		}
		n := reflect.TypeOf(l).Name()
		if n != tc.TypeName {
			t.Errorf("Expecting NewLinter expected %s, not %s ", tc.TypeName, n)
		}
	}
}

func TestUnknownLinterType(t *testing.T) {
	ruleSet := loadRulesForTest("./testdata/rules/unknown.yml", t)
	_, err := NewLinter(ruleSet, []string{})
	if err == nil {
		t.Errorf("Expecting NewLinter to return an error for unsupported type")
	}
}
