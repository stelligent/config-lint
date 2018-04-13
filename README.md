[![Build Status](https://circleci.com/gh/stelligent/config-lint.svg?style=shield)](https://circleci.com/gh/stelligent/config-lint)

# config-lint

A command line tool to validate configurations using rules specified in a YAML file.
The data being validated can come from template files, such as a Terraform file.
There is also an example of a Linter that runs agains data returned from an AWS API call.

# Installation 
You can use [Homebrew](https://brew.sh/) to install the latest version:

```
brew tap stelligent/tap
brew install config-lint
```

Alternatively, you can install manually from the [releases](https://github.com/stelligent/config-lint/releases).

# Build Command Line tool

```
make config-lint
```

# Run

The program currently supports scanning of the following types of files:

* Terraform
* Kubernetes
* LintRules
* YAML

And also the scanning of information from AWS Descibe API calls for:

* Security Groups
* IAM Users

## Example invocations

### Validate Terraform files

```
./config-lint --rules example-files/rules/terraform.yml example-files/config/*
```

### Validate Kubernetes files

```
./config-lint --rules example-files/rules/kubernetes.yml example-files/config/*
```

### Validate LintRules files

This type of linting allows the tool to lint its own rules.

```
./config-lint --rules example-files/rules/lint-rules.yml example-files/rules/*
```

### Validate Existing Security Groups

```
./config-lint --rules example-files/rules/security-groups.yml
```

### Validate Existing IAM Users

```
./config-lint --rules example-files/rules/iam-users.yml
```

# Rules File

A YAML file that specifies what kinds of files to process, and what validations to perform.

[Documented Here](docs/rules.md)

## Examples

To test that an AWS instance type has one of two values:
```
Version: 1
Description: Example rules
Type: Terraform
Files:
  - "*.tf"
Rules:
  - id: EC2_INSTANCE_TYPE
    message: Instance type should be t2.micro or m3.medium
    resource: aws_instance
    assertions:
      - key: instance_type
        op: in
        value: t2.micro,m3.medium
    severity: WARNING
```

This could also be done by using the or operation with two different assertions:

```
Version: 1
Description: Example rules
Type: Terraform
Files:
  - "*.tf"
Rules:
  - id: EC2_INSTANCE_TYPE
    message: Instance type should be t2.micro or m3.medium
    resource: aws_instance
    assertions:
      or:
        - key: instance_type
          op: eq
          value: t2.micro
        - key: instance_type
          op: eq
          value: m3.medium
    severity: WARNING
```

And this could also be done by looking up the valid values in an S3 object (HTTP endpoints are also supported)

```
Version: 1
Description: Example rules
Type: Terraform
Files:
  - "*.tf"
Rules:
  - id: EC2_INSTANCE_TYPE
    message: Instance type should be t2.micro or m3.medium
    resource: aws_instance
    assertions:
      - key: instance_type
        op: eq
        value_from: s3://your-bucket/instance-types.txt
    severity: FAILURE
```

The assertions and operations were inspired by those in Cloud Custodian: http://capitalone.github.io/cloud-custodian/docs/


## Valid Operations

[Documented Here](docs/operations.md)

# Output

The program outputs a JSON string with the results. The JSON object has the following attributes:

* FilesScanned - a list of the filenames evaluated
* Violations - an object whose keys are the severity of any violations detected. The value for each key is an array with an entry for every violation of that severity.

## Using --query to limit the output

You can limit the output by specifying a JMESPath expression for the --query command line option. For example, if you just wanted to see the ResourceId attribute for failed checks, you can do the following:

```
./config-lint --rules example-files/rules/terraform.yml --query 'Violations.FAILURE[].ResourceId' example-files/config/*
```

# Exit Code

If at least one rule with a severity of FAILURE was triggered the exit code will be 1, otherwise it will be 0.

# Developing new rules using --search

Each rule requires a JMESPath key that it will use to search resources. Documentation for JMESPATH is here: http://jmespath.org/

The expressions can be tricky to get right, so this tool provides a --search option which takes a JMESPath expression. The expression is evaluated against all the resources in the files provided on the command line. The results are written to stdout.

This example will scan the example terraform file and print the "ami" attribute for each resource:

```
./config-lint --rules example-files/rules/terraform.yml --search 'ami' example-files/config/terraform.tf
```

If you specify --search, the rules files is only used to determine the type of configuration files.
The files will *not* be scanned for violations.


# Support for AWS Config Custom Rules

It is also possible to use a rules files in a Lambda that handles events from AWS Config.

[Documented Here](docs/lambda.md)

# Releasing
To release a new version, run `make bumpversion` to increment the patch version and push a tag to GitHub to start the release process.

# TODO

* Add an optional YAML file for project settings, such as ignoring certain rules for certain resources
* Figure out how dependency management works in go
* The lambda function does not handle OverSizedChangeNotification
* The lambda function name is hard-coded in the Makefile
* Region is hard-coded to us-east-1 for GetValueFromS3
* Use type switch as more idiomatic way to handle multiple types in match.go
* Start using go testing coverage tools
* Use log package for error reporting
* Deal with a few FIXME comments in code, mostly error handling
* Should there be some pre-defined RuleSets?
* Would it be useful to have helper utilities to send output to CloudWatch/SNS/Kinesis?
* Add variable interpolation for Terraform files
* Update value_from to handle JSON return values
* Create a Provider interface for AWS calls, create a mock for testing SecurityGroupLinter
* Starting to have inconsistent naming in ops: is-true, is-false, has-properties vs. present, absent, empty, null
* Add options to Assertion type, for things like 'ignore-case' for string compares? Or just use a regex?
* Provide a default -query of 'Violations[]', and add an option for a full report
