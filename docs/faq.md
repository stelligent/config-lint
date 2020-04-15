# FAQs

1. With the newest version of config-lint being able to handle configuration files written in Terraform v0.12 syntax, is it still
backwards compatible with configuration files written in the Terraform v0.11?
    - Yes the new version of config-lint is able to handle parsing Terraform configuration files written in both v0.11 and v0.12 syntax.
    - To choose the between parsing Terraform 0.11 vs 0.12 syntax, you can pass in the flag option `-tfparser` followed
    by either `tf11` or `tf12`. For example:
        - `config-lint -tfparser tf12 -rules example_rule.yml example_config/example_file.tf`
2. I'm running into errors when trying to run the newest version of config-lint against configuration files
written in Terraform v0.12 syntax. Where should be the first place to check for resolving this?
    - The first thing to check is to make sure you're passing in the correct `-tfparser` flag option.
    Depending on which Terraform syntax the configuration file is written in, refer to the FAQ #1 above for 
    passing in the correct flag option values.
    - For configuration files that contain Terraform v0.12 syntax, you should confirm that whatever rule.yml file/files you pass in
    have the `type:` key set to `Terraform12`. For example in this rule.yml file:
    ```
   version: 1
   description: Rules for Terraform configuration files
   type: Terraform12
   files:
     - "*.tf"
   rules:
     - id: AMI_SET
       message: Testing
       resource: aws_instance
       assertions:
         - key: ami
           op: eq
           value: ami-f2d3638a
    ```
