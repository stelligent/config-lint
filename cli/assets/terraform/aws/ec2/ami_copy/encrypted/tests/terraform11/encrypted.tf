# Pass
resource "aws_ami_copy" "encrypted_set_to_true" {
  name              = "foo"
  source_ami_id     = "ami-xxxxxxxx"
  source_ami_region = "us-east-1"
  encrypted         = true
}

# Fail
resource "aws_ami_copy" "encrypted_set_to_false" {
  name              = "foo"
  source_ami_id     = "ami-xxxxxxxx"
  source_ami_region = "us-east-1"
  encrypted         = false
}

# Fail
resource "aws_ami_copy" "encrypted_not_set" {
  name              = "foo"
  source_ami_id     = "ami-xxxxxxxx"
  source_ami_region = "us-east-1"
}
