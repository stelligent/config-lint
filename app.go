package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/ghodss/yaml"
	"github.com/hashicorp/hcl"
	"github.com/jmespath/go-jmespath"
	"regexp"
	"strings"
)

var yamlTemplate = `
Resources:
  Instance:
    Type: "AWS::EC2::Instance"
    Properties:
      ImageId: ami-f2d3638a
      InstanceType: t2.micro
  AnotherInstance:
    Type: "AWS::EC2::Instance"
    Properties:
      ImageId: ami-f2d3638a
      InstanceType: m3.medium
  ThirdInstance:
    Type: "AWS::EC2::Instance"
    Properties:
      ImageId: ami-f2d3638b
      InstanceType: c4.large
`

var hclTemplate = `
resource "aws_instance" "first" {
	ami = "ami-f2d3638a"
	instance_type = "t2.micro"
}
resource "aws_instance" "second" {
	ami = "ami-f2d3638a"
	instance_type = "m3.medium"
	tags {
		Department = "Operations"
	}
}
resource "aws_instance" "third" {
	ami = "ami-f2d3638b"
	instance_type = "c4.large"
}
resource "aws_iam_role" "role1" {
    name = "role1"
    assume_role_policy = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
     {
        "Action": "*",
        "Principal": { "Service": "ec2.amazonaws.com" }
        "Effect": "Allow"
        "Resources": "*"
     }
  ]
}
EOF
}
`

var cloudFormationRules = `
Rules:
  - id: R1
    message: Check instance type
    filters:
      - type: value
        key: "Properties.InstanceType"
        op: in
        value: t2.micro,m3.medium
    severity: WARNING
  - id: R2
    message: Check image id
    filters:
      - type: value
        key: "Properties.ImageId"
        op: eq
        value: ami-f2d3638a
    severity: FAILURE
`

var terraformRules = `
Rules:
  - id: R1
    message: Check instance type
    resource: aws_instance
    filters:
      - type: value
        key: instance_type
        op: in
        value: t2.micro,m3.medium
    severity: WARNING
  - id: R2
    message: Check image id
    resource: aws_instance
    filters:
      - type: value
        key: ami
        op: in
        value: ami-f2d3638a
    severity: FAILURE
  - id: R3
    message: Check tags
    resource: aws_instance
    filters:
      - type: value
        key: "tags[].Department | [0]"
        op: regex
        value: Operations
    severity: WARNING
  - id: R4
    message: Check role name
    resource: aws_iam_role
    filters:
      - type: value
        key: name
        op: regex
        value: "role*"
    severity: WARNING
`

func loadYAML(template string) map[string]interface{} {
	jsonData, err := yaml.YAMLToJSON([]byte(template))
	if err != nil {
		panic(err)
	}

	var data interface{}
	err = yaml.Unmarshal(jsonData, &data)
	if err != nil {
		panic(err)
	}
	m := data.(map[string]interface{})
	return m["Resources"].(map[string]interface{})
}

func loadHCL(template string) []interface{} {
	var v interface{}
	err := hcl.Unmarshal([]byte(template), &v)
	if err != nil {
		panic(err)
	}
	jsonData, err := json.MarshalIndent(v, "", "  ")
	fmt.Println(string(jsonData))

	var data interface{}
	err = yaml.Unmarshal(jsonData, &data)
	if err != nil {
		panic(err)
	}
	m := data.(map[string]interface{})
	return m["resource"].([]interface{})
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

func cloudFormation() {
	resources := loadYAML(yamlTemplate)
	ruleData := MustParseRules(cloudFormationRules)
	for _, rule := range ruleData.Rules {
		fmt.Printf("Rule %s: %s\n", rule.Id, rule.Message)
		for _, filter := range rule.Filters {
			for resourceId, resource := range resources {
				o := searchData(filter.Key, resource)
				//fmt.Printf("Key: %s Output: %s Looking for %s %s\n", filter.Key, o, filter.Op, filter.Value)
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
	}
}

func terraform() {
	hclResources := loadHCL(hclTemplate)

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
	ruleData := MustParseRules(terraformRules)
	for _, rule := range ruleData.Rules {
		fmt.Printf("Rule %s: %s\n", rule.Id, rule.Message)
		for _, filter := range rule.Filters {
			for resourceId, resource := range resources {
				if rule.Resource == resourceTypes[resourceId] {
					o := searchData(filter.Key, resource)
					fmt.Printf("Key: %s Output: %s Looking for %s %s\n", filter.Key, o, filter.Op, filter.Value)
					fmt.Printf("ResourceId: %s %s\n",
						resourceId,
						isValid(o, filter.Op, filter.Value, rule.Severity))
				} else {
					fmt.Printf("Skipping rule %s for %s %s\n", rule.Id, resourceId, resourceTypes[resourceId])
				}
			}
		}
	}
}

func main() {
	parseCloudFormation := flag.Bool("cloudformation", false, "Validate CloudFormation template")
	parseTerraform := flag.Bool("terraform", false, "Validate Terraform template")
	flag.Parse()

	if *parseCloudFormation {
		cloudFormation()
	}
	if *parseTerraform {
		terraform()
	}
}
