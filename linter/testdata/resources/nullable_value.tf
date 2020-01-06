variable "is_jump_host" {
  default = false
}

resource "aws_instance" "nullable" {
  ami = "ami-f2d3638a"
  instance_type = "t2.micro"
  key_name = var.is_jump_host ? "my_aws_key" : null
}