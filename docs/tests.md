# Tests

You can run the project tests by invoking the correct make command:

* `make test` -> Runs all tests residing inside the config-lint project
* `make testtf` -> Runs all tests defined within the `TestTerraformBuiltInRules` test function.

## Testing Best Practices

It is best practice to always come up with **at least 2** scenarios for each test (ideally more if applicable). You want a test case that will *pass* and a test case that will *fail*. This covers the bare minimum to ensure that a rule and test case are working as expected.

## Terraform

Terraform rules have their own set of tests that you can use to verify that a new rule or configuration is working as expected. As noted above, the `make testtf` command will run the tests defined within the `TestTerraformBuiltInRules` test function.

### Creating Terraform Tests

To create a new test to validate a Terraform built in rule you need to do the following:
* Add test case inside `TestTerraformBuiltInRules` function in the `cli/builtin_terraform_test.go` file.
  * Example: `{"security-groups.tf", "SG_WORLD_INGRESS", 1, 0},` will run the `SG_WORLD_INGRESS` rule against the contents of the `cli/testdata/builtin/terraform/security-groups.tf` file. It is expecting there to be **1** *Warning* and **0** *Failures*
  * The test must follow the `struct` format for `BuiltInTestCase`

  ``` go
  type BuiltInTestCase struct {
	Filename     string
	RuleID       string
	WarningCount int
	FailureCount int
  }
  ```

  * This test case must match a valid Rule `id` from within the `cli/assets/terraform.yml` file.
