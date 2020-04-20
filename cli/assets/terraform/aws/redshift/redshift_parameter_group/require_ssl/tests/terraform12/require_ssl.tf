# Test that require_ssl parameter is present and set to true
# https://www.terraform.io/docs/providers/aws/r/redshift_parameter_group.html

# WARN require_ssl is not set
resource "aws_redshift_parameter_group" "parameter_and_require_ssl_not_set" {
  name   = "foobar"
  family = "redshift-1.0"
}

# WARN: require_ssl is false
resource "aws_redshift_parameter_group" "require_ssl_set_to_false" {
  name   = "foobar"
  family = "redshift-1.0"

  parameter {
    name  = "enable_user_activity_logging"
    value = "true"
  }

  parameter {
    name  = "require_ssl"
    value = "false"
  }

  parameter {
    name  = "query_group"
    value = "example"
  }
}

# PASS: require_ssl is set to true
resource "aws_redshift_parameter_group" "require_ssl_set_to_true" {
  name   = "foobar"
  family = "redshift-1.0"

  parameter {
    name  = "enable_user_activity_logging"
    value = "true"
  }

  parameter {
    name  = "require_ssl"
    value = "true"
  }

  parameter {
    name  = "query_group"
    value = "example"
  }
}
