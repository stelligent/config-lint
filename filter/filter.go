package filter

import (
	"fmt"
)

func searchAndMatch(filter Filter, resource Resource, valueSource ValueSource, log LoggingFunction) bool {
	v, err := SearchData(filter.Key, resource.Properties)
	if err != nil {
		panic(err)
	}
	value := valueSource.GetValue(filter)
	match := isMatch(unquoted(v), filter.Op, value)
	log(fmt.Sprintf("Key: %s Output: %s Looking for %s %s", filter.Key, v, filter.Op, value))
	log(fmt.Sprintf("ResourceId: %s Type: %s %t",
		resource.Id,
		resource.Type,
		match))
	return match
}

func orOperation(filters []Filter, resource Resource, valueSource ValueSource, log LoggingFunction) bool {
	for _, childFilter := range filters {
		if booleanOperation(childFilter, resource, valueSource, log) {
			return true
		}
	}
	return false
}

func andOperation(filters []Filter, resource Resource, valueSource ValueSource, log LoggingFunction) bool {
	for _, childFilter := range filters {
		if !booleanOperation(childFilter, resource, valueSource, log) {
			return false
		}
	}
	return true
}

func notOperation(filters []Filter, resource Resource, valueSource ValueSource, log LoggingFunction) bool {
	for _, childFilter := range filters {
		if booleanOperation(childFilter, resource, valueSource, log) {
			return false
		}
	}
	return true
}

func booleanOperation(filter Filter, resource Resource, valueSource ValueSource, log LoggingFunction) bool {
	if filter.Or != nil && len(filter.Or) > 0 {
		return orOperation(filter.Or, resource, valueSource, log)
	}
	if filter.And != nil && len(filter.And) > 0 {
		return andOperation(filter.And, resource, valueSource, log)
	}
	if filter.Not != nil && len(filter.Not) > 0 {
		return notOperation(filter.Not, resource, valueSource, log)
	}
	return searchAndMatch(filter, resource, valueSource, log)
}

func ExcludeResource(rule Rule, resource Resource) bool {
	for _, id := range rule.Except {
		if id == resource.Id {
			return true
		}
	}
	return false
}

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

func ApplyFilter(rule Rule, filter Filter, resource Resource, valueSource ValueSource, log LoggingFunction) string {
	status := "OK"
	if !booleanOperation(filter, resource, valueSource, log) {
		status = rule.Severity
	}
	return status
}
