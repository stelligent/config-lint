version: 1
description: Rules for Terraform configuration files
type: Terraform
files:
  - "*.tf"
rules:

  - id: INSTANCE_TYPE
    message: Instance type should be t2.micro or m3.medium
    resource: aws_instance
    assertions:
      - key: instance_type
        op: in
        value: t2.micro,m3.medium
    severity: WARNING

  - id: PROJECT_TAG
    message: Check project tag
    resource: aws_instance
    assertions:
      - key: "tags[].project|[0]"
        op: eq
        value: demo
    severity: WARNING
