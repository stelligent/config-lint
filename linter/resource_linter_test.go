package linter

import (
	"testing"
)

func TestIgnoreResource(t *testing.T) {
	options := Options{
		Tags:    []string{},
		RuleIDs: []string{},
	}
	filenames := []string{"./testdata/resources/terraform_instance.tf"}
	linter := FileLinter{Filenames: filenames, ValueSource: TestingValueSource{}, Loader: TerraformResourceLoader{}}
	ruleSet := loadRulesForTest("./testdata/rules/exclude_resource.yml", t)
	report, err := linter.Validate(ruleSet, options)
	if err != nil {
		t.Error("Expecting TestIgnoreResource to not return an error")
	}
	if len(report.ResourcesScanned) != 0 {
		t.Errorf("TestIgnoreResource scanned %d resources, expecting 0", len(report.ResourcesScanned))
	}
}
