package linter

import (
	"fmt"
	"io"

	"github.com/stelligent/config-lint/assertion"
)

type (
	// Linter provides the interface for all supported linters
	Linter interface {
		Validate(ruleSet assertion.RuleSet, options Options) (assertion.ValidationReport, error)
		Search(ruleSet assertion.RuleSet, searchExpression string, w io.Writer)
	}

	// Options configures what resources will be linted
	Options struct {
		Tags          []string
		RuleIDs       []string
		IgnoreRuleIDs []string
	}
)

// NewLinter create the right kind of Linter based on the type argument
func NewLinter(ruleSet assertion.RuleSet, vs assertion.ValueSource, filenames []string, parser string) (Linter, error) {
	assertion.Debugf("Filenames to scan: %v\n", filenames)
	switch ruleSet.Type {
	case "Kubernetes":
		return FileLinter{Filenames: filenames, ValueSource: vs, Loader: KubernetesResourceLoader{}}, nil
	case "Terraform":
		if parser == "tf12" {
			return FileLinter{Filenames: filenames, ValueSource: vs, Loader: Terraform12ResourceLoader{}}, nil
		} else {
			return FileLinter{Filenames: filenames, ValueSource: vs, Loader: TerraformResourceLoader{}}, nil
		}
	case "Terraform12":
		if parser == "tf11" {
			return FileLinter{Filenames: filenames, ValueSource: vs, Loader: TerraformResourceLoader{}}, nil
		} else {
			return FileLinter{Filenames: filenames, ValueSource: vs, Loader: Terraform12ResourceLoader{}}, nil
		}
	case "LintRules":
		return FileLinter{Filenames: filenames, ValueSource: vs, Loader: RulesResourceLoader{}}, nil
	case "YAML":
		return FileLinter{Filenames: filenames, ValueSource: vs, Loader: YAMLResourceLoader{Resources: ruleSet.Resources}}, nil
	case "JSON":
		return FileLinter{Filenames: filenames, ValueSource: vs, Loader: JSONResourceLoader{Resources: ruleSet.Resources}}, nil
	case "CSV":
		return FileLinter{Filenames: filenames, ValueSource: vs, Loader: CSVResourceLoader{Columns: ruleSet.Columns}}, nil
	default:
		return nil, fmt.Errorf("Type not supported: %s", ruleSet.Type)
	}
}
