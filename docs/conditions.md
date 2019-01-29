# Conditions

Sometimes you want to apply a rule only to a subset of a resource type, based on some attribute in that resource. You can add a list of expressions using the condition attribute. These must all evaluate to true for the rule to be applied. The format for these expressions is the same as that used for the assertions list.

If no condition is specified then the rule will be applied to all matching resources.

It's possible to accomplish the same end thing by using additional expressions directly in the assertions for a rule. But logically it's a little bit different to decide whether or not to apply a rule and what that rule should check. It also reads a little better.

### Example

```
...
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
...
```

