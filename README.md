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
./config-lint --terraform files/*
```

## Validate Kubernetes files

```
./config-lint --kubernetes files/*
```


# Rules

The rules file includes a list of objects with the following attributes:

* id: unique identifier for the rule
* message: string to be printed when a validation error is detected
* resource: the resource type to which the rule will be applied
* severity: whether the validation generates a WARNING or a FAILURE
* filters: a list of filters used to detect validation errors
* tags: optional list of tags, command line has option to limit scans to a subset of tags

# Filters

Each filter contains the following attributes:

* type: everything should be "value" for now
* key: the JMES path used to find data in a resource
* op: the operation to be performed on the data returned by searching for the JMES path
* value: needed for most operations

For example, to test that an AWS instance type has one of two values:
```
Version: 1
Description: Example rules
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

# Developing new rules using --search

Each rule requires a JMESPath key that it will use to search resources. Documentation for JMESPATH is here: http://jmespath.org/

The expressions can be tricky to get right, so this tool provides a --search option which takes a JMESPath expression. The expression is evaluated against all the resources in the files provided on the command line. The results are written to stdout.

This example will scan the example terraform file and print the "ami" attribute for each resource:
```
./config-lint --terraform --search 'ami' files/terraform.tf
```

If you specify --search, the rules files are ignored and the files are *not* scanned for violations.

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
