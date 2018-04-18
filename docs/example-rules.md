# Example Rules 

Add these rules to a YAML file, and pass the filename to config-lint using the -rules option.
Each rule contains a list of assertions, and these assertions use operations that are [documented here](operations.md).


* [Simple Expressions](#simple-expressions)
* [Boolean Expressions](#boolean-expressions)
* [Collection Expressions](#collection-expressions)
* [Dynamic Values](#dynamic-values)
* [Conditions](#conditions)
* [Macros](#macros)


## simple-expressions

To test that an AWS instance type has one of two values:

```
version: 1
description: Simple expression example
type: Terraform
files:
  - "*.tf"
rules:
  - id: EC2_INSTANCE_TYPE
    message: Instance type should be t2.micro or m3.medium
    resource: aws_instance
    assertions:
      - key: instance_type
        op: in
        value: t2.micro,m3.medium
    severity: WARNING
```

## boolean-expressions

This could also be done by using the or operation with two different assertions:

```
version: 1
description: Boolean expression example
type: Terraform
files:
  - "*.tf"
Rules:
  - id: EC2_INSTANCE_TYPE
    message: Instance type shouldG be t2.micro or m3.medium
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

## collection-expressions

There are three operators that simplify working with collections: [every](operations.md#every), [some](operations.md#some) and [none](operations.md#none).
You provide a JMESPath expression to extract the entire collection from the resource properties.
Then a separate set of expressions are applied to each element of the collection.
The expressions are the same as used for rule assertions.
The every operator requires all elements to return true for the expression, the some operator requires at least one element to return true, and the none operator requires all of the elements to return false for the expression.


```
version: 1
description: Collection expression example
type: YAML
files:
  - "*.config"

resources:
  - type: customer
    key: customers[]
    id: id

rules:

  - id: CUSTOMER_LOCATIONS
    message: Every customer location needs an address and a zip_code
    resource: customer
    severity: FAILURE
    assertions:
      - every:
          key: locations
          assertions:
            - key: address
              op: present
            - key: zip_code
              op: present

```

## dynamic-values

Instead of including a list of values directly in the rules file, it can be retrieved
from an S3 object at runtime. HTTP endpoints are also supported.

```
version: 1
description: Dynamic value example
type: Terraform
files:
  - "*.tf"
rules:
  - id: EC2_INSTANCE_TYPE
    message: Instance type should be t2.micro or m3.medium
    resource: aws_instance
    assertions:
      - key: instance_type
        op: in
        value_from: s3://your-bucket/instance-types.txt
    severity: FAILURE
```

## conditions

Rules always have a condition based on the resource type, but you can add additional conditions. Here is an example
from the internal rule set used for the -validate option. This rule will only apply when a resource of type LintRuleSet
has a type equal to YAML. For that type of LintRuleSet, another attribute called resources must be present:

```
version: 1
description: Condition example
type: YAML
files:
  - *.config
rules:
  - id: YAML_RULES_HAVE_RESOURCES_SECTION
    message: RuleSet for YAML required resources section
    resource: LintRuleSet
    severity: FAILURE
    conditions:
      - key: type
        op: eq
        value: YAML
    assertions:
      - key: resources
        op: present
```

## macros

Because the rules are specified in YAML format, it is possible to use anchors and aliases as a simple kind of macro language to
eliminate duplicate expressions in a rule set. If you are familiar with Ruby on Rails configuration files, this might be familiar.

Here is an example where the rules for different resource types check for the same attribute names.
You can use the "&" to define an anchor. In this example there
are three of these: &has_name, &has_description and &has_name_and_description (the first two are actually used to
define the third one). Elsewhere in the file you can use the "*" (alias) and the "<<" (merge key) to insert
these expressions into multiple rules, without copying the entire expression map.

```
version: 1
description: Macro example
type: YAML
files:
  - "*.config"

# some anchors to keep rules DRY

has_name: &has_name
  key: name
  op: present

has_name: &has_description
  key: description
  op: present

has_name_and_description: &has_name_and_description
  and:
    - <<: *has_name
    - <<: *has_description

# resources to find in the config file

resources:
  - type: widget
    key: widgets[]
    id: id
  - type: gadget
    key: gadgets[]
    id: name

# rules to apply

rules:

  - id: WIDGET_PROPERTIES
    message: Widget needs name and description
    severity: FAILURE
    resource: widget
    assertions:
      - <<: *has_name_and_description

  - id: GADGET_PROPERTIES
    message: Gadget needs name a description
    severity: FAILURE
    resource: gadget
    assertions:
      - <<: *has_name_and_description
```

This feature of YAML can even be used for partial expressions. The JMESPath expression for find a specific tag in a Terraform resource is not particularly friendly, so that could be hidden in an anchor called "has_tag". Then a rule assertion can reference that, and include the tag name that is required.

```
Version: 1
Description: Use an alias to hide a complicated JMESPath expression
Type: Terraform
Files:
  - "*.tf"

has_tag: &has_tag
  key: "tags[]|[0].keys(@)"
  op: contains

rules:

  - id: HAS_NAME_TAG
    message: Tags are required
    resource: aws_ebs_volume
    assertions:
      - <<: *has_tag
        value: Name
```

The assertions and operations were inspired by those in Cloud Custodian: http://capitalone.github.io/cloud-custodian/docs/
