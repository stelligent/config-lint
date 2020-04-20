# Test that user activity logging is enabled
# https://www.terraform.io/docs/providers/aws/r/redshift_parameter_group.html

# FAIL: enable_user_activity_logging is not set
resource "aws_redshift_parameter_group" "logging_not_set" {
  name   = "foobar"
  family = "redshift-1.0"

  parameter {
    name  = "require_ssl"
    value = "true"
  }
}

# FAIL: enable_user_activity_logging is false
resource "aws_redshift_parameter_group" "logging_set_to_false" {
  name   = "foobar"
  family = "redshift-1.0"

  parameter {
    name  = "require_ssl"
    value = "false"
  }

  parameter {
    name  = "enable_user_activity_logging"
    value = "false"
  }

  parameter {
    name  = "query_group"
    value = "example"
  }
}

# PASS: enable_user_activity_logging is set to true
resource "aws_redshift_parameter_group" "logging_set_to_true" {
  name   = "foobar"
  family = "redshift-1.0"

  parameter {
    name  = "require_ssl"
    value = "true"
  }

  parameter {
    name  = "enable_user_activity_logging"
    value = "true"
  }

  parameter {
    name  = "query_group"
    value = "example"
  }
}
