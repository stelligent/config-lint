# Terraform Linting

## Validate Terraform files with built-in rules

There is a set of [built-in rules](cli/assets/terraform.yml) that cover some best practices for AWS resources.
There

```
config-lint -terraform <FILE_OR_DIRECTORY_OF_TF_FILES>
```

If you want to run most of the built-in rules, but not all, you can use a [profile](docs/profiles.md) to exclude some rules or resources.

## Custom Terraform rules for your project or organization

```
config-lint -rules <CUSTOM_RULE_YML_FILE> <FILE_OR_DIRECTORY_OF_TF_FILES>
```

You can specify the -rules option multiple times if you have multiple custom rule files. It is also possible to specify both the -terraorm option as well as one or more -rules options, if you want the built-in rules as well as some custom rules.

### Categories

The default category for resources that can be linter is "resource", which covers the most common use case. This is for things like aws__instances, or s3_buckets, etc. But there are some additional categories available for Terraform linting. The current list of supported categories is:

* resource
* data
* provider
* module

### Resource Example

```
---
version: 1
description: Check for tags in Terraform file
type: Terraform
files:
  - "*.tf"
rules:
  - id: REQUIRED_TAGS
    message: "A required tag is missing"
    resources:
      - aws_s3_bucket
      - aws_instance
    assertions:
      - key: tags[0]
        op: has-properties
        value: environment,cost_center
  - id: VALID_ENVIRONMENT_TAG
    message: "The environment tag is not valid"
    resources:
      - aws_s3_bucket
      - aws_instance
    assertions:
      - key: tags[0].environment
        op: in
        value: dev,prod,stage
  - id: VALID_COST_CENTER_TAG
    message: "Cost center must be a 4 digit number"
    resources:
      - aws_s3_bucket
      - aws_instance
    assertions:
      - key: tags[0].cost_center
        op: regex
        value: "^[0-9]{4}$"
```

### Provider Example

For providers, set the category to "provider" and the resource attribute to the name of the provider.

```
---
version: 1
description: Terraform provider example
type: Terraform
files:
  - "*.tf"
rules:
  - id: NO_SECRETS_IN_AWS_PROVIDER
    category: provider
    resource: aws
    assertions:
      - key: access_key
        op: absent
      - key: secret_key
        op: absent
```

### Module Example

For modules, use "module" for category, and for resource use the "source" attribute.

This allows checking of parameters being used when a module is referenced.


```
---
version: 1
description: Terraform module invocation example
type: Terraform
files:
  - "*.tf"
rules:
  - id: MODULE_EXAMPLE
    message:
    category: module
    resource: "example/website"
    assertions:
      - key: num_servers
        op: present
```


