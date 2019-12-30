variable "instance_type" {
  default = "t2.micro"
}

variable "ami" {
  default = "ami-f2d3638a"
}

variable "project" {
  default = "demo"
}

variable "list_variable" {
  default = [ "foo", "bar" ]
}

variable "default_tags" {
  default = {
    project = "demo"
    environment = "test"
  }
}

variable "environment" {
  default = "test"
}

variable "department" {}

resource "aws_instance" "first" {
  ami = "${var.ami}"
  instance_type = "${var.instance_type}"
  tags = {
    project = "${var.project}"
    environment = "${lookup(var.default_tags,"environment","dev")}"
    comment = "${var.list_variable[1]}"
    department = "${var.department}"
  }
}

resource "aws_vpc" "example" {
  cidr_block = "10.0.0.0/16"
}