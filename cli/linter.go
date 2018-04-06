package main

import (
	"fmt"
	"github.com/stelligent/config-lint/assertion"
)

type (
	// Linter provides the interface for all supported linters
	Linter interface {
		Validate(ruleSet assertion.RuleSet, options LinterOptions) (assertion.ValidationReport, error)
		Search(ruleSet assertion.RuleSet, searchExpression string)
	}

	// LinterOptions configures what resources will be linted
	LinterOptions struct {
		Tags    []string
		RuleIDs []string
	}
)

func makeLinter(linterType string, args []string) (Linter, error) {
	vs := assertion.StandardValueSource{}
	switch linterType {
	case "Kubernetes":
		return KubernetesLinter{Filenames: args, ValueSource: vs}, nil
	case "Terraform":
		return TerraformLinter{Filenames: args, ValueSource: vs}, nil
	case "SecurityGroup":
		return AWSResourceLinter{Loader: SecurityGroupLoader{}, ValueSource: vs}, nil
	case "IAMUser":
		return AWSResourceLinter{Loader: IAMUserLoader{}, ValueSource: vs}, nil
	case "LintRules":
		return RulesLinter{Filenames: args, ValueSource: vs}, nil
	case "YAML":
		return YAMLLinter{Filenames: args, ValueSource: vs}, nil
	default:
		return nil, fmt.Errorf("Type not supported: %s", linterType)
	}
}
