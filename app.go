package main

import (
	"flag"
	"fmt"
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

type ValidationResult struct {
	RuleId       string
	ResourceId   string
	ResourceType string
	Status       string
	Message      string
	Filename     string
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
	rulesFilename := flag.String("rules", "./rules/terraform.yml", "Rules file")
	tags := flag.String("tags", "", "Run only tests with tags in this comma separated list")
	ids := flag.String("ids", "", "Run only the rules in this comma separated list")
	flag.Parse()

	for _, filename := range flag.Args() {
		terraform(filename, *rulesFilename, makeTagList(*tags), makeRulesList(*ids), makeLogger(*verboseLogging))
	}
}
