# Rules File

## Attributes for the Rule Set

|Name       |Description                                                                         |
|-----------|------------------------------------------------------------------------------------|
|version    |Currently ignored                                                                   |
|description|Text description for the file, not currently used                                   |
|type       |Terraform, Terraform12, Kubernetes, SecurityGroups, AWSConfig                                    |
|files      |Filenames must match one of these patterns to be processed by this set of rules     |
|rules      |A list of rules, see next section                                                   |

## Attributes for each Rule

Each rule contains the following attributes:

|Name             |Description                                                                         |
|-----------------|------------------------------------------------------------------------------------|
|id               | A unique identifier for the rule                                                   |
|message          | A string to be printed when a validation error is detected                         |
|resource         | The resource type to which the rule will be applied                                |
|resources        | A list of resources types to which the rule will be applied                        |
|except_resources | A list of resource types to exclude                                                |
|category         | Optional value used for Terraform: resource(default), data, provider               |
|[conditions](conditions.md)       | Expressions (in addition to resource) that determine if a rule should apply        |
|except           | An optional list of resource ids that should not be validated                      |
|severity         | FAILURE, WARNING, NON_COMPLIANT                                                    |
|assertions       | A list of expressions used to detect validation errors, see next section           |
|invoke           | Alternative to assertions for a custom external API call to validate, see below    |
|tags             | Optional list of tags, command line has option to limit scans to a subset of tags  |

## Attributes for each Expression

Each expression contains the following attributes:

|Name       |Description                                                                         |
|-----------|------------------------------------------------------------------------------------|
|key        | JMES path used to find data in a resource                                          |
|[op](operations.md)         | Operation to perform on the data return. [See here for valid operations](operations.md) |
|value      | Literal value needed for most operations                                           |
|[value_from](value_from.md) | Endpoint for loading values dynamically [See here for dynamic values](value_from.md) |

## Invoke external API for validation

|Name       | Description                                                                        |
|-----------|------------------------------------------------------------------------------------|
|Url        | HTTP endpoint to invoke                                                            |
|Payload    | Optional JMESPATH to use for payload, default is '@'                               |

