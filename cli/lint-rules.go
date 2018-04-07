package main

// LintRules string containing YAML for -validate option
var LintRules = `---
version: 1
description: Rules for config-lint
type: LintRules
files:
  - "*.yml"
rules:

  - id: VALID_TYPE
    message: Not a valid linter type
    resource: LintRuleSet
    severity: FAILURE
    assertions:
      - key: type
        op: in
        value: Terraform,Kubernetes,SecurityGroup,IAMUser,AWSConfig,LintRules,YAML

  - id: VALID_VERSION
    message: RuleSet must have a supported version
    resource: LintRuleSet
    severity: WARNING
    assertions:
      - key: version
        op: eq
        value: 1

  - id: HAS_RULES
    message: RuleSet needs at least one rule
    resource: LintRuleSet
    severity: WARNING
    assertions:
      - key: rules
        op: not-empty

  - id: EVERY_RULE_HAS_ID
    message: Event rule in rule set must have an id
    resource: LintRuleSet
    severity: FAILURE
    assertions:
      - every:
          key: rules
          assertions:
            - key: id
              op: present

  - id: ID_PRESENT
    message: Rule must have an ID
    resource: LintRule
    severity: FAILURE
    assertions:
      - key: id
        op: present

  - id: SEVERITY_PRESENT
    message: Rule must have a severity
    resource: LintRule
    severity: FAILURE
    assertions:
      - key: severity
        op: present

  - id: RESOURCE_PRESENT
    message: Rule must have a resource filter
    severity: FAILURE
    resource: LintRule
    assertions:
      - key: resource
        op: present
    tags:
      - resource

  - id: ASSERTIONS_OR_INVOKE
    message: Rule must have assertions or invoke
    resource: LintRule
    severity: FAILURE
    assertions:
      - or:
          - key: assertions
            op: present
          - key: invoke
            op: present

  - id: VALID_EXPRESSION
    message: "There are mutually exclusive in the same expression: key,or,xor,and,not,every,some,none"
    resource: LintRule
    severity: FAILURE
    assertions:
      - every:
          key: "assertions[]"
          assertions:
            - xor:
              - key: "@"
                op: has-properties
                value: key,op
              - key: "@"
                op: has-properties
                value: or
              - key: "@"
                op: has-properties
                value: xor
              - key: "@"
                op: has-properties
                value: and
              - key: "@"
                op: has-properties
                value: not
              - key: "@"
                op: has-properties
                value: every
              - key: "@"
                op: has-properties
                value: some
              - key: "@"
                op: has-properties
                value: none
`
