## Setup Helper
resource "aws_vpc" "test_vpc" {
  cidr_block = "10.0.0.0/16"
}

# Pass
resource "aws_subnet" "map_public_ip_on_launch_not_set" {
  vpc_id     = "${aws_vpc.test_vpc.id}"
  cidr_block = "172.2.0.0/24"
}

# Pass
resource "aws_subnet" "map_public_ip_on_launch_set_to_false" {
  vpc_id                  = "${aws_vpc.test_vpc.id}"
  cidr_block              = "172.2.0.0/24"
  map_public_ip_on_launch = false
}

# Warn
resource "aws_subnet" "map_public_ip_on_launch_set_to_true" {
  vpc_id                  = "${aws_vpc.test_vpc.id}"
  cidr_block              = "172.2.0.0/24"
  map_public_ip_on_launch = true
}
