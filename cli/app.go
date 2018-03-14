package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/ghodss/yaml"
	"github.com/lhitchon/config-lint/assertion"
	"os"
	"strings"
)

type Linter interface {
	Validate(filenames []string, ruleSet assertion.RuleSet, tags []string, ruleIds []string) assertion.ValidationReport
	Search(filenames []string, searchExpression string)
}

func printReport(report assertion.ValidationReport, queryExpression string) int {
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
		v, err := assertion.SearchData(queryExpression, data)
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

func makeLinter(linterType string, log assertion.LoggingFunction) Linter {
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

	ruleSet := assertion.MustParseRules(assertion.LoadRules(*rulesFilename))
	linter := makeLinter(ruleSet.Type, assertion.MakeLogger(*verboseLogging))
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
