package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/ghodss/yaml"
	"github.com/hashicorp/hcl"
	"github.com/jmespath/go-jmespath"
	"io/ioutil"
	"regexp"
	"strings"
)

type LoggingFunction func(string)

func makeLogger(verbose bool) LoggingFunction {
	if verbose {
		return func(message string) {
			fmt.Println(message)
		}
	}
	return func(message string) {}
}

func loadYAML(template string, log LoggingFunction) map[string]interface{} {
	jsonData, err := yaml.YAMLToJSON([]byte(template))
	if err != nil {
		panic(err)
	}
	log(string(jsonData))

	var data interface{}
	err = yaml.Unmarshal(jsonData, &data)
	if err != nil {
		panic(err)
	}
	m := data.(map[string]interface{})
	return m["Resources"].(map[string]interface{})
}

func loadHCL(template string, log LoggingFunction) []interface{} {
	var v interface{}
	err := hcl.Unmarshal([]byte(template), &v)
	if err != nil {
		panic(err)
	}
	jsonData, err := json.MarshalIndent(v, "", "  ")
	log(string(jsonData))

	var hclData interface{}
	err = yaml.Unmarshal(jsonData, &hclData)
	if err != nil {
		panic(err)
	}
	m := hclData.(map[string]interface{})
	results := make([]interface{}, 0)
	if m["resource"] != nil {
		log("Adding resources")
		results = append(results, m["resource"].([]interface{})...)
	}
	if m["data"] != nil {
		log("Adding data")
		results = append(results, m["data"].([]interface{})...)
	}
	return results
}

func searchData(expression string, data interface{}) string {
	if len(expression) == 0 {
		return "null"
	}
	result, err := jmespath.Search(expression, data)
	if err != nil {
		panic(err)
	}
	toJSON, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		panic(err)
	}
	return string(toJSON)
}

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

func MustParseRules(rules string) Rules {
	r := Rules{}
	err := yaml.Unmarshal([]byte(rules), &r)
	if err != nil {
		panic(err)
	}
	return r
}

func unquoted(s string) string {
	if s[0] == '"' {
		return s[1 : len(s)-1]
	}
	return s
}

func isAbsent(s string) bool {
	if s == "null" || s == "[]" {
		return true
	}
	return false
}

func isMatch(searchResult string, op string, value string) bool {
	// TODO see Cloud Custodian for ideas
	// ADD gt, ge, lt, le, not-null, empty, intersect, glob
	switch op {
	case "eq":
		if searchResult == value {
			return true
		}
	case "ne":
		if searchResult != value {
			return true
		}
	case "in":
		for _, v := range strings.Split(value, ",") {
			if v == searchResult {
				return true
			}
		}
	case "notin":
		for _, v := range strings.Split(value, ",") {
			if v == searchResult {
				return false
			}
		}
		return true
	case "absent":
		if isAbsent(searchResult) {
			return true
		}
	case "present":
		if searchResult != "null" {
			return true
		}
	case "contains":
		if strings.Contains(searchResult, value) {
			return true
		}
		return false
	case "regex":
		if regexp.MustCompile(value).MatchString(unquoted(searchResult)) {
			return true
		}
		return false
	}
	return false
}

type TerraformResource struct {
	Id         string
	Type       string
	Properties interface{}
	Filename   string
}

func loadTerraformResources(filename string, hclResources []interface{}) []TerraformResource {
	resources := make([]TerraformResource, 0)
	for _, resource := range hclResources {
		for resourceType, templateResources := range resource.(map[string]interface{}) {
			if templateResources != nil {
				for _, templateResource := range templateResources.([]interface{}) {
					for resourceId, resource := range templateResource.(map[string]interface{}) {
						tr := TerraformResource{
							Id:         resourceId,
							Type:       resourceType,
							Properties: resource.([]interface{})[0],
							Filename:   filename,
						}
						resources = append(resources, tr)
					}
				}
			}
		}
	}
	return resources
}

