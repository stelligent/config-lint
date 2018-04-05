package main

import (
	"fmt"
	"github.com/stelligent/config-lint/assertion"
)

// ResourceLinter provides the basic validation logic used by all linters
type ResourceLinter struct {
	Log         assertion.LoggingFunction
	ValueSource assertion.ValueSource
}

// ValidateResources evaluates a list of Rule objects to a list of Resource objects
func (r ResourceLinter) ValidateResources(resources []assertion.Resource, rules []assertion.Rule) (assertion.ValidationReport, error) {

	report := assertion.ValidationReport{
		ResourcesScanned: []assertion.ScannedResource{},
		Violations:       []assertion.Violation{},
	}

	resolvedRules := assertion.ResolveRules(rules, r.ValueSource, r.Log)
	externalRules := assertion.StandardExternalRuleInvoker{Log: r.Log}

	for _, rule := range resolvedRules {
		r.Log(fmt.Sprintf("Rule %s: %s", rule.ID, rule.Message))
		for _, resource := range assertion.FilterResourcesByType(resources, rule.Resource) {
			if assertion.ExcludeResource(rule, resource) {
				r.Log(fmt.Sprintf("Ignoring resource %s", resource.ID))
			} else {
				status, violations, err := assertion.CheckRule(rule, resource, externalRules, r.Log)
				if err != nil {
					return report, nil
				}
				report.ResourcesScanned = append(report.ResourcesScanned, assertion.ScannedResource{
					ResourceID:   resource.ID,
					ResourceType: resource.Type,
					RuleID:       rule.ID,
					Status:       status,
				})
				report.Violations = append(report.Violations, violations...)
			}
		}
	}
	return report, nil
}
