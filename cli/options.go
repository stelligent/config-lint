package main

import (
	"fmt"
	"github.com/ghodss/yaml"
	"github.com/stelligent/config-lint/assertion"
	"io/ioutil"
	"path/filepath"
	"strings"
)

func getLinterOptions(o CommandLineOptions, p ProfileOptions) (LinterOptions, error) {
	allExcludePatterns, err := loadExcludePatterns(o.ExcludePatterns, o.ExcludeFromFilenames)
	if err != nil {
		return LinterOptions{}, err
	}
	linterOptions := LinterOptions{
		Tags:             makeTagList(*o.Tags, p.Tags),
		RuleIDs:          makeRulesList(*o.Ids, p.IDs),
		IgnoreRuleIDs:    makeRulesList(*o.IgnoreIds, p.IgnoreIDs),
		QueryExpression:  makeQueryExpression(*o.QueryExpression, *o.VerboseReport, p.Query),
		SearchExpression: *o.SearchExpression,
		ExcludePatterns:  allExcludePatterns,
		Variables:        mergeVariables(p.Variables, parseVariables(o.Variables)),
	}
	return linterOptions, nil
}

func loadProfile(filename string) (ProfileOptions, error) {
	defaultFilename := "config-lint.yml"
	var options ProfileOptions
	if filename == "" {
		filename = defaultFilename
	}
	bb, err := ioutil.ReadFile(filename)
	if err != nil {
		if filename == defaultFilename {
			return options, nil
		}
		return options, err
	}
	err = yaml.Unmarshal(bb, &options)
	if err != nil {
		return options, err
	}
	if len(options.Files) > 0 {
		patterns := options.Files
		options.Files = []string{}
		for _, pattern := range patterns {
			matches, err := filepath.Glob(pattern)
			if err != nil {
				return options, err
			}
			options.Files = append(options.Files, matches...)
		}
	}
	return options, nil
}

func makeTagList(tags string, profileOptions []string) []string {
	if tags != "" {
		return strings.Split(tags, ",")
	}
	if len(profileOptions) != 0 {
		return profileOptions
	}
	return nil
}

func makeRulesList(ruleIDs string, profileOptions []string) []string {
	if ruleIDs != "" {
		return strings.Split(ruleIDs, ",")
	}
	if len(profileOptions) != 0 {
		return profileOptions
	}
	return nil
}

func makeQueryExpression(queryExpression string, verboseReport bool, profileOptions string) string {
	if queryExpression != "" {
		return queryExpression
	}
	// return complete report when -verbose option is used
	if verboseReport {
		return ""
	}
	if profileOptions != "" {
		return profileOptions
	}
	// default is to only report Violations
	return "Violations[]"
}

func parseVariables(vars []string) map[string]string {
	m := map[string]string{}
	for _, kv := range vars {
		parts := strings.Split(kv, "=")
		if len(parts) == 2 {
			m[parts[0]] = parts[1]
		} else {
			fmt.Printf("Cannot parse command line variable: %s\n", kv)
		}
	}
	return m
}

func mergeVariables(a, b map[string]string) map[string]string {
	if a == nil {
		return b
	}
	if b == nil {
		return map[string]string{}
	}
	for k, v := range b {
		a[k] = v
	}
	return a
}

func loadExcludePatterns(patterns []string, excludeFromFilenames []string) ([]string, error) {
	if len(excludeFromFilenames) == 0 {
		return patterns, nil
	}
	for _, filename := range excludeFromFilenames {
		lines, err := ioutil.ReadFile(filename)
		if err != nil {
			return patterns, err
		}
		for _, patternFromFile := range strings.Split(string(lines), "\n") {
			if patternFromFile != "" {
				assertion.Debugf("Pattern from file %s: %s\n", filename, patternFromFile)
				patterns = append(patterns, patternFromFile)
			}
		}
	}
	return patterns, nil
}
