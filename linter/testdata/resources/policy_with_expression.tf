variable "resources" {
    default = "None"
}

resource "aws_iam_role" "role_with_variable" {
    name = "non_compliant"
    assume_role_policy = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
     {
        "Action": "*",
        "Principal": { "Service": "ec2.amazonaws.com" },
        "Effect": "Allow", 
        "Resources": "${var.resources == "foo" ? "foo": "*"}"
     }
  ]
}
EOF
}

