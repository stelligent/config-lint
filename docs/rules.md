# Rules File

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
|category   | Optional value used for Terraform: resource(default), data, provider               |
|conditions | Expressions (in addition to resource) that determine if a rule should apply        |
|except     | An optional list of resource ids that should not be validated                      |
|severity   | FAILURE, WARNING, NON_COMPLIANT                                                    |
|assertions | A list of expressions used to detect validation errors, see next section           |
|invoke     | Alternative to assertions for a custom external API call to validate, see below    |
|tags       | Optional list of tags, command line has option to limit scans to a subset of tags  |
|conditions | Optional expressions that must be true for the rule to be applied                  |

## Attributes for each Expression

Each expression contains the following attributes:

|Name       |Description                                                                         |
|-----------|------------------------------------------------------------------------------------|
|key        | JMES path used to find data in a resource                                          |
|op         | Operation to perform on the data return. [See here for valid operations](operations.md) |
|value      | Literal value needed for most operations                                           |
|value_from | Endpoint for loading values dynamically                                            |

## Invoke external API for validation

|Name       | Description                                                                        |
|-----------|------------------------------------------------------------------------------------|
|Url        | HTTP endpoint to invoke                                                            |
|Payload    | Optional JMESPATH to use for payload, default is '@'                               |

