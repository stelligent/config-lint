# Example Rules 

Add these rules to a YAML file, and pass the filename to config-lint using the -rules option.
Each rule contains a list of assertions, and these assertions use operations that are [documented here](operations.md)

To test that an AWS instance type has one of two values:

```
Version: 1
Description: Example rules
Type: Terraform
Files:
  - "*.tf"
Rules:
  - id: EC2_INSTANCE_TYPE
    message: Instance type should be t2.micro or m3.medium
    resource: aws_instance
    assertions:
      - key: instance_type
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
  - id: EC2_INSTANCE_TYPE
    message: Instance type should be t2.micro or m3.medium
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

And this could also be done by looking up the valid values in an S3 object (HTTP endpoints are also supported)

```
Version: 1
Description: Example rules
Type: Terraform
Files:
  - "*.tf"
Rules:
  - id: EC2_INSTANCE_TYPE
    message: Instance type should be t2.micro or m3.medium
    resource: aws_instance
    assertions:
      - key: instance_type
        op: eq
        value_from: s3://your-bucket/instance-types.txt
    severity: FAILURE
```

The assertions and operations were inspired by those in Cloud Custodian: http://capitalone.github.io/cloud-custodian/docs/




