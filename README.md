# config-lint

Validate configuration files using rules specified in a YAML file.

# Build

```
make
```

# Run

The program currently supports Terraform and Kubernetes files.

## Validate Terraform files

```
./config-lint --rules rules/terraform.yml files/*
```

## Validate Kubernetes files

```
./config-lint --rules rules/kubernetes.yml files/*
```


# Rules File

The rules file specifies what files to process, and what validations to perform.

## Attributes for the Rule Set

|Name       |Description                                                                         |
|-----------|------------------------------------------------------------------------------------|
|Version    |Currently ignored                                                                   |
|Description|Text description for the file, not currently used                                   |
|Type       |Should be 'Terraform' or 'Kubernetes'                                               |
|Files      |Filenames must match one of these patterns to be processed by this set of rules     |
|Rules      |A list of rules, see next section                                                   |

## Attributes for each Rule

Each rule contains the following attributes:

|Name       |Description                                                                         |
|-----------|------------------------------------------------------------------------------------|
|id         | A unique identifier for the rule                                                   |
|message    | A string to be printed when a validation error is detected                         |
|resource   | The resource type to which the rule will be applied                                |
|severity   | Should be 'WARNING' or 'FAILURE'                                                   |
|filters    | A list of filters used to detect validation errors, see next section               |
|tags       | Optional list of tags, command line has option to limit scans to a subset of tags  |

## Attributes for each Filter

Each filter contains the following attributes:

|Name       |Description                                                                         |
|-----------|------------------------------------------------------------------------------------|
|type       | Should always be "value" for now                                                   |
|key        | JMES path used to find data in a resource                                          |
|op         | Operation to perform on the data returned by searching for the key                 |
|value      | Literal value needed for most operations                                           |

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
    filters:
      - type: value
        key: instance_type
        op: in
        value: t2.micro,m3.medium
    severity: WARNING
```

This could also be done by using the or operation with two different filters:

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
    filters:
      or:
        - type: value
          key: instance_type
          op: eq
          value: t2.micro
        - type: value
          key: instance_type
          op: eq
          value: m3.medium
    severity: WARNING
```

The filters and operations are modeled after those used by CloudCustodian: http://capitalone.github.io/cloud-custodian/docs/


## Operations supported for a Filter

* eq - Equals

Example:
```
...
  - id: VOLUME1
    resource: aws_ebs_volume
    message: EBS Volumes must be encrypted
    severity: FAILURE
    filters:
      - type: value
        key: encrypted
        op: eq
        value: true
...
```

* ne - Not equals

Example:
```
...
  - id: SG1
    resource: aws_security_group
    message: Security group should not allow ingress from 0.0.0.0/0
    severity: FAILURE
    filters:
      - type: value
        key: "ingress[].cidr_blocks[] | [0]"
        op: ne
        value: "0.0.0.0/0"
...
```

* in - In list of values

Example:
```
...
  - id: R1
    message: Instance type should be t2.micro or m3.medium
    resource: aws_instance
    filters:
      - type: value
        key: instance_type
        op: in
        value: t2.micro,m3.medium
    severity: WARNING
...
```

* notin - Not in list of values

* present - Attribute is present

Example:
```
...
  - id: R6
    message: Department tag is required
    resource: aws_instance
    filters:
      - type: value
        key: "tags[].Department | [0]"
        op: present
    severity: FAILURE
...
```

* absent - Attribute is not present

* contains - Attribute contains a substring

* regex - Attribute matches a regular expression

* not - Logical not of another filter

Example:
```
...
  - id: NOTTEST
    resource: aws_instance
    message: Should not have instance type of c4.large
    severity: WARNING
    filters:
      - not:
        - type: value
          key: instance_type
          op: eq
          value: c4.large
...
```

* and - Logical and of a list of filters

Exmaple:
```
...
  - id: ANDTEST
    resource: aws_instance
    message: Should have both Project and Department tags
    severity: WARNING
    filters:
      - and:
        - type: value
          key: "tags[].Department | [0]"
          op: present
        - type: value
          key: "tags[].Project | [0]"
          op: present
    tags:
      - and-test
...
```

* or - Logical or  of a list of filters

Example:

```
...
  - id: ORTEST
    resource: aws_instance
    message: Should have instance_type of t2.micro or m3.medium
    severity: WARNING
    filters:
      - or:
        - type: value
          key: instance_type
          op: eq
          value: t2.micro
        - type: value
          key: instance_type
          op: eq
          value: m3.medium
...
```

# Output

The program outputs a JSON string with the results. The JSON object has the following attributes:

* FilesScanned - a list of the filenames evaluated
* Violations - an object whose keys are the severity of any violations detected. The value for each key is an array with an entry for every violation of that severity.

## Using --query to limit the output

You can limit the output by specifying a JMESPath expression for the --query command line option. For example, if you just wanted to see the RecordIds for failed checks, you can do the following:

```
./config-lint --rules rules/terraform.yml --query 'Violations.FAILURE[].ResourceId' files/*
```

# Exit Code

If at least one rule with a severity of FAILURE was triggered the exit code will be 1, otherwise it will be 0.

# Developing new rules using --search

Each rule requires a JMESPath key that it will use to search resources. Documentation for JMESPATH is here: http://jmespath.org/

The expressions can be tricky to get right, so this tool provides a --search option which takes a JMESPath expression. The expression is evaluated against all the resources in the files provided on the command line. The results are written to stdout.

This example will scan the example terraform file and print the "ami" attribute for each resource:

```
./config-lint --rules rules/terraform.yml --search 'ami' files/terraform.tf
```

If you specify --search, the rules files is only used to determine the type of configuration files.
The files will *not* be scanned for violations.


# TODO

Lots to do. This is just a proof-of-concept.

* Embedded JSON for IAM policies should be parsed and made available for JMESPath query
* Add an optional YAML file for project settings, such as ignoring certain rules for certain resources
* Implement more of the operators from CloudCustodian
* Figure out what other filter types might be needed (if any)
* Output should be grouped by resource id
* Add value_from to allow for dynamic data (again, see CloudCustodian)
* Add ability to extend with a Lambda function
* It should be possible to nest the and, or, not operators
* Instead of iterating through rules, filters, then resources, make resources the outer loop, so results are reports by resource id
* Add examples to this file for Kubernetes files
* Support multiple rules files, or a rules directory
* Add a --table option which uses tablewriter for more readable report
* Not operator takes a list (to match the way Cloud Custodian works), should make sure size == 1
