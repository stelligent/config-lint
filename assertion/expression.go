package assertion

func searchAndMatch(expression Expression, resource Resource) (MatchResult, error) {
	v, err := SearchData(expression.Key, resource.Properties)
	if err != nil {
		return matchError(err)
	}
	match, err := isMatch(v, expression)
	Debugf("Key: %s Output: %v Looking for %v %v\n", expression.Key, v, expression.Op, expression.Value)
	Debugf("ResourceID: %s Type: %s %v\n",
		resource.ID,
		resource.Type,
		match)
	return match, err
}

func orExpression(expressions []Expression, resource Resource) (MatchResult, error) {
	for _, childExpression := range expressions {
		match, err := booleanExpression(childExpression, resource)
		if err != nil {
			return matchError(err)
		}
		if match.Match {
			return matches()
		}
	}
	return doesNotMatch("Or expression fails") // TODO needs more information
}

func xorExpression(expressions []Expression, resource Resource) (MatchResult, error) {
	matchCount := 0
	for _, childExpression := range expressions {
		match, err := booleanExpression(childExpression, resource)
		if err != nil {
			return matchError(err)
		}
		if match.Match {
			matchCount++
		}
	}
	if matchCount == 1 {
		return matches()
	}
	return doesNotMatch("Xor expression fails") // TODO needs more information
}

func andExpression(expressions []Expression, resource Resource) (MatchResult, error) {
	for _, childExpression := range expressions {
		match, err := booleanExpression(childExpression, resource)
		if err != nil {
			return matchError(err)
		}
		if !match.Match {
			return doesNotMatch("And expression fails: %s", match.Message)
		}
	}
	return matches()
}

func notExpression(expressions []Expression, resource Resource) (MatchResult, error) {
	// more than one child filter treated as not any
	for _, childExpression := range expressions {
		match, err := booleanExpression(childExpression, resource)
		if err != nil {
			return matchError(err)
		}
		if match.Match {
			return doesNotMatch("Not expression fails") // TODO needs more information
		}
	}
	return matches()
}

func collectResources(key string, resource Resource) ([]Resource, error) {
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

func everyExpression(collectionExpression CollectionExpression, resource Resource) (MatchResult, error) {
	resources, err := collectResources(collectionExpression.Key, resource)
	if err != nil {
		return matchError(err)
	}
	for _, collectionResource := range resources {
		match, err := andExpression(collectionExpression.Expressions, collectionResource)
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

func someExpression(collectionExpression CollectionExpression, resource Resource) (MatchResult, error) {
	resources, err := collectResources(collectionExpression.Key, resource)
	if err != nil {
		return matchError(err)
	}
	for _, collectionResource := range resources {
		match, err := andExpression(collectionExpression.Expressions, collectionResource)
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

func noneExpression(collectionExpression CollectionExpression, resource Resource) (MatchResult, error) {
	resources, err := collectResources(collectionExpression.Key, resource)
	if err != nil {
		return matchError(err)
	}
	for _, collectionResource := range resources {
		match, err := andExpression(collectionExpression.Expressions, collectionResource)
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

func exactlyOneExpression(collectionExpression CollectionExpression, resource Resource) (MatchResult, error) {
	resources, err := collectResources(collectionExpression.Key, resource)
	if err != nil {
		return matchError(err)
	}
	matchCount := 0
	for _, collectionResource := range resources {
		match, err := andExpression(collectionExpression.Expressions, collectionResource)
		if err != nil {
			return matchError(err)
		}
		if match.Match {
			matchCount++
		}
	}
	if matchCount == 1 {
		return matches()
	}
	return doesNotMatch("ExactlyOne expression fails")
}

func booleanExpression(expression Expression, resource Resource) (MatchResult, error) {
	if expression.Or != nil && len(expression.Or) > 0 {
		return orExpression(expression.Or, resource)
	}
	if expression.Xor != nil && len(expression.Xor) > 0 {
		return xorExpression(expression.Xor, resource)
	}
	if expression.And != nil && len(expression.And) > 0 {
		return andExpression(expression.And, resource)
	}
	if expression.Not != nil && len(expression.Not) > 0 {
		return notExpression(expression.Not, resource)
	}
	if expression.Every.Key != "" {
		return everyExpression(expression.Every, resource)
	}
	if expression.Some.Key != "" {
		return someExpression(expression.Some, resource)
	}
	if expression.None.Key != "" {
		return noneExpression(expression.None, resource)
	}
	if expression.ExactlyOne.Key != "" {
		return exactlyOneExpression(expression.ExactlyOne, resource)
	}
	return searchAndMatch(expression, resource)
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

// CheckExpression validates a single Resource using a single Expression
func CheckExpression(rule Rule, expression Expression, resource Resource) (Result, error) {
	result := Result{
		Status:  "OK",
		Message: "",
	}
	match, err := booleanExpression(expression, resource)
	if err != nil {
		result.Status = "FAILURE"
		result.Message = err.Error()
		return result, err
	}
	if !match.Match {
		if rule.Severity == "" {
			result.Status = "FAILURE"
		} else {
			result.Status = rule.Severity
		}
		result.Message = match.Message
	}
	return result, nil
}
