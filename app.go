package main

import (
	"encoding/json"
	"flag"
	"fmt"
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
	Warnings      []Violation
	Failures      []Violation
	AllViolations []Violation
	FilesScanned  []string
}

func printReport(report ValidationReport, queryExpression string) int {
	if queryExpression != "" {
		v := searchData(queryExpression, report)
		if v != "null" {
			fmt.Println(v)
		}
	} else {
		jsonData, err := json.MarshalIndent(report, "", "  ")
		if err != nil {
			panic(err)
		}
		fmt.Println(string(jsonData))
	}
	if len(report.Failures) > 0 {
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

func main() {
	kubernetesFiles := flag.Bool("kubernetes", false, "Process kubernetes files")
	terraformFiles := flag.Bool("terraform", false, "Process terraform files")
	verboseLogging := flag.Bool("verbose", false, "Verbose logging")
	rulesFilename := flag.String("rules", "./rules/terraform.yml", "Rules file")
	tags := flag.String("tags", "", "Run only tests with tags in this comma separated list")
	ids := flag.String("ids", "", "Run only the rules in this comma separated list")
	queryExpression := flag.String("query", "", "JMESPath expression to query the results")
	searchExpression := flag.String("search", "", "JMESPath expression to evaluation against the files")
	flag.Parse()

	logger := makeLogger(*verboseLogging)

	exitCode := 0

	if *kubernetesFiles {
		if *searchExpression != "" {
			kubernetesSearch(flag.Args(), *searchExpression, logger)
		} else {
			report := kubernetes(flag.Args(), *rulesFilename, makeTagList(*tags), makeRulesList(*ids), logger)
			exitCode = printReport(report, *queryExpression)
		}
	}
	if *terraformFiles {
		if *searchExpression != "" {
			terraformSearch(flag.Args(), *searchExpression, logger)
		} else {
			report := terraform(flag.Args(), *rulesFilename, makeTagList(*tags), makeRulesList(*ids), logger)
			exitCode = printReport(report, *queryExpression)
		}
	}
	os.Exit(exitCode)
}
