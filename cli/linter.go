package main

import (
	"fmt"
	"github.com/stelligent/config-lint/assertion"
)

// Linter provides the interface for all supported linters
type Linter interface {
	Validate(ruleSet assertion.RuleSet, tags []string, ruleIDs []string) (assertion.ValidationReport, error)
	Search(ruleSet assertion.RuleSet, searchExpression string)
}

func makeLinter(linterType string, args []string, log assertion.LoggingFunction) Linter {
	vs := assertion.StandardValueSource{}
	switch linterType {
	case "Kubernetes":
		return KubernetesLinter{Filenames: args, Log: log, ValueSource: vs}
	case "Terraform":
		return TerraformLinter{Filenames: args, Log: log, ValueSource: vs}
	case "SecurityGroup":
		return AWSResourceLinter{Loader: SecurityGroupLoader{}, Log: log, ValueSource: vs}
	case "IAMUser":
		return AWSResourceLinter{Loader: IAMUserLoader{}, Log: log, ValueSource: vs}
	case "LintRules":
		return RulesLinter{Filenames: args, Log: log, ValueSource: vs}
	case "YAML":
		return YAMLLinter{Filenames: args, Log: log, ValueSource: vs}
	default:
		fmt.Printf("Type not supported: %s\n", linterType)
		return nil
	}
}
