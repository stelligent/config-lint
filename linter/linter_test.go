package linter

import (
	"github.com/stelligent/config-lint/assertion"
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
		{"./testdata/rules/generic-yaml.yml", "FileLinter"},
		{"./testdata/rules/generic-json.yml", "FileLinter"},
		{"./testdata/rules/generic-csv.yml", "FileLinter"},
		{"./testdata/rules/kubernetes.yml", "FileLinter"},
		{"./testdata/rules/rules.yml", "FileLinter"},
	}

	vs := MockValueSource{}
	for _, tc := range testCases {
		tfParserOptions := []string{"", "tf11", "tf12"}
		ruleSet := loadRulesForTest(tc.Filename, t)
		for _, tfPO := range tfParserOptions {
			l, err := NewLinter(ruleSet, vs, []string{}, tfPO)
			if err != nil {
				t.Errorf("Expecting TestNewLinter to not return an error: %s", err.Error())
			}
			n := reflect.TypeOf(l).Name()
			if n != tc.TypeName {
				t.Errorf("Expecting NewLinter expected %s, not %s ", tc.TypeName, n)
			}
		}
	}
}

func TestUnknownLinterType(t *testing.T) {
	ruleSet := loadRulesForTest("./testdata/rules/unknown.yml", t)
	vs := MockValueSource{}
	tfParserOptions := []string{"", "tf11", "tf12"}
	for _, tfPO := range tfParserOptions {
		_, err := NewLinter(ruleSet, vs, []string{}, tfPO)
		if err == nil {
			t.Errorf("Expecting NewLinter to return an error for unsupported type")
		}
	}
}

type MockValueSource struct{}

func (m MockValueSource) GetValue(e assertion.Expression) (string, error) {
	return "", nil
}
