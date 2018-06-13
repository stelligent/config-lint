package linter

import (
	"fmt"
	"github.com/stelligent/config-lint/assertion"
	"io"
	"os"
	"path/filepath"
)

type (
	// Linter provides the interface for all supported linters
	Linter interface {
		Validate(ruleSet assertion.RuleSet, options Options) (assertion.ValidationReport, error)
		Search(ruleSet assertion.RuleSet, searchExpression string, w io.Writer)
	}

	// Options configures what resources will be linted
	Options struct {
		Tags          []string
		RuleIDs       []string
		IgnoreRuleIDs []string
	}
)

// NewLinter create the right kind of Linter based on the type argument
func NewLinter(ruleSet assertion.RuleSet, args []string) (Linter, error) {
	vs := assertion.StandardValueSource{}
	filenames := getFilenames(args)
	switch ruleSet.Type {
	case "Kubernetes":
		return FileLinter{Filenames: filenames, ValueSource: vs, Loader: KubernetesResourceLoader{}}, nil
	case "Terraform":
		return FileLinter{Filenames: filenames, ValueSource: vs, Loader: TerraformResourceLoader{}}, nil
	case "SecurityGroup":
		return AWSResourceLinter{Loader: SecurityGroupLoader{}, ValueSource: vs}, nil
	case "IAMUser":
		return AWSResourceLinter{Loader: IAMUserLoader{}, ValueSource: vs}, nil
	case "LintRules":
		return FileLinter{Filenames: filenames, ValueSource: vs, Loader: RulesResourceLoader{}}, nil
	case "YAML":
		return FileLinter{Filenames: filenames, ValueSource: vs, Loader: YAMLResourceLoader{Resources: ruleSet.Resources}}, nil
	default:
		return nil, fmt.Errorf("Type not supported: %s", ruleSet.Type)
	}
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
