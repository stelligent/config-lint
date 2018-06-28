package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/ghodss/yaml"
	"github.com/stelligent/config-lint/assertion"
	"github.com/stelligent/config-lint/linter"
)

var version string

type (
	// LinterOptions for applying rules
	LinterOptions struct {
		Tags             []string
		RuleIDs          []string
		IgnoreRuleIDs    []string
		QueryExpression  string
		SearchExpression string
		ExcludePatterns  []string
		Variables        map[string]string
	}

	// ProfileOptions for default options from a project file
	ProfileOptions struct {
		Rules      []string
		IDs        []string
		IgnoreIDs  []string `json:"ignore_ids"`
		Tags       []string
		Query      string
		Files      []string
		Terraform  bool
		Exceptions []RuleException
		Variables  map[string]string
	}

	// RuleException optional list allowing a project to ignore specific rules for specific resources
	RuleException struct {
		RuleID           string
		ResourceCategory string
		ResourceType     string
		ResourceID       string
		Comments         string
	}

	// CommandLineOptions for collecting options from the command line
	CommandLineOptions struct {
		RulesFilenames        arrayFlags
		ExcludePatterns       arrayFlags
		ExcludeFromFilenames  arrayFlags
		Variables             arrayFlags
		ProfileFilename       *string
		TerraformBuiltInRules *bool
		Tags                  *string
		Ids                   *string
		IgnoreIds             *string
		QueryExpression       *string
		VerboseReport         *bool
		SearchExpression      *string
		Validate              *bool
		Version               *bool
		Debug                 *bool
	}
)

//go:generate go-bindata -pkg $GOPACKAGE -o assets.go assets/
func main() {
	commandLineOptions := CommandLineOptions{}
	commandLineOptions.TerraformBuiltInRules = flag.Bool("terraform", false, "Use built-in rules for Terraform")
	flag.Var(&commandLineOptions.RulesFilenames, "rules", "Rules file, can be specified multiple times")
	commandLineOptions.Tags = flag.String("tags", "", "Run only tests with tags in this comma separated list")
	commandLineOptions.Ids = flag.String("ids", "", "Run only the rules in this comma separated list")
	commandLineOptions.IgnoreIds = flag.String("ignore-ids", "", "Ignore the rules in this comma separated list")
	commandLineOptions.QueryExpression = flag.String("query", "", "JMESPath expression to query the results")
	commandLineOptions.VerboseReport = flag.Bool("verbose", false, "Output a verbose report")
	commandLineOptions.SearchExpression = flag.String("search", "", "JMESPath expression to evaluation against the files")
	commandLineOptions.Validate = flag.Bool("validate", false, "Validate rules file")
	commandLineOptions.Version = flag.Bool("version", false, "Get program version")
	commandLineOptions.ProfileFilename = flag.String("profile", "", "Provide default options")
	flag.Var(&commandLineOptions.ExcludePatterns, "exclude", "Filename patterns to exclude")
	flag.Var(&commandLineOptions.ExcludeFromFilenames, "exclude-from", "Filename containing patterns to exclude")
	flag.Var(&commandLineOptions.Variables, "var", "Variable values for rules with ValueFrom.Variable")
	commandLineOptions.Debug = flag.Bool("debug", false, "Debug logging")

	flag.Parse()

	if *commandLineOptions.Version == true {
		fmt.Println(version)
		return
	}

	if *commandLineOptions.Debug == true {
		assertion.SetDebug(true)
	}

	if *commandLineOptions.Validate {
		validateRules(flag.Args())
		return
	}

	profileOptions, err := loadProfile(*commandLineOptions.ProfileFilename)
	if err != nil {
		fmt.Printf("Error loading profile: %v\n", err)
		return
	}

	rulesFilenames := loadFilenames(commandLineOptions.RulesFilenames, profileOptions.Rules)
	configFilenames := loadFilenames(flag.Args(), profileOptions.Files)
	useTerraformBuiltInRules := *commandLineOptions.TerraformBuiltInRules || profileOptions.Terraform

	if err != nil {
		fmt.Printf("Unable to load exclude patterns: %s\n", err)
		return
	}

	linterOptions, err := getLinterOptions(commandLineOptions, profileOptions)
	if err != nil {
		fmt.Printf("Failed to parse options: %v\n", err)
		return
	}

	ruleSets, err := loadRuleSets(rulesFilenames)
	if err != nil {
		fmt.Printf("Failed to load rules: %v\n", err)
		return
	}
	ruleSets = addExceptions(ruleSets, profileOptions.Exceptions)
	if useTerraformBuiltInRules {
		builtInRuleSet, err := loadBuiltInRuleSet("assets/terraform.yml")
		if err != nil {
			fmt.Printf("Failed to load built-in rules for Terraform: %v\n", err)
			return
		}
		ruleSets = append(ruleSets, builtInRuleSet)
	}
	applyRules(ruleSets, configFilenames, linterOptions)
}

