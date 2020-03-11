# Warn
resource "aws_redshift_parameter_group" "parameter_and_require_ssl_not_set" {
  name   = "foobar"
  family = "redshift-1.0"
}

# Warn
resource "aws_redshift_parameter_group" "require_ssl_set_to_false" {
  name   = "foobar"
  family = "redshift-1.0"

  parameter {
    name  = "require_ssl"
    value = "false"
  }
}

# Pass
resource "aws_redshift_parameter_group" "require_ssl_set_to_true" {
  name   = "foobar"
  family = "redshift-1.0"

  parameter {
    name  = "require_ssl"
    value = "true"
  }
}
