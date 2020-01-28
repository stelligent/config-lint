package linter

import (
	"fmt"
	"io"
	"time"

	"github.com/stelligent/config-lint/assertion"
)

type (
	// Variable contains a key value pair for expressions in a Terraform configuration file
	Variable struct {
		Name  string
		Value interface{}
	}
	// FileResources contains the variables and resources loaded from a collection of files
	FileResources struct {
		Resources []assertion.Resource
		Variables []Variable
	}
	// FileResourceLoader provides the interface that a Linter needs to load a collection of Resource objects
	FileResourceLoader interface {
		Load(filename string) (FileResources, error)
		PostLoad(resources FileResources) ([]assertion.Resource, error)
	}
	// FileLinter provides implementation for some common functions that are used by multiple Linter implementations
	FileLinter struct {
		Filenames   []string
		ValueSource assertion.ValueSource
		Loader      FileResourceLoader
	}
)

// Validate validates a collection of filenames using a RuleSet
func (fl FileLinter) Validate(ruleSet assertion.RuleSet, options Options) (assertion.ValidationReport, error) {
	rules := assertion.FilterRulesByTagAndID(ruleSet.Rules, options.Tags, options.RuleIDs, options.IgnoreRuleIDs)
	rl := ResourceLinter{ValueSource: fl.ValueSource}

	resources := []assertion.Resource{}
	variables := []Variable{}
	filesScanned := []string{}

	loadViolations := []assertion.Violation{}
	var resourcesToValidate []assertion.Resource

	//TODO: This is ugly in several ways
	if tf12Loader, ok := fl.Loader.(Terraform12ResourceLoader); ok {
		filteredFilenames := filterFiles(fl.Filenames, ruleSet.Files)
		result, err := tf12Loader.LoadMany(filteredFilenames)
		if err != nil {
			//TODO: It would probably be nice if we mapped this back to the correct file
			loadViolations = append(loadViolations, makeLoadViolation(fl.Filenames[0], err))
		}
		resourcesToValidate = result.Resources
		filesScanned = append(filesScanned, fl.Filenames...)
	} else {
		filesScanned, loadViolations, resources, variables = iterateFiles(fl, ruleSet, filesScanned, loadViolations, resources, variables)
		var err error
		resourcesToValidate, err = fl.Loader.PostLoad(FileResources{Resources: resources, Variables: variables})
		if err != nil {
			return assertion.ValidationReport{}, err
		}
	}

	report, err := rl.ValidateResources(resourcesToValidate, rules)
	if err != nil {
		return report, err
	}
	report.FilesScanned = filesScanned
	report.Violations = append(report.Violations, loadViolations...)
	return report, nil
}

func filterFiles(fileNames []string, patterns []string) (ret []string) {
	for _, s := range fileNames {
		ok, _ := assertion.ShouldIncludeFile(patterns, s)
		if ok {
			ret = append(ret, s)
		}
	}
	return
}

func iterateFiles(fl FileLinter, ruleSet assertion.RuleSet, filesScanned []string, loadViolations []assertion.Violation, resources []assertion.Resource, variables []Variable) ([]string, []assertion.Violation, []assertion.Resource, []Variable) {
	for _, filename := range fl.Filenames {
		include, err := assertion.ShouldIncludeFile(ruleSet.Files, filename)
		if err == nil && include {
			assertion.Debugf("Processing %s\n", filename)
			filesScanned = append(filesScanned, filename)
			loaded, err := fl.Loader.Load(filename)
			if err != nil {
				loadViolations = append(loadViolations, makeLoadViolation(filename, err))
				continue
			}
			assertion.Debugf("Found variables %v\n", loaded.Variables)
			resources = append(resources, loaded.Resources...)
			variables = append(variables, loaded.Variables...)
		}
	}
	return filesScanned, loadViolations, resources, variables
}

func makeLoadViolation(filename string, err error) assertion.Violation {
	return assertion.Violation{
		RuleID:           "FILE_LOAD",
		ResourceID:       filename,
		ResourceType:     "file",
		Category:         "load",
		Status:           "FAILURE",
		RuleMessage:      "Unable to load file",
		AssertionMessage: err.Error(),
		Filename:         filename,
		CreatedAt:        time.Now().UTC().Format(time.RFC3339),
	}
}

// Search evaluates a JMESPath expression against resources in a collection of filenames
func (fl FileLinter) Search(ruleSet assertion.RuleSet, searchExpression string, w io.Writer) {
	resources := []assertion.Resource{}
	variables := []Variable{}
	for _, filename := range fl.Filenames {
		include, _ := assertion.ShouldIncludeFile(ruleSet.Files, filename) // FIXME what about error?
		if include {
			fmt.Fprintf(w, "Searching %s:\n", filename)
			loaded, err := fl.Loader.Load(filename)
			if err != nil {
				fmt.Fprintln(w, "Error for file:", filename)
				fmt.Fprintln(w, err.Error())
			}
			resources = append(resources, loaded.Resources...)
			variables = append(variables, loaded.Variables...)
		}
	}
	resourcesToSearch, err := fl.Loader.PostLoad(FileResources{Resources: resources, Variables: variables})
	if err != nil {
		fmt.Fprintln(w, err.Error())
		return
	}
	for _, resource := range resourcesToSearch {
		v, err := assertion.SearchData(searchExpression, resource.Properties)
		if err != nil {
			fmt.Fprintln(w, err)
		} else {
			s, err := assertion.JSONStringify(v)
			if err != nil {
				fmt.Fprintln(w, err)
			} else {
				fmt.Fprintf(w, "%s (%s): %s\n", resource.ID, resource.Type, s)
			}
		}
	}
}
