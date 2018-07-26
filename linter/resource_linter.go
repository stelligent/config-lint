package linter

import (
	"github.com/stelligent/config-lint/assertion"
)

// ResourceLinter provides the basic validation logic used by all linters
type ResourceLinter struct {
	ValueSource assertion.ValueSource
}

// ValidateResources evaluates a list of Rule objects to a list of Resource objects
func (r ResourceLinter) ValidateResources(resources []assertion.Resource, rules []assertion.Rule) (assertion.ValidationReport, error) {

	report := assertion.ValidationReport{
		ResourcesScanned: []assertion.ScannedResource{},
		Violations:       []assertion.Violation{},
	}

	resolvedRules, violations := assertion.ResolveRules(rules, r.ValueSource)
	report.Violations = append(report.Violations, violations...)
	externalRules := assertion.StandardExternalRuleInvoker{}

	for _, rule := range resolvedRules {
		assertion.Debugf("Rule: ID: %v Message: %s\n", rule.ID, rule.Message)
		filteredResources := assertion.FilterResourcesForRule(resources, rule)
		if len(filteredResources) == 0 {
			assertion.Debugf("No resources to check\n")
		}
		for _, resource := range filteredResources {
			if assertion.ExcludeResource(rule, resource) {
				assertion.Debugf("Ignoring resource %s\n", resource.ID)
			} else {
				assertion.Debugf("Checking resource %s\n", resource.ID)
				status, violations, err := assertion.CheckRule(rule, resource, externalRules)
				if err != nil {
					return report, nil
				}
				report.ResourcesScanned = append(report.ResourcesScanned, assertion.ScannedResource{
					ResourceID:   resource.ID,
					ResourceType: resource.Type,
					RuleID:       rule.ID,
					Status:       status,
					Filename:     resource.Filename,
					LineNumber:   resource.LineNumber,
				})
				report.Violations = append(report.Violations, violations...)
			}
		}
	}
	return report, nil
}
