# Pass
resource "aws_ebs_volume" "encrypted_set_to_true" {
  availability_zone = "us-west-2a"
  size              = 20
  encrypted         = true
}

# Fail
resource "aws_ebs_volume" "encrypted_set_to_false" {
  availability_zone = "us-west-2a"
  size              = 20
  encrypted         = false
}

# Fail
resource "aws_ebs_volume" "encrypted_not_set" {
  availability_zone = "us-west-2a"
  size              = 20
}
