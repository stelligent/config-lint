package linter

import (
	"bytes"
	"github.com/stelligent/config-lint/assertion"
	"strings"
	"testing"
)

func TestAWSResourceLinterValidate(t *testing.T) {
	options := Options{
		Tags:    []string{},
		RuleIDs: []string{},
	}
	ruleSet := loadRulesForTest("./testdata/rules/aws_sg_resource.yml", t)
	mockLoader := AWSMockLoader{}
	linter := AWSResourceLinter{Loader: mockLoader, ValueSource: TestingValueSource{}}
	report, err := linter.Validate(ruleSet, options)
	if err != nil {
		t.Error("Expecting TestYAMLLinter to not return an error")
	}
	if len(report.ResourcesScanned) != 1 {
		t.Errorf("AWSResourceLinter scanned %d resources, expecting 1", len(report.ResourcesScanned))
	}
	if len(report.Violations) != 0 {
		t.Errorf("AWSResourceLinter returned %d violations, expecting 0", len(report.Violations))
	}
}

func TestAWSResourceLinterSearch(t *testing.T) {
	ruleSet := loadRulesForTest("./testdata/rules/aws_sg_resource.yml", t)
	mockLoader := AWSMockLoader{}
	linter := AWSResourceLinter{Loader: mockLoader, ValueSource: TestingValueSource{}}
	var b bytes.Buffer
	linter.Search(ruleSet, "@", &b)
	if !strings.Contains(b.String(), "Name") {
		t.Errorf("Expecting Search to find Name in %s", b.String())
	}
}

type AWSMockLoader struct{}

func (l AWSMockLoader) Load() ([]assertion.Resource, error) {
	r := assertion.Resource{
		ID:   "1",
		Type: "AWS::S3::Bucket",
		Properties: map[string]interface{}{
			"Name": "Testing",
		},
	}
	return []assertion.Resource{r}, nil
}
