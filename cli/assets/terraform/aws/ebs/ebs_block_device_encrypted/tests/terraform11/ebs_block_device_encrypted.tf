## Setup Helper
data "aws_ami" "test_ami" {
  most_recent = true

  filter {
    name   = "name"
    values = ["ubuntu/images/hvm-ssd/ubuntu-bionic-18.04-amd64-server-*"]
  }

  filter {
    name   = "virtualization-type"
    values = ["hvm"]
  }

  owners = ["099720109477"] # Canonical
}

# Pass
resource "aws_instance" "ebs_block_device_not_set" {
  ami           = "${data.aws_ami.ubuntu.id}"
  instance_type = "t2.micro"
}

# Pass
resource "aws_instance" "ebs_block_device_encrypted_set_to_true" {
  ami           = "${data.aws_ami.ubuntu.id}"
  instance_type = "t2.micro"

  ebs_block_device {
    device_name = "/dev/xvda"
    volume_size = 20
    encrypted   = true
  }
}

# Fail
resource "aws_instance" "ebs_block_device_encrypted_set_to_false" {
  ami           = "${data.aws_ami.ubuntu.id}"
  instance_type = "t2.micro"

  ebs_block_device {
    device_name = "/dev/xvda"
    volume_size = 20
    encrypted   = false
  }
}

# Fail
resource "aws_instance" "ebs_block_device_encrypted_not_set" {
  ami           = "${data.aws_ami.ubuntu.id}"
  instance_type = "t2.micro"

  ebs_block_device {
    device_name = "/dev/xvda"
    volume_size = 20
  }
}
