# Terraform Linting

## Validate Terraform files with built-in rules

There is a set of [built-in rules](/cli/assets/terraform) that cover some best practices for AWS resources.

```
config-lint -terraform <FILE_OR_DIRECTORY_OF_TF_FILES>
```

If you want to run most of the built-in rules, but not all, you can use a [profile](profiles.md) to exclude some rules or resources.

The Terraform12 parser is fully backwards compatible with previous versions of Terraform. By default, Terraform files will be validated with Terraform 0.12 standards. 

If you wish to force a specific parser version, add the `-tfparser tf11|tf12` flag. This is useful if you have a lot of rules with `Type: Terraform` but your Terraform files include Terraform 12 syntax. 

## Custom Terraform rules for your project or organization

```
config-lint -rules <CUSTOM_RULE_YML_FILE> <FILE_OR_DIRECTORY_OF_TF_FILES>
```

You can specify the -rules option multiple times if you have multiple custom rule files. It is also possible to specify both the -terraform option as well as one or more -rules options, if you want the built-in rules as well as some custom rules.

### Categories

The default category for resources that can be linter is "resource", which covers the most common use case. This is for things like aws_instances, or s3_buckets, etc. But all other block types are available for Terraform linting.

* data
* locals
* module
* output
* provider
* resource
* terraform
* variable

### Rule Structure

Rules are divided into their respective resource directory starting under `assets/terraform`. Each rule is organized following the same tiered directory structure `{ Provider }} / {{ Major Family }} / {{ Resource Name }} / {{ Rule Name }} / rule.yml` where Major Family and Resource Name follow the same naming conventions defined by Terraform. For example, `cli/assets/terraform/aws/batch/batch_job_definition/container_properties_privileged/rule.yml`. The rule configuration itself must be named `rule.yml`.

```
└── terraform
    ├── aws
    │   ├── batch
    │   │   └── batch_job_definition
    │   │       └── container_properties_privileged
    │   │           ├── rule.yml
    │   │           └── tests
    │   │               ├── terraform11
    │   │               │   └── container_properties_privileged.tf
    │   │               ├── terraform12
    │   │               │   └── container_properties_privileged.tf
    │   │               └── test.yml
    ...
```


### Rule Example

```yaml
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

### Top level blocks

Config-lint mainly works on [block arguments and expressions](https://www.terraform.io/docs/configuration/index.html#arguments-blocks-and-expressions) level (the inner part of Terraform blocks), however top level block types and names could linted using `__type__` and `__name__` keys.

Here is an example to follow [Terraform best practices](https://www.terraform.io/docs/extend/best-practices/naming.html) for naming:

```yaml
---
version: 1
description: Make sure Terraform top level blocks follow best practices.
type: Terraform
files:
  - "*.tf"
rules:
  - id: TF_RESOURCE_NAMING_CONVENTION
    message: "Terraform resource block name should match the naming convention. Name should be: not more 64 chars, starts with letter, doesn't have dash, and ends with letter or number"
    severity: FAILURE
    category: resource
    assertions:
    - key: __name__
      op: regex
      value: '^[a-z][a-z0-9_]{0,62}[a-z0-9]$'
  - id: TF_DATA_NAMING_CONVENTION
    message: "Terraform data block name should match the naming convention. Name should be: not more 64 chars, starts with letter, doesn't have dash, and ends with letter or number"
    severity: FAILURE
    category: data
    assertions:
    - key: __name__
      op: regex
      value: '^[a-z][a-z0-9_]{0,62}[a-z0-9]$'
```

Another example, maybe you want to make sure there are no beta providers (e.g. [google-beta](https://www.terraform.io/docs/providers/google/guides/provider_versions.html#google-beta)) used in production:
```yaml
rules:
  - id: TF_PROVIDER_NO_BETA
    message: "No beta feature providers should be used in production"
    severity: FAILURE
    category: data
    assertions:
    - not:
      - key: __type__
        op: regex
        value: '.*?beta.*'
```

**Please note:**

Not all blocks have `__type__` and `__name__` keys. It depends on the block itself. For example, `variable` blocks have name but not type.

There are 4 groups in that regard:
* **Type and name:** data, resource.
* **Type only:** provider.
* **Name only:** module, output, variable.
* **No type and no name:** locals, terraform (they are linted using normal keys defined within them).

### Provider Example

For providers, set the category to "provider" and the resource attribute to the name of the provider.

```yaml
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


```yaml
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

### Terraform 12 Example

Note the `type: Terraform12` item below. Rules targeting templates with Terraform 12-specific features must use the Terraform12 type.

```yaml
version: 1
description: Rules for Terraform configuration files
type: Terraform12
files:
  - "*.tf"
rules:

  - id: CIDR_SET
    message: Testing
    resource: aws_security_group
    assertions:
      - every:
          key: "ingress"
          expressions:
            # this says that it either must be a private IP, or not have IP regex (eg sg string, interpolation)
            - every:
                key: cidr_blocks
                expressions:
                  - key: "@"
                    op: contains
                    value: "/24"
```

### Evaluating Terraform 12 Dynamic Blocks

Dynamic blocks are a new feature introduced in Terraform 12 that enables users to dynamically construct repeatable nested blocks such as ingress rules in an AWS Security Group.

Writing rules for dynamic blocks is a little tricky, as the structure that Terraform parses the .tf file into is different than you may expect.

This Terraform config will generate an `ingress` block for reach item in the `service_ports` list variable.
```hcl-terraform
variable "service_ports" {
  default = [22, 80, 1433, 6379]
}

resource "aws_security_group" "example" {
  name = "example"

  dynamic "ingress" {
    for_each = var.service_ports
    content {
      from_port = ingress.value
      to_port   = ingress.value
      protocol  = "tcp"
    }
  }

  egress = "-1"
}
```

The following rule will result in an error if port 22 (SSH) is included as an ingress for the security group.

The JMESPATH expression refers to keys ("dynamic" and "for_each") that are generated by Terraform, rather than what is present in the configuration.

```yaml
version: 1
description: Rules for Terraform configuration files
type: Terraform12
files:
  - "dynamic_block.tf"
rules:
  - id: NO_SSH_ACCESS
    message: Testing
    resource: aws_security_group
    assertions:
      - key: "dynamic[*].for_each[]"
        op: not-contains
        value: 22
```

### Testing Builtin Rules

All rules need to be tested. All tests for a given rule will be included at the same rule path as the rule configuration itself and live under the `tests` folder. That test folder must include a configuration for the tests, named `test.yml`, and the resources required for testing. The test configuration file must follow the following format:

```yaml
---
version: 1
description: Terraform 11 and 12 tests
type: Terraform
files:
  - "*.tf"
  - "*.tfvars"
tests:
  -
    ruleId: RULE_1_TO_BE_TESTED
    warnings: 0
    failures: 2
    tags:
      - "terraform11"
      - "terraform12"
  -
    ruleId: RULE_2_TO_BE_TESTED
    warnings: 2
    failures: 0
    tags:
      - "terraform11"
  -
    ruleId: RULE_2_TO_BE_TESTED
    warnings: 3
    failures: 0
    tags:
      - "terraform12"
```

The `ruleId` must match the RuleID given in the rule configuration. Warnings or Failures will check against any resource files included under the named `tag` directory within the same tests directory. For example `RULE_1_TO_BE_TESTED` will run the rule against resources in both folders terraform11 and terraform12 and check for the same number of warnings and failures for both. Whereas `RULE_2_TO_BE_TESTED` will check for a different number of warnings for the two versions.
