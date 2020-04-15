# Output from config-lint

The program outputs a JSON string with the results. The JSON object has the following attributes:

* FilesScanned - a list of the filenames evaluated
* Violations - an object whose keys are the severity of any violations detected. The value for each key is an array with an entry for every violation of that severity.

## Using -query to limit the output

You can limit the output by specifying a JMESPath expression for the -query command line option. For example, if you just wanted to see the ResourceId attribute for failed checks, you can do the following:

```
./config-lint -rules example-files/rules/terraform.yml -query 'Violations.FAILURE[].ResourceId' example-files/config/*
```