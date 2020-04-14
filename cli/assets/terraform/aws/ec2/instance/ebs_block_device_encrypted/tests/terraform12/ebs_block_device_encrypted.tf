# Test that EBS block device is using encrpytion and specifies a KMS key
# https://www.terraform.io/docs/providers/aws/r/instance.html#encrypted
# https://www.terraform.io/docs/providers/aws/r/instance.html#kms_key_id

provider "aws" {
  region = "us-east-1"
}

## Setup Helper
data "aws_ami" "ubuntu" {
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

# PASS: Not specifiying an EBS block device
resource "aws_instance" "ebs_block_device_not_set" {
  ami           = data.aws_ami.ubuntu.id
  instance_type = "t2.micro"
}

# PASS: Block device specified with encryption enabled and KMS key
resource "aws_instance" "ebs_block_device_encrypted_set_to_true" {
  ami           = data.aws_ami.ubuntu.id
  instance_type = "t2.micro"

  ebs_block_device {
    device_name = "/dev/xvda"
    volume_size = 20
    encrypted   = true
    kms_key_id  = "arn:aws:kms:us-west-2:111122223333:key/1234abcd-12ab-34cd-56ef-1234567890a"
  }
}

# FAIL: Encryption disabled
# WARN: KMS key not specified
resource "aws_instance" "ebs_block_device_encrypted_set_to_false" {
  ami           = data.aws_ami.ubuntu.id
  instance_type = "t2.micro"

  ebs_block_device {
    device_name = "/dev/xvda"
    volume_size = 20
    encrypted   = false
  }
}

# FAIL: Encryption not specified
# WARN: KMS key not specified
resource "aws_instance" "ebs_block_device_encrypted_not_set" {
  ami           = data.aws_ami.ubuntu.id
  instance_type = "t2.micro"

  ebs_block_device {
    device_name = "/dev/xvda"
    volume_size = 20
  }
}
