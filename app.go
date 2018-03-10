package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"strings"
)

type Filter struct {
	Type  string
	Key   string
	Op    string
	Value string
	Or    []Filter
	And   []Filter
	Not   []Filter
}

type Rule struct {
	Id       string
	Message  string
	Severity string
	Resource string
	Filters  []Filter
	Tags     []string
}

type Rules struct {
	Rules []Rule
}

type ValidationResult struct {
	RuleId       string
	ResourceId   string
	ResourceType string
	Status       string
	Message      string
	Filename     string
}

func searchAndMatch(filter Filter, resource TerraformResource, log LoggingFunction) bool {
	o := unquoted(searchData(filter.Key, resource.Properties))
	status := isMatch(o, filter.Op, filter.Value)
	log(fmt.Sprintf("Key: %s Output: %s Looking for %s %s", filter.Key, o, filter.Op, filter.Value))
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

func searchAndTest(rule Rule, filter Filter, resource TerraformResource, log LoggingFunction) string {
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

func printResults(results []ValidationResult) {
	for _, result := range results {
		fmt.Printf("%s %s '%s' in '%s': %s (%s)\n",
			result.Status,
			result.ResourceType,
			result.ResourceId,
			result.Filename,
			result.Message,
			result.RuleId)
	}
}

func terraform(filename string, tags []string, ruleIds []string, log LoggingFunction) {
	hclTemplate, err := ioutil.ReadFile(filename)
	if err != nil {
		panic(err)
	}
	resources := loadTerraformResources(filename, loadHCL(string(hclTemplate), log))
	rules := filterRulesById(MustParseRules(loadTerraformRules()), ruleIds)
	results := validateTerraformResources(resources, rules, tags, log)
	printResults(results)
}

func makeTagList(tags string) []string {
	if tags == "" {
		return nil
	}
	return strings.Split(tags, ",")
}

func makeRulesList(ruleIds string) []string {
	if ruleIds == "" {
		return nil
	}
	return strings.Split(ruleIds, ",")
}

func main() {
	verboseLogging := flag.Bool("verbose", false, "Verbose logging")
	tags := flag.String("tags", "", "Run only tests with tags in this comma separated list")
	rules := flag.String("rules", "", "Run only the rules in this comma separated list")
	flag.Parse()

	for _, filename := range flag.Args() {
		terraform(filename, makeTagList(*tags), makeRulesList(*rules), makeLogger(*verboseLogging))
	}
}
