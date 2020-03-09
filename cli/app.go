package main

//go:generate packr

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/ghodss/yaml"
	"github.com/gobuffalo/packr"
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
		TerraformParser  string
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
		TerraformParser       *string
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
		Args                  []string
	}

	// ReportWriter formats and displays a ValidationReport
	ReportWriter interface {
		WriteReport(assertion.ValidationReport, LinterOptions)
	}

	// DefaultReportWriter writes the report to Stdout
	DefaultReportWriter struct {
		Writer io.Writer
	}
)

func main() {

	commandLineOptions := getCommandLineOptions()

	if *commandLineOptions.Version == true {
		fmt.Println(version)
		return
	}

	if *commandLineOptions.Debug == true {
		assertion.SetDebug(true)
	}

	if *commandLineOptions.Validate {
		exitCode, err := validateRules(commandLineOptions.Args, DefaultReportWriter{Writer: os.Stdout})
		if err != nil {
			fmt.Println(err.Error())
		}
		os.Exit(exitCode)
	}

	profileOptions, err := loadProfile(*commandLineOptions.ProfileFilename)
	if err != nil {
		fmt.Printf("Error loading profile: %v\n", err)
		os.Exit(-1)
	}

	rulesFilenames := loadFilenames(commandLineOptions.RulesFilenames, profileOptions.Rules)
	configFilenames := defaultToCurrentDirectory(loadFilenames(commandLineOptions.Args, profileOptions.Files))
	useTerraformBuiltInRules := *commandLineOptions.TerraformBuiltInRules || profileOptions.Terraform

	if err != nil {
		fmt.Printf("Unable to load exclude patterns: %s\n", err)
		os.Exit(-1)
	}

	linterOptions, err := getLinterOptions(commandLineOptions, profileOptions)
	if err != nil {
		fmt.Printf("Failed to parse options: %v\n", err)
		os.Exit(-1)
	}

	ruleSets, err := loadRuleSets(rulesFilenames)
	if err != nil {
		fmt.Printf("Failed to load rules: %v\n", err)
		os.Exit(-1)
	}
	ruleSets = addExceptions(ruleSets, profileOptions.Exceptions)
	// Same rule set applies to both TerraformBuiltInRules and Terraform11BuiltInRules
	// Run for terraform12 by default
	if useTerraformBuiltInRules {
		builtInRuleSet, err := loadBuiltInRuleSet("terraform.yml")
		if err != nil {
			fmt.Printf("Failed to load built-in rules for Terraform: %v\n", err)
			os.Exit(-1)
		}
		ruleSets = append(ruleSets, builtInRuleSet)
	}
	if len(ruleSets) == 0 {
		fmt.Println("No rules")
		os.Exit(-1)
	}
	os.Exit(applyRules(ruleSets, configFilenames, linterOptions, DefaultReportWriter{Writer: os.Stdout}))
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
			if rule.ID == e.RuleID && resourceMatch(rule, e) && categoryMatch(rule, e) {
				rule.Except = append(rule.Except, e.ResourceID)
			}
		}
		rules = append(rules, rule)
	}
	ruleSet.Rules = rules
	return ruleSet
}

func resourceMatch(rule assertion.Rule, exception RuleException) bool {
	if assertion.SliceContains(rule.Resources, exception.ResourceType) || rule.Resource == exception.ResourceType {
		return true
	}
	return false
}

func categoryMatch(rule assertion.Rule, exception RuleException) bool {
	return rule.Category == exception.ResourceCategory || exception.ResourceCategory == "resources" || rule.Category == ""
}

func validateRules(filenames []string, w ReportWriter) (int, error) {
	builtInRuleSet, err := loadBuiltInRuleSet("lint-rules.yml")
	if err != nil {
		return -1, err
	}
	ruleSets := []assertion.RuleSet{builtInRuleSet}
	linterOptions := LinterOptions{
		QueryExpression: "Violations[]",
	}
	return applyRules(ruleSets, filenames, linterOptions, w), nil
}

func loadRuleSets(args arrayFlags) ([]assertion.RuleSet, error) {
	rulesFilenames := yamlFilesOnly(getFilenames(args))
	ruleSets := []assertion.RuleSet{}
	for _, rulesFilename := range rulesFilenames {
		rulesContent, err := ioutil.ReadFile(rulesFilename)
		if err != nil {
			return ruleSets, err
		}
		ruleSet, err := assertion.ParseRules(string(rulesContent))
		if err != nil {
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
	box := packr.NewBox("./assets")
	rulesContent, err := box.FindString(name)
	if err != nil {
		return emptyRuleSet, err
	}
	ruleSet, err := assertion.ParseRules(string(rulesContent))
	if err != nil {
		return emptyRuleSet, err
	}
	return ruleSet, nil
}

func applyRules(ruleSets []assertion.RuleSet, args arrayFlags, options LinterOptions, w ReportWriter) int {

	report := assertion.ValidationReport{
		Violations:       []assertion.Violation{},
		FilesScanned:     []string{},
		ResourcesScanned: []assertion.ScannedResource{},
	}

	tfParser := options.TerraformParser
	filenames := excludeFilenames(getFilenames(args), options.ExcludePatterns)
	vs := assertion.StandardValueSource{Variables: options.Variables}

	for _, ruleSet := range ruleSets {
		l, err := linter.NewLinter(ruleSet, vs, filenames, tfParser)
		if err != nil {
			fmt.Println(err)
			return -1
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
	w.WriteReport(report, options)
	return generateExitCode(report)
}

func printReport(w io.Writer, report assertion.ValidationReport, queryExpression string) error {
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
			fmt.Fprintln(w, s)
		}
	} else {
		fmt.Fprintln(w, string(jsonData))
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
			return -1
		}
	}
	return 0
}

func loadFilenames(commandLineFilenames []string, profileFilenames []string) []string {
	if len(commandLineFilenames) > 0 {
		return commandLineFilenames
	}
	if len(profileFilenames) > 0 {
		return profileFilenames
	}
	return []string{}
}

func defaultToCurrentDirectory(filenames []string) []string {
	if len(filenames) == 0 {
		return []string{"."}
	}
	return filenames
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
		if arg == "-" {
			filenames = append(filenames, arg)
			continue
		}
		fi, err := os.Stat(arg)
		if err != nil {
			// append as is, error reported later when file cannot be opened
			filenames = append(filenames, arg)
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
