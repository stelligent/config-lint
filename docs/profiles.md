# Profiles

The -profile command line option takes a filename which contains a set of default values for various command line options.
If there is a file in the working directory called `config-lint.yml`, it will be loaded automatically.
All values in the profile are optional, and are overriden by anything specified on the command line.
An example profile:

```
# A list of files containing rules for linting
rules:
  - example-files/rules/generic-yaml.yml

# A list of files to scan
files:
  - example-files/config/*.config

# An optional list of rules to check, the default is all rules
ids:
  - RULE_1
  - RULE_2

# An optional list of tags used to select what rules to apply, the default is all rules
tags:
  - s3

# A list of resources and rules that should not be applied
# This is useful if you want to turn off some rules for some resources, especially
# when using built-in rules
# (For custom rules files, you can use the Except attribute on a rule)
exceptions:
  - RuleID: S3_BUCKET_ACL
    ResourceCategory: resource
    ResourceType: aws_s3_bucket
    ResouceID: simple_website
    Comments: This bucket hosts a public website
```

