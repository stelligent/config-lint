# Pass
resource "aws_efs_file_system" "encrypted_set_to_true" {
  creation_token = "foo"
  encrypted      = true
}

# Fail
resource "aws_efs_file_system" "encrypted_set_to_false" {
  creation_token = "foo"
  encrypted      = false
}

# Fail
resource "aws_efs_file_system" "encrypted_not_set" {
  creation_token = "foo"
}
