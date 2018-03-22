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

The program currently supports Terraform and Kubernetes files.

## Validate Terraform files

```
./config-lint --rules example-files/rules/terraform.yml example-files/config/*
```

## Validate Kubernetes files

```
./config-lint --rules example-files/rules/kubernetes.yml example-files/config/*
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
|type       |Should be 'Terraform' or 'Kubernetes'                                               |
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
|severity   | Should be 'WARNING' or 'FAILURE'                                                   |
|assertions | A list of assertions used to detect validation errors, see next section            |
|tags       | Optional list of tags, command line has option to limit scans to a subset of tags  |

## Attributes for each Assertion

Each assertion contains the following attributes:

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
    assertions:
      - type: value
        key: instance_type
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

The assertions and operations are modeled after those used by Cloud Custodian: http://capitalone.github.io/cloud-custodian/docs/


## Operations supported for an Assertion

* eq - Equals

Example:
```
...
  - id: VOLUME1
    resource: aws_ebs_volume
    message: EBS Volumes must be encrypted
    severity: FAILURE
    assertions:
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
    assertions:
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
    assertions:
      - type: value
        key: instance_type
        op: in
        value: t2.micro,m3.medium
    severity: WARNING
...
```

* not-in - Not in list of values

* present - Attribute is present

Example:
```
...
  - id: R6
    message: Department tag is required
    resource: aws_instance
    assertions:
      - type: value
        key: "tags[].Department | [0]"
        op: present
    severity: FAILURE
...
```

* absent - Attribute is not present

* contains - Attribute contains a substring

* regex - Attribute matches a regular expression

* not - Logical not of another assertions

Example:
```
...
  - id: NOTTEST
    resource: aws_instance
    message: Should not have instance type of c4.large
    severity: WARNING
    assertions:
      - not:
        - type: value
          key: instance_type
          op: eq
          value: c4.large
...
```

* and - Logical and of a list of assertions

Example:
```
...
  - id: ANDTEST
    resource: aws_instance
    message: Should have both Project and Department tags
    severity: WARNING
    assertions:
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

* or - Logical or  of a list of assertions

Example:

```
...
  - id: ORTEST
    resource: aws_instance
    message: Should have instance_type of t2.micro or m3.medium
    severity: WARNING
    assertions:
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

## Invoking an external API for more difficult cases

If the combination of JMESPath and the simple expression DSL are not sufficient, it is possible to have the
rules engine make an API call to validate a resource. Instead of the list of assertions, set the invoke
attribute to an object containg these attributes:

* url - An HTTP GET request will be made to this URL
* payload - A JMESPath expression used to generate the JSON payload to include in the GET request. If not provided, the entire resource will be included (same as using '@' in JMESPath)

The return value should look like this:
```
{
   "Violations": [
       { "Message": "First Violation" }
   ]
}
```

Example:
```
...
  - id: CUSTOM
    severity: FAILURE
    message: Custom
    resource: Policy
    invoke:
      url: https://19kfojjbi2.execute-api.us-east-1.amazonaws.com/dev
      payload: "{ user: spec.user, namespace: spec.namespace }"
    tags:
      - custom
...
```


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
              - type: value
                key: ipPermissions[].fromPort[]
                op: contains
                value: "22"
              - type: value
                key: ipPermissions[].ipRanges[]
                op: contains
                value: 0.0.0.0/0
```

# TODO

Lots to do. This is just a proof-of-concept.

* Add an optional YAML file for project settings, such as ignoring certain rules for certain resources
* Figure out what other assertion types might be needed (if any)
* Finish implementing value_from to allow for dynamic data (again, see Cloud Custodian)
* Figure out how dependency management works in go
* The lambda function does not handle OverSizedChangeNotification
* The lambda function name is hard-coded in the Makefile
