## Setup Helper
variable "test_ami" {
  default = "ami-xxxxxxxx"
}

variable "test_region" {
  default = "us-east-1"
}

# Pass
resource "aws_ami_copy" "encrypted_set_to_true" {
  name              = "foo"
  source_ami_id     = var.test_ami
  source_ami_region = var.test_region
  encrypted         = true
}

# Pass
resource "aws_ami_copy" "encrypted_set_to_false" {
  name              = "foo"
  source_ami_id     = var.test_ami
  source_ami_region = var.test_region
  encrypted         = false
}

# Pass
resource "aws_ami_copy" "encrypted_not_set" {
  name              = "foo"
  source_ami_id     = var.test_ami
  source_ami_region = var.test_region
}
