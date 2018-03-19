resource "aws_instance" "first" {
    ami = "ami-f2d3638a"
    instance_type = "t2.micro"
    tags {
        "CostCenter" = "1001"
        "Project" = "Web"
    }
}
resource "aws_instance" "second" {
    ami = "ami-f2d3638a"
    instance_type = "m3.medium"
    tags {
    Department = "Operations"
        CostCenter = "2001"
        "Project" = "Web"
    }
}
resource "aws_instance" "third" {
    ami = "ami-f2d3638b"
    instance_type = "c4.large"
}
resource "aws_instance" "foo" {
    ami = "ami-f2d3638b"
    instance_type = "c4.large"
    tags {
        Foo = "Foo"
        "Project" = "Web"
    }
}
resource "aws_iam_role" "role1" {
    name = "role1"
    assume_role_policy = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Action": "sts:AssumeRole",
      "Principal": {
        "Service": "ec2.amazonaws.com"
      },
      "Effect": "Allow",
      "Sid": ""
    }
  ]
}
EOF
}
resource "aws_iam_role" "role2" {
    name = "non_compliant"
    assume_role_policy = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
     {
        "Action": "*",
        "Principal": { "Service": "ec2.amazonaws.com" },
        "Effect": "Allow",
        "Resources": "*"
     }
  ]
}
EOF
}
