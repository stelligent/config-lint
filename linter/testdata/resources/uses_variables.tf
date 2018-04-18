variable "instance_type" {
  default = "t2.micro"
}

variable "ami" {
  default = "ami-f2d3638a"
}

variable "project" {
  default = "demo"
}

resource "aws_instance" "first" {
  ami = "${var.ami}"
  instance_type = "${var.instance_type}"
  tags = {
    project = "${var.project}"
  }
}
