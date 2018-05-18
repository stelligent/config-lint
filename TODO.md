# TODO

* Region is hard-coded to us-east-1 for GetValueFromS3
* Use type switch as more idiomatic way to handle multiple types in match.go
* Use log package for error reporting
* Deal with a few FIXME comments in code, mostly error handling
* Would it be useful to have helper utilities to send output to CloudWatch/SNS/Kinesis?
* Update value_from to handle JSON return values
* Create a Provider interface for AWS calls, create a mock for testing SecurityGroupLinter
* Starting to have inconsistent naming in ops: is-true, is-false, has-properties vs. present, absent, empty, null
* Add options to Assertion type, for things like 'ignore-case' for string compares? Or just use a regex?
* Terraform converter wraps every map in an array - apparently it is valid HCL to have, e.g. "tags" appear multiple times in a resource
