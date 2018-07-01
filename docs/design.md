# Design

## Motivation

* Static analysis of Terraform configuration files, similar to [cfn_nag](https://github.com/stelligent/cfn_nag) for CloudFormation templates
* Analysis of Kubernetes spec files
* Scanning of AWS resources already provisioned, via AWS API describe* calls
* Processing of event data in AWS Config custom rules

While the data source in each of these cases is different (files vs API calls vs event parameters), the application of the rules follows a similar pattern:

* Extract a data element about a resource
* Make assertions about the value or values found

## Goals

### Rules written in a DSL, rather than code

Using a DSL in YAML has some advantages: 

* Common format that is easy to read
* Easy to add new rules
* Easy to scan an existing rule set

The DSL was modeled after that found in [Cloud Custodian](https://github.com/capitalone/cloud-custodian). Instead of specifying filters to select resources, the DSL makes assertions about values discovered about resources.

The DSL does have the ability to invoke HTTP endpoints. This is intended for logic that is too complex to specify in the DSL.

### Dynamic data

In addition to using data embedded directly in a rule, the rules should be able to reference dynamic data. Typical examples are lists of IP addresses, and EC2 instance types. 

The DSL can reference HTTP endpoints or S3 objects for such dynamic data. This idea was also inspired by Cloud Custodian.

### Not aimed at remediation. 

The primary use case is as part of a CI/CD pipeline. If violations are detected, the exit code of the program is set so the pipeline can be terminated.

The JSON output of the tool can be read by a other tools to trigger notifications or automatic remediation.

## Implementation 

The tool itself is written in go, and is a self contained binary. This simplifies its use in pipelines as well as for local development.  

### Built-in rules

The implementation includes a set of rules for Terraform that can be turned on with a command line option. These implement the same rules found in [cfn_nag](https://github.com/stelligent/cfn_nag) as well as those found in [terrascan](https://github.com/cesar-rodriguez/terrascan)

### Packages

There are three packages in the repository:

* cli
* linter
* assertion

#### cli package

This processes command line arguments and loads project files before using the linter package to do the actual linting.

#### linter package

Defines a linter interface, and provides a factory function to create linters that can discover resources that will be analyzed.
Currently supported:

* File based configurations - Terraform, Kubernetes, YAML, and self validation of config-lint rules files
* API based congifurations - AWS Security Groups and IAM users

The cli package in this repo uses the linter package. The linter package might also be used in tools that are not command line driven, such as a website or an AWS Lambda used for Config rules. Implementations of these can be found in the early commit history of this repository, but future work will be moved to separate repositories.

The linter then uses the assertion package to analyze the collection of resources, and then generate a report.

#### assertion

Applies a set of rules to a collection of JSON objects and returns a report that includes any violations found. 
This package works on JSON objects and has no knowledge of how the resources were loaded. That work is done in the linter package.

