# Running config-lint

The program has a set of built-in rules for scanning the following types of files:

* [Terraform](terraform.md)

The program can also read files from a separate YAML file, and can scan these types of files:

* [Terraform](terraform.md)
* Kubernetes
* LintRules
* YAML
* JSON

## Example invocations

### Validate Terraform files with built-in rules

```
config-lint -terraform example-files/config
```

### Validate Terraform files with custom rules

```
config-lint -rules examples-files/rules/terraform.yml example-files/config
```

### Validate Kubernetes files

```
config-lint -rules example-files/rules/kubernetes.yml example-files/config
```

### Validate LintRules files

This type of linting allows the tool to lint its own rules.

```
config-lint -rules example-files/rules/lint-rules.yml example-files/rules
```

### Validate a custom YAML file

```
config-lint -rules example-files/rules/generic-yaml.yml example-files/config/generic.config
```

## Using STDIN

You can use "-" for the filename if you want the configuration data read from STDIN.

```
cat example-files/resources/s3.tf | config-lint -terraform -
```

## Exit Code

If at least one rule with a severity of FAILURE was triggered the exit code will be 1, otherwise it will be 0.

## Options

Here are all the different command line options that can be used with config-lint. You can also
view them via the -help option.

 * -debug - Debug logging
    	
 * -exclude value - Filename patterns to exclude
 
 * -exclude-from value - Filename containing patterns to exclude
 
 * -ids string - Run only the rules in this comma separated list
 
 * -ignore-ids string - Ignore the rules in this comma separated list
 
 * -profile string- Provide default options
 
 * -query string - JMESPath expression to query the results
 
 * -rules value - Rules file, can be specified multiple times
 
 * -search string - JMESPath expression to evaluation against the files
 
 * -tags string - Run only tests with tags in this comma separated list
 
 * -terraform - Use built-in rules for Terraform
 
 * -validate - Validate rules file
 
 * -var value - Variable values for rules with ValueFrom.Variable
 
 * -verbose - Output a verbose report
 
 * -version - Get program version
 
 * -tfparser - (Optional) Set the Terraform parser version. Options are `tf11` or `tf12`. By default, `tf12` will be used.