func loadTerraformRules() string {
	terraformRules, err := ioutil.ReadFile("./rules/terraform.yml")
	if err != nil {
		panic(err)
	}
	return string(terraformRules)
}

type ValidationResult struct {
	RuleId       string
	ResourceId   string
	ResourceType string
	Status       string
	Message      string
	Filename     string
}

func listsIntersect(list1 []string, list2 []string) bool {
	for _, a := range list1 {
		for _, b := range list2 {
			if a == b {
				return true
			}
		}
	}
	return false
}

func filterRulesByTag(rules []Rule, tags []string) []Rule {
	filteredRules := make([]Rule, 0)
	for _, rule := range rules {
		if tags == nil || listsIntersect(tags, rule.Tags) {
			filteredRules = append(filteredRules, rule)
		}
	}
	return filteredRules
}

func searchAndMatch(filter Filter, resource TerraformResource, log LoggingFunction) bool {
	o := unquoted(searchData(filter.Key, resource.Properties))
	status := isMatch(o, filter.Op, filter.Value)
	log(fmt.Sprintf("Key: %s Output: %s Looking for %s %s", filter.Key, o, filter.Op, filter.Value))
	log(fmt.Sprintf("ResourceId: %s Type: %s %s",
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
	if searchAndMatch(filter, resource, log) {
		status = rule.Severity
	}
	return status
}

func filterResourcesByType(resources []TerraformResource, resourceType string) []TerraformResource {
	filtered := make([]TerraformResource, 0)
	for _, resource := range resources {
		if resource.Type == resourceType {
			filtered = append(filtered, resource)
		}
	}
	return filtered
}

func validateTerraformResources(resources []TerraformResource, ruleData Rules, tags []string, log LoggingFunction) []ValidationResult {
	results := make([]ValidationResult, 0)
	for _, rule := range filterRulesByTag(ruleData.Rules, tags) {
		log(fmt.Sprintf("Rule %s: %s", rule.Id, rule.Message))
		for _, filter := range rule.Filters {
			for _, resource := range filterResourcesByType(resources, rule.Resource) {
				log(fmt.Sprintf("Checking resource %s", resource.Id))
				status := searchAndTest(rule, filter, resource, log)
				if status != "OK" {
					results = append(results, ValidationResult{
						RuleId:       rule.Id,
						ResourceId:   resource.Id,
						ResourceType: resource.Type,
						Status:       status,
						Message:      rule.Message,
						Filename:     resource.Filename,
					})
				}
			}
		}
	}
	return results
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

func filterRules(allRules Rules, ruleIds []string) Rules {
	if len(ruleIds) == 0 {
		return allRules
	}
	filteredRules := make([]Rule, 0)
	for _, rule := range allRules.Rules {
		for _, id := range ruleIds {
			if id == rule.Id {
				filteredRules = append(filteredRules, rule)
			}
		}
	}
	return Rules{Rules: filteredRules}
}

func terraform(filename string, tags []string, ruleIds []string, log LoggingFunction) {
	hclTemplate, err := ioutil.ReadFile(filename)
	if err != nil {
		panic(err)
	}
	resources := loadTerraformResources(filename, loadHCL(string(hclTemplate), log))
	rules := filterRules(MustParseRules(loadTerraformRules()), ruleIds)
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
	parseTerraform := flag.Bool("terraform", true, "Validate Terraform template")
	verboseLogging := flag.Bool("verbose", false, "Verbose logging")
	tags := flag.String("tags", "", "Run only tests with tags in this comma separated list")
	rules := flag.String("rules", "", "Run only the rules in this comma separated list")
	flag.Parse()

	for _, filename := range flag.Args() {
		if *parseTerraform {
			terraform(filename, makeTagList(*tags), makeRulesList(*rules), makeLogger(*verboseLogging))
		}
	}
}
