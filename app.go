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
	resources := m["resource"].([]interface{})
	data := m["data"].([]interface{})
	return append(resources, data...)
}

func search(expression string, data interface{}) interface{} {
	result, err := jmespath.Search(expression, data)
	if err != nil {
		panic(err)
	}
	return result
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
	// ADD gt, ge, lt, le
	// absent, present, not-null, empty
	// and, or, not, intersect
	// glob
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
	case "regex":
		if regexp.MustCompile(value).MatchString(unquoted(searchResult)) {
			return "OK"
		}
	}
	return severity
}

func cloudFormation(filename string, log LoggingFunction) {
	yamlTemplate, err := ioutil.ReadFile(filename)
	if err != nil {
		panic(err)
	}
	resources := loadYAML(string(yamlTemplate), log)
	cloudFormationRules, err := ioutil.ReadFile("./rules/cloudformation.yml")
	if err != nil {
		panic(err)
	}
	ruleData := MustParseRules(string(cloudFormationRules))
	for _, rule := range ruleData.Rules {
		fmt.Printf("Rule %s: %s\n", rule.Id, rule.Message)
		for _, filter := range rule.Filters {
			for resourceId, resource := range resources {
				o := searchData(filter.Key, resource)
				log(fmt.Sprintf("Key: %s Output: %s Looking for %s %s\n", filter.Key, o, filter.Op, filter.Value))
				fmt.Printf("ResourceId: %s %s\n",
					resourceId,
					isValid(o, filter.Op, filter.Value, rule.Severity))
			}
		}
	}
}

func terraformResourceTypes() []string {
	return []string{
		"aws_instance",
		"aws_iam_role",
		"aws_s3_bucket",
	}
}

func terraform(filename string, log LoggingFunction) {
	hclTemplate, err := ioutil.ReadFile(filename)
	if err != nil {
		panic(err)
	}
	hclResources := loadHCL(string(hclTemplate), log)

	resources := make(map[string]interface{})
	resourceTypes := make(map[string]interface{})
	for _, resource := range hclResources {
		for _, resourceType := range terraformResourceTypes() {
			templateResources := resource.(map[string]interface{})[resourceType]
			if templateResources != nil {
				for _, templateResource := range templateResources.([]interface{}) {
					for resourceId, resource := range templateResource.(map[string]interface{}) {
						resources[resourceId] = resource.([]interface{})[0]
						resourceTypes[resourceId] = resourceType
					}
				}
			}
		}
	}
	terraformRules, err := ioutil.ReadFile("./rules/terraform.yml")
	if err != nil {
		panic(err)
	}
	ruleData := MustParseRules(string(terraformRules))
	for _, rule := range ruleData.Rules {
		fmt.Printf("Rule %s: %s\n", rule.Id, rule.Message)
		for _, filter := range rule.Filters {
			for resourceId, resource := range resources {
				if rule.Resource == resourceTypes[resourceId] {
					o := searchData(filter.Key, resource)
					log(fmt.Sprintf("Key: %s Output: %s Looking for %s %s\n", filter.Key, o, filter.Op, filter.Value))
					fmt.Printf("ResourceId: %s %s\n",
						resourceId,
						isValid(o, filter.Op, filter.Value, rule.Severity))
				} else {
					log(fmt.Sprintf("Skipping rule %s for %s %s\n", rule.Id, resourceId, resourceTypes[resourceId]))
				}
			}
		}
	}
}

func main() {
	parseCloudFormation := flag.Bool("cloudformation", false, "Validate CloudFormation template")
	parseTerraform := flag.Bool("terraform", false, "Validate Terraform template")
	verboseLogging := flag.Bool("verbose", false, "Verbose logging")
	flag.Parse()

	for _, filename := range flag.Args() {
		if *parseCloudFormation {
			cloudFormation(filename, makeLogger(*verboseLogging))
		}
		if *parseTerraform {
			terraform(filename, makeLogger(*verboseLogging))
		}
	}
}
