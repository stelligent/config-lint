package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/ghodss/yaml"
	"os"
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

type RuleSet struct {
	Type        string
	Description string
	Files       []string
	Rules       []Rule
	Version     string
}

type Violation struct {
	RuleId       string
	ResourceId   string
	ResourceType string
	Status       string
	Message      string
	Filename     string
}

type ValidationReport struct {
	Violations   map[string]([]Violation)
	FilesScanned []string
}

type Linter interface {
	Validate(filenames []string, ruleSet RuleSet, tags []string, ruleIds []string) ValidationReport
	Search(filenames []string, searchExpression string)
}

func printReport(report ValidationReport, queryExpression string) int {
	jsonData, err := json.MarshalIndent(report, "", "  ")
	if err != nil {
		panic(err)
	}
	if queryExpression != "" {
		var data interface{}
		err = yaml.Unmarshal(jsonData, &data)
		if err != nil {
			panic(err)
		}
		v, err := searchData(queryExpression, data)
		if err == nil && v != "null" {
			fmt.Println(v)
		}
	} else {
		fmt.Println(string(jsonData))
	}
	if len(report.Violations["FAILURE"]) > 0 {
		return 1
	}
	return 0
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

func makeLinter(linterType string, log LoggingFunction) Linter {
	switch linterType {
	case "Kubernetes":
		return KubernetesLinter{Log: log}
	case "Terraform":
		return TerraformLinter{Log: log}
	default:
		fmt.Printf("Type not supported: %s\n", linterType)
		return nil
	}
}

func main() {
	verboseLogging := flag.Bool("verbose", false, "Verbose logging")
	rulesFilename := flag.String("rules", "./rules/terraform.yml", "Rules file")
	tags := flag.String("tags", "", "Run only tests with tags in this comma separated list")
	ids := flag.String("ids", "", "Run only the rules in this comma separated list")
	queryExpression := flag.String("query", "", "JMESPath expression to query the results")
	searchExpression := flag.String("search", "", "JMESPath expression to evaluation against the files")
	flag.Parse()

	exitCode := 0

	ruleSet := MustParseRules(loadRules(*rulesFilename))
	linter := makeLinter(ruleSet.Type, makeLogger(*verboseLogging))
	if linter != nil {
		if *searchExpression != "" {
			linter.Search(flag.Args(), *searchExpression)
		} else {
			report := linter.Validate(flag.Args(), ruleSet, makeTagList(*tags), makeRulesList(*ids))
			exitCode = printReport(report, *queryExpression)
		}
	}
	os.Exit(exitCode)
}
