resource "aws_instance" "with_tags" {
    ami = "ami-f2d3638a"
    instance_type = "t2.micro"
    tags {
        "CostCenter" = "1001"
        "Project" = "Web"
    }
}
resource "aws_instance" "without_tags" {
    ami = "ami-f2d3638a"
    instance_type = "m3.medium"
}
