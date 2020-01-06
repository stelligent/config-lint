# Configuration for Terraform 0.12

variable "vpc_id" {
  description = "ID for the AWS VPC where a security group is to be created."
}

variable "subnet_numbers" {
  description = "List of 8-bit numbers of subnets of base_cidr_block that should be granted access."
  default = [1, 2, 3]
}

data "aws_vpc" "example" {
  id = var.vpc_id
}

resource "aws_security_group" "example" {
  name        = "friendly_subnets"
  description = "Allows access from friendly subnets"
  vpc_id      = var.vpc_id

  ingress {
    from_port = 0
    to_port   = 0
    protocol  = -1

    # For each number in subnet_numbers, extend the CIDR prefix of the
    # requested VPC to produce a subnet CIDR prefix.
    # For the default value of subnet_numbers above and a VPC CIDR prefix
    # of 10.1.0.0/16, this would produce:
    #   ["10.1.1.0/24", "10.1.2.0/24", "10.1.3.0/24"]
    cidr_blocks = [
    for num in var.subnet_numbers:
    cidrsubnet(data.aws_vpc.example.cidr_block, 8, num)
    ]
  }
}