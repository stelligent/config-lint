package assertion

import (
	"fmt"
)

func searchAndMatch(assertion Assertion, resource Resource, log LoggingFunction) (MatchResult, error) {
	v, err := SearchData(assertion.Key, resource.Properties)
	if err != nil {
		return matchError(err)
	}
	match, err := isMatch(v, assertion.Op, assertion.Value, assertion.ValueType)
	log(fmt.Sprintf("Key: %s Output: %v Looking for %v %v", assertion.Key, v, assertion.Op, assertion.Value))
	log(fmt.Sprintf("ResourceID: %s Type: %s %v",
		resource.ID,
		resource.Type,
		match))
	return match, err
}

func orExpression(assertions []Assertion, resource Resource, log LoggingFunction) (MatchResult, error) {
	for _, childAssertion := range assertions {
		match, err := booleanExpression(childAssertion, resource, log)
		if err != nil {
			return matchError(err)
		}
		if match.Match {
			return matches()
		}
	}
	return doesNotMatch("Or expression fails") // TODO needs more information
}

func andExpression(assertions []Assertion, resource Resource, log LoggingFunction) (MatchResult, error) {
	for _, childAssertion := range assertions {
		match, err := booleanExpression(childAssertion, resource, log)
		if err != nil {
			return matchError(err)
		}
		if !match.Match {
			return doesNotMatch("And expression fails: %s", match.Message)
		}
	}
	return matches()
}

func notExpression(assertions []Assertion, resource Resource, log LoggingFunction) (MatchResult, error) {
	// more than one child filter treated as not any
	for _, childAssertion := range assertions {
		match, err := booleanExpression(childAssertion, resource, log)
		if err != nil {
			return matchError(err)
		}
		if match.Match {
			return doesNotMatch("Not expression failsL %s", match.Message)
		}
	}
	return matches()
}

func collectResources(key string, resource Resource, log LoggingFunction) ([]Resource, error) {
	resources := make([]Resource, 0)
	value, err := SearchData(key, resource.Properties)
	if err != nil {
		return resources, err
	}
	if collection, ok := value.([]interface{}); ok {
		for _, properties := range collection {
			collectionResource := Resource{
				ID:         resource.ID,
				Type:       resource.Type,
				Properties: properties,
				Filename:   resource.Filename,
			}
			resources = append(resources, collectionResource)
		}
	}
	return resources, nil
}

func everyExpression(collectionAssertion CollectionAssertion, resource Resource, log LoggingFunction) (MatchResult, error) {
	resources, err := collectResources(collectionAssertion.Key, resource, log)
	if err != nil {
		return matchError(err)
	}
	for _, collectionResource := range resources {
		match, err := andExpression(collectionAssertion.Assertions, collectionResource, log)
		if err != nil {
			return matchError(err)
		}
		if !match.Match {
			// at least one element is false, so entire expression is false
			return doesNotMatch("Every expression fails: %s", match.Message)
		}
	}
	// every element passes, so entire expression is true
	return matches()
}

func someExpression(collectionAssertion CollectionAssertion, resource Resource, log LoggingFunction) (MatchResult, error) {
	resources, err := collectResources(collectionAssertion.Key, resource, log)
	if err != nil {
		return matchError(err)
	}
	for _, collectionResource := range resources {
		match, err := andExpression(collectionAssertion.Assertions, collectionResource, log)
		if err != nil {
			return matchError(err)
		}
		// at least one element passes, so entire expression is true
		if match.Match {
			return matches()
		}
	}
	// no element passes, so entire expression is false
	return doesNotMatch("Some expression fails") // TODO needs more information
}

func noneExpression(collectionAssertion CollectionAssertion, resource Resource, log LoggingFunction) (MatchResult, error) {
	resources, err := collectResources(collectionAssertion.Key, resource, log)
	if err != nil {
		return matchError(err)
	}
	for _, collectionResource := range resources {
		match, err := andExpression(collectionAssertion.Assertions, collectionResource, log)
		if err != nil {
			return matchError(err)
		}
		// at least one element passes, so entire expression is false
		if match.Match {
			return doesNotMatch("None expression fails: %s", match.Message)
		}
	}
	// no element passes, so entire expression is true
	return matches()
}

func booleanExpression(assertion Assertion, resource Resource, log LoggingFunction) (MatchResult, error) {
	if assertion.Or != nil && len(assertion.Or) > 0 {
		return orExpression(assertion.Or, resource, log)
	}
	if assertion.And != nil && len(assertion.And) > 0 {
		return andExpression(assertion.And, resource, log)
	}
	if assertion.Not != nil && len(assertion.Not) > 0 {
		return notExpression(assertion.Not, resource, log)
	}
	if assertion.Every.Key != "" {
		return everyExpression(assertion.Every, resource, log)
	}
	if assertion.Some.Key != "" {
		return someExpression(assertion.Some, resource, log)
	}
	if assertion.None.Key != "" {
		return noneExpression(assertion.None, resource, log)
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
func CheckAssertion(rule Rule, assertion Assertion, resource Resource, log LoggingFunction) (string, error) {
	status := "OK"
	match, err := booleanExpression(assertion, resource, log)
	if err != nil {
		return "FAILURE", err
	}
	if !match.Match {
		status = rule.Severity
	}
	return status, nil
}
