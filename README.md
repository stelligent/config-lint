# config-lint

A command line tool to validate configurations using rules specified in a YAML file.
The data being validated can come from template files, such as a Terraform file.
There is also an example of a Linter that runs agains data returned from an AWS API call.

There is also the ability to deploy an AWS Lambda that can be used as a custom rule
AWS Config. The compliance tests are written in YAML, using the same format. This
YAML is stored in an S3 object, and the bucket and key of the object are passed as 
parameters to the AWS Config fule


# Build Command Line tool

```
make config-lint
```

# Run

The program currently supports scanning of the following types of files:

* Terraform
* Kubernetes
* LintRules

## Validate Terraform files

```
./config-lint --rules example-files/rules/terraform.yml example-files/config/*
```

## Validate Kubernetes files

```
./config-lint --rules example-files/rules/kubernetes.yml example-files/config/*
```

## Validate LintRules files

This type of linting allows the tool to lint its own rules.

```
./config-lint --rules example-files/rules/lint-rules.yml example-files/rules/*
```

## Validate Existing Security Groups

```
./config-lint --rules example-files/rules/security-groups.yml
```


# Rules File

The rules file specifies what files to process, and what validations to perform.

## Attributes for the Rule Set

|Name       |Description                                                                         |
|-----------|------------------------------------------------------------------------------------|
|version    |Currently ignored                                                                   |
|description|Text description for the file, not currently used                                   |
|type       |Terraform, Kubernetes, SecurityGroups, AWSConfig                                    |
|files      |Filenames must match one of these patterns to be processed by this set of rules     |
|rules      |A list of rules, see next section                                                   |

## Attributes for each Rule

Each rule contains the following attributes:

|Name       |Description                                                                         |
|-----------|------------------------------------------------------------------------------------|
|id         | A unique identifier for the rule                                                   |
|message    | A string to be printed when a validation error is detected                         |
|resource   | The resource type to which the rule will be applied                                |
|except     | An optional list of resource ids that should not be validated                      |
|severity   | FAILURE, WARNING, NON_COMPLIANT                                                    |
|assertions | A list of assertions used to detect validation errors, see next section            |
|invoke     | Alternative to assertions for a custom external API call to validate, see below    |
|tags       | Optional list of tags, command line has option to limit scans to a subset of tags  |

## Attributes for each Assertion

Each assertion contains the following attributes:

|Name       |Description                                                                         |
|-----------|------------------------------------------------------------------------------------|
|key        | JMES path used to find data in a resource                                          |
|op         | Operation to perform on the data returned by searching for the key                 |
|value      | Literal value needed for most operations                                           |

## Invoke external API for validation

|Name       | Description                                                                        |
|-----------|------------------------------------------------------------------------------------|
|Url        | HTTP endpoint to invoke                                                            |
|Payload    | Optional JMESPATH to use for payload, default is '@'                               |

## Examples

To test that an AWS instance type has one of two values:
```
Version: 1
Description: Example rules
Type: Terraform
Files:
  - "*.tf"
Rules:
  - id: R1
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
  - id: R2
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

The assertions and operations are modeled after those used by Cloud Custodian: http://capitalone.github.io/cloud-custodian/docs/


## [Valid Operations](docs/operations.md)

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


# Custom Rule for AWS Config

```
make lambda
```

This builds and deploys an AWS Lambda function. The ARN for the Lambda is used to set up a custom AWS Config rule. 
The same YAML format is used to specify the rules to test for  compliance. The severity of the rules for
this use case should be set to NON_COMPLIANT

There are two parameters that need to also be configured for the AWS Config rule:

|Name       |Description                                                                         |
|-----------|------------------------------------------------------------------------------------|
|bucket     | S3 bucket that contains the S3 object with the YAML rules                          |
|key        | Key of the S3 object                                                               |


## AWS Config example

Here's an example of an AWS Config rule that checks for port 22 being open to all IP addresses.
It also includes the 'except:' option which allows the check to be ignored for some resources.

```
Version: 1
Description: Rules for AWS Config
Type: AWSConfig
Rules:
  - id: SG1
    message: Security group should not allow ingress from 0.0.0.0/0
    resource: AWS::EC2::SecurityGroup
    except:
      - sg-88206cff
    severity: NON_COMPLIANT
    assertions:
      - not:
          - and:
              - key: ipPermissions[].fromPort[]
                op: contains
                value: "22"
              - key: ipPermissions[].ipRanges[]
                op: contains
                value: 0.0.0.0/0
```

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
* Use the LintRules linter to implement a -validate option
* Should there be some pre-defined RuleSets?
* Would it be useful to have helper utilities to send output to CloudWatch/SNS/Kinesis?
* Add variable interpolation for Terraform files
* Update value_from to handle JSON return values
* Create a Provider interface for AWS calls, create a mock for testing SecurityGroupLinter
* Refactor lambda.go so unit tests can be written
