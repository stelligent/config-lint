resource "aws_instance" "foo" {
  ebs_block_device {
    encrypted   = false
  }

  ebs_block_device {
    encrypted   = true
  }

  credit_specification {
    cpu_credits = "unlimited"
  }
}
