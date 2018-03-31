package main

import (
	"fmt"
	"github.com/stelligent/config-lint/assertion"
)

// Linter provides the interface for all supported linters
type Linter interface {
	Validate(filenames []string, ruleSet assertion.RuleSet, tags []string, ruleIDs []string) ([]string, []assertion.Violation, error)
	Search(filenames []string, ruleSet assertion.RuleSet, searchExpression string)
}

// ResourceLoader provides the interface that a Linter needs to load a collection of Resource objects
type ResourceLoader interface {
	Load(filename string) ([]assertion.Resource, error)
}

func makeLinter(linterType string, log assertion.LoggingFunction) Linter {
	switch linterType {
	case "Kubernetes":
		return KubernetesLinter{Log: log}
	case "Terraform":
		return TerraformLinter{Log: log}
	case "SecurityGroup":
		return SecurityGroupLinter{Log: log}
	case "LintRules":
		return RulesLinter{Log: log}
	case "YAML":
		return YAMLLinter{Log: log}
	default:
		fmt.Printf("Type not supported: %s\n", linterType)
		return nil
	}
}
