package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/ghodss/yaml"
	"github.com/stelligent/config-lint/assertion"
	"github.com/stelligent/config-lint/linter"
	"os"
	"strings"
)

type (
	// ApplyOptions for applying rules
	ApplyOptions struct {
		Tags             []string
		RuleIDs          []string
		QueryExpression  string
		SearchExpression string
	}
)

func main() {
	var rulesFilenames arrayFlags
	verboseLogging := flag.Bool("verbose", false, "Verbose logging")
	flag.Var(&rulesFilenames, "rules", "Rules file, can be specified multiple times")
	tags := flag.String("tags", "", "Run only tests with tags in this comma separated list")
	ids := flag.String("ids", "", "Run only the rules in this comma separated list")
	queryExpression := flag.String("query", "", "JMESPath expression to query the results")
	searchExpression := flag.String("search", "", "JMESPath expression to evaluation against the files")
	validate := flag.Bool("validate", false, "Validate rules file")
	flag.Parse()

	if *verboseLogging == true {
		assertion.SetVerbose(true)
	}

	if *validate {
		validateRules(flag.Args())
		return
	}

	applyOptions := ApplyOptions{
		Tags:             makeTagList(*tags),
		RuleIDs:          makeRulesList(*ids),
		QueryExpression:  *queryExpression,
		SearchExpression: *searchExpression,
	}
	applyRules(rulesFilenames, flag.Args(), applyOptions)

}

func validateRules(filenames []string) {
	rulesFilenames := []string{"./builtin-rules/lint-rules.yml"} // FIXME how to embed this file in the binary?
	applyOptions := ApplyOptions{
		QueryExpression: "Violations[]",
	}
	applyRules(rulesFilenames, filenames, applyOptions)
}

func applyRules(rulesFilenames arrayFlags, args arrayFlags, options ApplyOptions) {

	report := assertion.ValidationReport{
		Violations:       []assertion.Violation{},
		FilesScanned:     []string{},
		ResourcesScanned: []assertion.ScannedResource{},
	}

	for _, rulesFilename := range rulesFilenames {
		rulesContent, err := assertion.LoadRules(rulesFilename)
		if err != nil {
			fmt.Println("Unable to load rules from:" + rulesFilename)
			fmt.Println(err.Error())
		}
		ruleSet, err := assertion.ParseRules(rulesContent)
		if err != nil {
			fmt.Println("Unable to parse rules in:" + rulesFilename)
			fmt.Println(err.Error())
		}
		l, err := linter.NewLinter(ruleSet.Type, args)
		if err != nil {
			fmt.Println(err)
			return
		}
		if l != nil {
			if options.SearchExpression != "" {
				l.Search(ruleSet, options.SearchExpression)
			} else {
				options := linter.Options{
					Tags:    options.Tags,
					RuleIDs: options.RuleIDs,
				}
				r, err := l.Validate(ruleSet, options)
				if err != nil {
					fmt.Println("Validate failed:", err)
				}
				report = linter.CombineValidationReports(report, r)
			}
		}
	}
	if options.SearchExpression == "" {
		err := printReport(report, options.QueryExpression)
		if err != nil {
			fmt.Println(err.Error())
		}
	}
	os.Exit(generateExitCode(report))
}

func printReport(report assertion.ValidationReport, queryExpression string) error {
	jsonData, err := json.MarshalIndent(report, "", "  ")
	if err != nil {
		return err
	}
	if queryExpression != "" {
		var data interface{}
		err = yaml.Unmarshal(jsonData, &data)
		if err != nil {
			return err
		}
		v, err := assertion.SearchData(queryExpression, data)
		if err != nil {
			return err
		}
		s, err := assertion.JSONStringify(v)
		if err == nil && s != "null" {
			fmt.Println(s)
		}
	} else {
		fmt.Println(string(jsonData))
	}
	return nil
}

func makeTagList(tags string) []string {
	if tags == "" {
		return nil
	}
	return strings.Split(tags, ",")
}

func makeRulesList(ruleIDs string) []string {
	if ruleIDs == "" {
		return nil
	}
	return strings.Split(ruleIDs, ",")
}

type arrayFlags []string

func (i *arrayFlags) String() string {
	if i != nil {
		return strings.Join(*i, ",")
	}
	return ""
}

func (i *arrayFlags) Set(value string) error {
	*i = append(*i, value)
	return nil
}

func generateExitCode(report assertion.ValidationReport) int {
	for _, v := range report.Violations {
		if v.Status == "FAILURE" {
			return 1
		}
	}
	return 0
}
