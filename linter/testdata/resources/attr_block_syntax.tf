# PASS with best syntax
resource "aws_redshift_parameter_group" "pass_best" {
  name   = "foo"
  family = "redshift-1.0"

  parameter {
    name  = "tls"
    value = "enabled"
  }
  parameter  {
    name = "audit_logs"
    value = "enabled"
  }
}

# PASS with Alternate syntax
resource "aws_redshift_parameter_group" "pass_alt" {
  name   = "bar"
  family = "redshift-1.0"

  parameter =
  [{
      name  = "tls"
      value = "enabled"
    },
    {
      name = "audit_logs"
      value = "enabled"
  }]
}

# FAIL with best syntax
resource "aws_redshift_parameter_group" "fail_best" {
  name   = "fail_foo"
  family = "redshift-1.0"

  parameter {
    name  = "tls"
    value = "disabled"
  }
  parameter  {
    name = "audit_logs"
    value = "enabled"
  }
}

# FAIL with Alternate syntax
resource "aws_redshift_parameter_group" "fail_alt" {
  name   = "fail_bar"
  family = "redshift-1.0"

  parameter =
  [{
      name  = "tls"
      value = "disabled"
    },
    {
      name = "audit_logs"
      value = "enabled"
  }]
}