func addExceptions(ruleSets []assertion.RuleSet, exceptions []RuleException) []assertion.RuleSet {
	sets := []assertion.RuleSet{}
	for _, ruleSet := range ruleSets {
		sets = append(sets, addExceptionsToRuleSet(ruleSet, exceptions))
	}
	return sets
}

func addExceptionsToRuleSet(ruleSet assertion.RuleSet, exceptions []RuleException) assertion.RuleSet {
	rules := []assertion.Rule{}
	for _, rule := range ruleSet.Rules {
		for _, e := range exceptions {
			if rule.ID == e.RuleID &&
				rule.Resource == e.ResourceType &&
				(rule.Category == e.ResourceCategory || rule.Category == "") {
				rule.Except = append(rule.Except, e.ResourceID)
			}
		}
		rules = append(rules, rule)
	}
	ruleSet.Rules = rules
	return ruleSet
}

func validateRules(filenames []string) {
	builtInRuleSet, err := loadBuiltInRuleSet("assets/lint-rules.yml")
	if err != nil {
		fmt.Printf("Unable to load build-in rules for validation: %v\n", err)
		return
	}
	ruleSets := []assertion.RuleSet{builtInRuleSet}
	linterOptions := LinterOptions{
		QueryExpression: "Violations[]",
	}
	applyRules(ruleSets, filenames, linterOptions)
}

func loadRuleSets(args arrayFlags) ([]assertion.RuleSet, error) {
	rulesFilenames := yamlFilesOnly(getFilenames(args))
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

func yamlFilesOnly(filenames []string) []string {
	configFiles := []string{}
	configPatterns := []string{"*yml", "*.yaml"}
	for _, filename := range filenames {
		match, _ := assertion.ShouldIncludeFile(configPatterns, filename)
		if match {
			configFiles = append(configFiles, filename)
		}
	}
	return configFiles
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

func applyRules(ruleSets []assertion.RuleSet, args arrayFlags, options LinterOptions) {

	report := assertion.ValidationReport{
		Violations:       []assertion.Violation{},
		FilesScanned:     []string{},
		ResourcesScanned: []assertion.ScannedResource{},
	}

	filenames := excludeFilenames(getFilenames(args), options.ExcludePatterns)
	vs := assertion.StandardValueSource{Variables: options.Variables}

	for _, ruleSet := range ruleSets {
		l, err := linter.NewLinter(ruleSet, vs, filenames)
		if err != nil {
			fmt.Println(err)
			return
		}
		if l != nil {
			if options.SearchExpression != "" {
				l.Search(ruleSet, options.SearchExpression, os.Stdout)
			} else {
				options := linter.Options{
					Tags:          options.Tags,
					RuleIDs:       options.RuleIDs,
					IgnoreRuleIDs: options.IgnoreRuleIDs,
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

func loadFilenames(commandLineOptions []string, profileOptions []string) []string {
	if len(commandLineOptions) > 0 {
		return commandLineOptions
	}
	return profileOptions
}

func excludeFilenames(filenames []string, excludePatterns []string) []string {
	assertion.Debugf("Exclude patterns: %v\n", excludePatterns)
	filteredFilenames := []string{}
	for _, filename := range filenames {
		if !excludeFilename(filename, excludePatterns) {
			filteredFilenames = append(filteredFilenames, filename)
		}
	}
	return filteredFilenames
}

func excludeFilename(filename string, excludePatterns []string) bool {
	for _, pattern := range excludePatterns {
		match, _ := filepath.Match(pattern, filename)
		if match {
			assertion.Debugf("Excluding file: %s using pattern: %s\n", filename, pattern)
			return true
		}
	}
	return false
}

func getFilenames(args []string) []string {
	filenames := []string{}
	for _, arg := range args {
		fi, err := os.Stat(arg)
		if err != nil {
			fmt.Printf("Cannot open %s\n", arg)
			continue
		}
		mode := fi.Mode()
		if mode.IsDir() {
			filenames = append(filenames, getFilesInDirectory(arg)...)
		} else {
			filenames = append(filenames, arg)
		}
	}
	return filenames
}

func getFilesInDirectory(root string) []string {
	directoryFiles := []string{}
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			fmt.Printf("Error processing %s: %s\n", path, err)
			return err
		}
		if !info.IsDir() {
			directoryFiles = append(directoryFiles, path)
		}
		return nil
	})
	if err != nil {
		fmt.Printf("Error walking directory %s: %s\n", root, err)
	}
	return directoryFiles
}
