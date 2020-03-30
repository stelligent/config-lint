# Pass
resource "aws_ami" "ebs_block_device_encrypted_set_to_true" {
  name = "foo"

  ebs_block_device {
    device_name = "/dev/xvda"
    volume_size = 8
    encrypted   = true
  }
}

# Fail
resource "aws_ami" "ebs_block_device_encrypted_set_to_false" {
  name = "foo"

  ebs_block_device {
    device_name = "/dev/xvda"
    volume_size = 8
    encrypted   = false
  }
}

# Fail
resource "aws_ami" "ebs_block_device_encrypted_not_set" {
  name = "foo"

  ebs_block_device {
    device_name = "/dev/xvda"
    volume_size = 8
  }
}
