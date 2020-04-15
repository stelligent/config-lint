# Developing rules for config-lint

## Developing new rules using -search

Each rule requires a JMESPath key that it will use to search resources. Documentation for JMESPATH is here: http://jmespath.org/

The expressions can be tricky to get right, so this tool provides a -search option which takes a JMESPath expression. The expression is evaluated against all the resources in the files provided on the command line. The results are written to stdout.

This example will scan the example terraform file and print the "ami" attribute for each resource:

```
./config-lint -rules example-files/rules/terraform.yml -search 'ami' example-files/config/terraform.tf
```

If you specify -search, the rules files is only used to determine the type of configuration files.
The files will *not* be scanned for violations.