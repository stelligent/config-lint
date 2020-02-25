# Pass
resource "aws_kms_key" "enable_key_rotation_set_to_true" {
  enable_key_rotation = true
}

# Warn
resource "aws_kms_key" "enable_key_rotation_set_to_false" {
  enable_key_rotation = false
}

# Warn
resource "aws_kms_key" "enable_key_rotation_not_set" {
}
