package linter

import (
	"fmt"
	"github.com/stelligent/config-lint/assertion"
	"io"
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
func NewLinter(ruleSet assertion.RuleSet, filenames []string) (Linter, error) {
	assertion.Debugf("Filenames to scan: %v\n", filenames)
	vs := assertion.StandardValueSource{}
	switch ruleSet.Type {
	case "Kubernetes":
		return FileLinter{Filenames: filenames, ValueSource: vs, Loader: KubernetesResourceLoader{}}, nil
	case "Terraform":
		return FileLinter{Filenames: filenames, ValueSource: vs, Loader: TerraformResourceLoader{}}, nil
	case "SecurityGroup":
		return AWSResourceLinter{Loader: SecurityGroupLoader{}, ValueSource: vs}, nil
	case "IAMUser":
		return AWSResourceLinter{Loader: IAMUserLoader{}, ValueSource: vs}, nil
	case "LintRules":
		return FileLinter{Filenames: filenames, ValueSource: vs, Loader: RulesResourceLoader{}}, nil
	case "YAML":
		return FileLinter{Filenames: filenames, ValueSource: vs, Loader: YAMLResourceLoader{Resources: ruleSet.Resources}}, nil
	default:
		return nil, fmt.Errorf("Type not supported: %s", ruleSet.Type)
	}
}
