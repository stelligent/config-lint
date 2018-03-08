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

func quoted(s string) string {
	return "\"" + s + "\""
}

func unquoted(s string) string {
	if s[0] == '"' {
		return s[1 : len(s)-1]
	}
	return s
}

func isValid(searchResult, op, value, severity string) string {
	// TODO see Cloud Custodian for ideas
	// ADD gt, ge, lt, le, not-null, empty, and, or, not, intersect, glob
	switch op {
	case "eq":
		if searchResult == quoted(value) {
			return "OK"
		}
	case "ne":
		if searchResult != quoted(value) {
			return "OK"
		}
	case "in":
		for _, v := range strings.Split(value, ",") {
			if quoted(v) == searchResult {
				return "OK"
			}
		}
	case "notin":
		for _, v := range strings.Split(value, ",") {
			if quoted(v) == searchResult {
				return severity
			}
		}
		return "OK"
	case "absent":
		if searchResult == "null" {
			return "OK"
		}
	case "present":
		if searchResult != "null" {
			return "OK"
		}
	case "contains":
		if strings.Contains(searchResult, value) {
			return "OK"
		}
	case "regex":
		if regexp.MustCompile(value).MatchString(unquoted(searchResult)) {
			return "OK"
		}
	}
	return severity
}

func terraformResourceTypes() []string {
	return []string{
		"aws_instance",
		"aws_iam_role",
		"aws_s3_bucket",
	}
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
		for _, resourceType := range terraformResourceTypes() {
			templateResources := resource.(map[string]interface{})[resourceType]
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
	Rule       string
	ResourceId string
	Status     string
	Message    string
	Filename   string
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

func validateTerraformResources(resources []TerraformResource, ruleData Rules, tags []string, log LoggingFunction) []ValidationResult {
	results := make([]ValidationResult, 0)
	for _, rule := range filterRulesByTag(ruleData.Rules, tags) {
		log(fmt.Sprintf("Rule %s: %s", rule.Id, rule.Message))
		for _, filter := range rule.Filters {
			for _, resource := range resources {
				if rule.Resource == resource.Type {
					o := searchData(filter.Key, resource.Properties)
					status := isValid(o, filter.Op, filter.Value, rule.Severity)
					log(fmt.Sprintf("Key: %s Output: %s Looking for %s %s", filter.Key, o, filter.Op, filter.Value))
					log(fmt.Sprintf("ResourceId: %s Type: %s %s",
						resource.Id,
						resource.Type,
						status))
					if status != "OK" {
						results = append(results, ValidationResult{
							Rule:       rule.Id,
							ResourceId: resource.Id,
							Status:     status,
							Message:    rule.Message,
							Filename:   resource.Filename,
						})
					}
				} else {
					log(fmt.Sprintf("Skipping rule %s for %s %s", rule.Id, resource.Id, resource.Type))
				}
			}
		}
	}
	return results
}

func printResults(results []ValidationResult) {
	for _, result := range results {
		fmt.Printf("%s Resource '%s' in '%s': %s (%s)\n",
			result.Status,
			result.ResourceId,
			result.Filename,
			result.Message,
			result.Rule)
	}
}

func terraform(filename string, tags []string, log LoggingFunction) {
	hclTemplate, err := ioutil.ReadFile(filename)
	if err != nil {
		panic(err)
	}
	resources := loadTerraformResources(filename, loadHCL(string(hclTemplate), log))
	rules := MustParseRules(loadTerraformRules())

	results := validateTerraformResources(resources, rules, tags, log)
	printResults(results)
}

func makeTagList(tags string) []string {
	if tags == "" {
		return nil
	}
	return strings.Split(tags, ",")
}

func main() {
	parseTerraform := flag.Bool("terraform", true, "Validate Terraform template")
	verboseLogging := flag.Bool("verbose", false, "Verbose logging")
	tags := flag.String("tags", "", "Run only tests with tags in this comma separated list")
	flag.Parse()

	for _, filename := range flag.Args() {
		if *parseTerraform {
			terraform(filename, makeTagList(*tags), makeLogger(*verboseLogging))
		}
	}
}
