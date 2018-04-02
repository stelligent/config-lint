# Custom Rule for AWS Config

Build and deploy the lambda function:

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
