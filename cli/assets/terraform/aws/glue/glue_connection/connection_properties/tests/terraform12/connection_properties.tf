# Test that connection_properties is not providing a plaintext password
# https://www.terraform.io/docs/providers/aws/r/glue_connection.html#connection_properties

provider "aws" {
  region = "us-east-1"
}

# PASS: connection_properties not used
resource "aws_glue_connection" "glue_connection_properties_password_not_used" {
  connection_properties = {
    JDBC_CONNECTION_URL = "jdbc:mysql://example.com/exampledatabase"
    USERNAME            = "exampleusername"
  }
  name = "example"
}

# FAIL: connection_properties are being used
resource "aws_glue_connection" "glue_connection_properties_password_used" {
  connection_properties = {
    JDBC_CONNECTION_URL = "jdbc:mysql://example.com/exampledatabase"
    PASSWORD            = "examplepassword"
    USERNAME            = "exampleusername"
  }
  name = "example"
}
