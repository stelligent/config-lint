package main

import (
	"fmt"
	"github.com/stelligent/config-lint/assertion"
)

// ResourceLinter provides the basic validation logic used by all linters
type ResourceLinter struct {
	Log assertion.LoggingFunction
}

// ValidateResources evaluates a list of Rule objects to a list of Resource objects
func (r ResourceLinter) ValidateResources(resources []assertion.Resource, rules []assertion.Rule) ([]assertion.ScannedResource, []assertion.Violation, error) {

	scannedResources := make([]assertion.ScannedResource, 0)

	valueSource := assertion.StandardValueSource{Log: r.Log}
	resolvedRules := assertion.ResolveRules(rules, valueSource, r.Log)
	externalRules := assertion.StandardExternalRuleInvoker{Log: r.Log}

	allViolations := make([]assertion.Violation, 0)
	for _, rule := range resolvedRules {
		r.Log(fmt.Sprintf("Rule %s: %s", rule.ID, rule.Message))
		for _, resource := range assertion.FilterResourcesByType(resources, rule.Resource) {
			if assertion.ExcludeResource(rule, resource) {
				r.Log(fmt.Sprintf("Ignoring resource %s", resource.ID))
			} else {
				status, violations, err := assertion.CheckRule(rule, resource, externalRules, r.Log)
				if err != nil {
					return scannedResources, allViolations, err
				}
				scannedResources = append(scannedResources, assertion.ScannedResource{
					ResourceID:   resource.ID,
					ResourceType: resource.Type,
					Status:       status,
				})
				allViolations = append(allViolations, violations...)
			}
		}
	}
	return scannedResources, allViolations, nil
}
