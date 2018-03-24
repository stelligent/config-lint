package assertion

import (
	"fmt"
)

func searchAndMatch(assertion Assertion, resource Resource, log LoggingFunction) bool {
	v, err := SearchData(assertion.Key, resource.Properties)
	if err != nil {
		panic(err)
	}
	match := isMatch(v, assertion.Op, assertion.Value)
	log(fmt.Sprintf("Key: %s Output: %s Looking for %s %s", assertion.Key, v, assertion.Op, assertion.Value))
	log(fmt.Sprintf("ResourceID: %s Type: %s %t",
		resource.ID,
		resource.Type,
		match))
	return match
}

func orOperation(assertions []Assertion, resource Resource, log LoggingFunction) bool {
	for _, childAssertion := range assertions {
		if booleanOperation(childAssertion, resource, log) {
			return true
		}
	}
	return false
}

func andOperation(assertions []Assertion, resource Resource, log LoggingFunction) bool {
	for _, childAssertion := range assertions {
		if !booleanOperation(childAssertion, resource, log) {
			return false
		}
	}
	return true
}

func notOperation(assertions []Assertion, resource Resource, log LoggingFunction) bool {
	// more than one child filter treated as not any
	for _, childAssertion := range assertions {
		if booleanOperation(childAssertion, resource, log) {
			return false
		}
	}
	return true
}

func booleanOperation(assertion Assertion, resource Resource, log LoggingFunction) bool {
	if assertion.Or != nil && len(assertion.Or) > 0 {
		return orOperation(assertion.Or, resource, log)
	}
	if assertion.And != nil && len(assertion.And) > 0 {
		return andOperation(assertion.And, resource, log)
	}
	if assertion.Not != nil && len(assertion.Not) > 0 {
		return notOperation(assertion.Not, resource, log)
	}
	return searchAndMatch(assertion, resource, log)
}

// ExcludeResource when resource.ID included in list of exceptions
func ExcludeResource(rule Rule, resource Resource) bool {
	for _, id := range rule.Except {
		if id == resource.ID {
			return true
		}
	}
	return false
}

// FilterResourceExceptions filters out resources that should not be validated
func FilterResourceExceptions(rule Rule, resources []Resource) []Resource {
	if rule.Except == nil || len(rule.Except) == 0 {
		return resources
	}
	filtered := make([]Resource, 0)
	for _, resource := range resources {
		if ExcludeResource(rule, resource) {
			filtered = append(filtered, resource)
		}
	}
	return filtered
}

// CheckAssertion validates a single Resource using a single Assertion
func CheckAssertion(rule Rule, assertion Assertion, resource Resource, log LoggingFunction) string {
	status := "OK"
	if !booleanOperation(assertion, resource, log) {
		status = rule.Severity
	}
	return status
}
