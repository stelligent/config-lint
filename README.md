# terraform-lint

Validate a terraform HCL file using rules specified in a YAML file.

# Run

```
go run app.go files/terraform.hcl
```

# Rules

The rules file is currently hard-coded to be 'rules/terraform.yml'. It is a list of objects with the following attributes:

* id: unique identifier for the rule
* message: string to be printed when a validation error is detected
* resource: the resource type to which the rule will be applied
* severity: whether the validation generates a WARNING or a FAILURE
* filters: a list of filters used to detect validation errors
* tags: optional list of tags, command line has option to limit scans to a subset of tags

My thought is to require a command line parameter with a file or directory name where the rules can be found.
If a directory name is given, load all the files in that directory. Maybe allow multiple directories to be specified

# Filters

Each filter contains the following attributes:

* type: everything should be "value" for now
* key: the JMES path used to find data in a resource
* op: the operation to be performed on the data returned by searching for the JMES path
* value: needed to most operations

For example, to test that an AWS instance type has one of two values:
```
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

The filters and operations are modeled after those used by CloudCustodian: http://capitalone.github.io/cloud-custodian/docs/

# TODO

Lots to do. This is just a proof-of-concept.

* Embedded JSON for IAM policies should be parsed and made available for JMESPath query
* Add an optional YAML file for project settings, such as ignoring certain rules for certain resources
* Implement more of the operators from CloudCustodian
* Figure out what other filter types might be needed (if any)
* Improve output - table format, JSON format
* Output should be grouped by resource id
* Add value_from to allow for dynamic data (again, see CloudCustodian)
* Add ability to extend with a Lambda function
* Add command line parameter for rules file or directory (currently hard-coded)
* It's already big enough to warrant some automated tests
* It should be possible to nest the and, or, not operators
