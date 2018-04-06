resource "aws_instance" "first" {
    ami = "ami-f2d3638a"
    instance_type = "t2.micro"
}
