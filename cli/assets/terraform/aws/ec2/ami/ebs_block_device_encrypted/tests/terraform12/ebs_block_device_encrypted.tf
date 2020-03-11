## Setup Helper
variable "test_device" {
  default = "/dev/xvda"
}

variable "test_volume" {
  default = 8
}

# Pass
resource "aws_ami" "ebs_block_device_encrypted_set_to_true" {
  name = "foo"

  ebs_block_device {
    device_name = var.test_device
    volume_size = var.test_volume
    encrypted   = true
  }
}

# Fail
resource "aws_ami" "ebs_block_device_encrypted_set_to_false" {
  name = "foo"

  ebs_block_device {
    device_name = var.test_device
    volume_size = var.test_volume
    encrypted   = false
  }
}

# Fail
resource "aws_ami" "ebs_block_device_encrypted_not_set" {
  name = "foo"

  ebs_block_device {
    device_name = var.test_device
    volume_size = var.test_volume
  }
}
