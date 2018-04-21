package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/ghodss/yaml"
	"github.com/stelligent/config-lint/assertion"
	"github.com/stelligent/config-lint/linter"
)

var version string

type (
	// ApplyOptions for applying rules
	ApplyOptions struct {
		Tags             []string
		RuleIDs          []string
		QueryExpression  string
		SearchExpression string
	}
)

//go:generate go-bindata -pkg $GOPACKAGE -o assets.go assets/
func main() {
	var rulesFilenames arrayFlags
	terraformBuiltInRules := flag.Bool("terraform", false, "Use built-in rules for Terraform")
	debug := flag.Bool("debug", false, "Debug logging")
	flag.Var(&rulesFilenames, "rules", "Rules file, can be specified multiple times")
	tags := flag.String("tags", "", "Run only tests with tags in this comma separated list")
	ids := flag.String("ids", "", "Run only the rules in this comma separated list")
	queryExpression := flag.String("query", "", "JMESPath expression to query the results")
	verboseReport := flag.Bool("verbose", false, "Output a verbose report")
	searchExpression := flag.String("search", "", "JMESPath expression to evaluation against the files")
	validate := flag.Bool("validate", false, "Validate rules file")
	versionFlag := flag.Bool("version", false, "Get program version")
	flag.Parse()

	if *versionFlag == true {
		fmt.Println(version)
		return
	}

	if *debug == true {
		assertion.SetDebug(true)
	}

	if *validate {
		validateRules(flag.Args())
		return
	}

	applyOptions := ApplyOptions{
		Tags:             makeTagList(*tags),
		RuleIDs:          makeRulesList(*ids),
		QueryExpression:  makeQueryExpression(*queryExpression, *verboseReport),
		SearchExpression: *searchExpression,
	}
	ruleSets, err := loadRuleSets(rulesFilenames)
	if err != nil {
		fmt.Printf("Failed to load rules: %v\n", err)
		return
	}
	if *terraformBuiltInRules {
		builtInRuleSet, err := loadBuiltInRuleSet("assets/terraform.yml")
		if err != nil {
			fmt.Printf("Failed to load built-in rules for Terraform: %v\n", err)
			return
		}
		ruleSets = append(ruleSets, builtInRuleSet)
	}
	applyRules(ruleSets, flag.Args(), applyOptions)
}

func validateRules(filenames []string) {
	builtInRuleSet, err := loadBuiltInRuleSet("assets/lint-rules.yml")
	if err != nil {
		fmt.Printf("Unable to load build-in rules for validation: %v\n", err)
		return
	}
	ruleSets := []assertion.RuleSet{builtInRuleSet}
	applyOptions := ApplyOptions{
		QueryExpression: "Violations[]",
	}
	applyRules(ruleSets, filenames, applyOptions)
}

func loadRuleSets(rulesFilenames arrayFlags) ([]assertion.RuleSet, error) {
	ruleSets := []assertion.RuleSet{}
	for _, rulesFilename := range rulesFilenames {
		rulesContent, err := ioutil.ReadFile(rulesFilename)
		if err != nil {
			fmt.Println("Unable to load rules from:" + rulesFilename)
			fmt.Println(err.Error())
			return ruleSets, err
		}
		ruleSet, err := assertion.ParseRules(string(rulesContent))
		if err != nil {
			fmt.Println("Unable to parse rules in:" + rulesFilename)
			fmt.Println(err.Error())
			return ruleSets, err
		}
		ruleSets = append(ruleSets, ruleSet)
	}
	return ruleSets, nil
}

func loadBuiltInRuleSet(name string) (assertion.RuleSet, error) {
	emptyRuleSet := assertion.RuleSet{}
	rulesContent, err := Asset(name)
	if err != nil {
		fmt.Println("Unable to find built-in rules:", name)
		fmt.Println(err.Error())
		return emptyRuleSet, err
	}
	ruleSet, err := assertion.ParseRules(string(rulesContent))
	if err != nil {
		fmt.Println("Unable to parse built-in rules:" + name)
		fmt.Println(err.Error())
		return emptyRuleSet, err
	}
	return ruleSet, nil
}

func applyRules(ruleSets []assertion.RuleSet, args arrayFlags, options ApplyOptions) {

	report := assertion.ValidationReport{
		Violations:       []assertion.Violation{},
		FilesScanned:     []string{},
		ResourcesScanned: []assertion.ScannedResource{},
	}

	for _, ruleSet := range ruleSets {
		l, err := linter.NewLinter(ruleSet, args)
		if err != nil {
			fmt.Println(err)
			return
		}
		if l != nil {
			if options.SearchExpression != "" {
				l.Search(ruleSet, options.SearchExpression, os.Stdout)
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

func makeQueryExpression(queryExpression string, verboseReport bool) string {
	if queryExpression != "" {
		return queryExpression
	}
	// return complete report when -verbose option is used
	if verboseReport {
		return ""
	}
	// default is to only report Violations
	return "Violations[]"
}
