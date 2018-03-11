package main

import (
	"fmt"
)

func searchAndMatch(filter Filter, resource TerraformResource, log LoggingFunction) bool {
	v, err := searchData(filter.Key, resource.Properties)
	if err != nil {
		panic(err)
	}
	status := isMatch(unquoted(v), filter.Op, filter.Value)
	log(fmt.Sprintf("Key: %s Output: %s Looking for %s %s", filter.Key, v, filter.Op, filter.Value))
	log(fmt.Sprintf("ResourceId: %s Type: %s %t",
		resource.Id,
		resource.Type,
		status))
	return status
}

func orOperation(rule Rule, filters []Filter, resource TerraformResource, log LoggingFunction) string {
	for _, childFilter := range filters {
		if searchAndMatch(childFilter, resource, log) {
			return "OK"
		}
	}
	return rule.Severity
}

func andOperation(rule Rule, filters []Filter, resource TerraformResource, log LoggingFunction) string {
	for _, childFilter := range filters {
		if !searchAndMatch(childFilter, resource, log) {
			return rule.Severity
		}
	}
	return "OK"
}

func notOperation(rule Rule, filters []Filter, resource TerraformResource, log LoggingFunction) string {
	for _, childFilter := range filters {
		if searchAndMatch(childFilter, resource, log) {
			return rule.Severity
		}
	}
	return "OK"
}

func applyFilter(rule Rule, filter Filter, resource TerraformResource, log LoggingFunction) string {
	status := "OK"
	if filter.Or != nil && len(filter.Or) > 0 {
		return orOperation(rule, filter.Or, resource, log)
	}
	if filter.And != nil && len(filter.And) > 0 {
		return andOperation(rule, filter.And, resource, log)
	}
	if filter.Not != nil && len(filter.Not) > 0 {
		return notOperation(rule, filter.Not, resource, log)
	}
	if !searchAndMatch(filter, resource, log) {
		status = rule.Severity
	}
	return status
}
