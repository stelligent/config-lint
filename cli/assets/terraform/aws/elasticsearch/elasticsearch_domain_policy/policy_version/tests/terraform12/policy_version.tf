# Test that an elasticsearch domain policy is using version 2012-10-17
# https://www.terraform.io/docs/providers/aws/r/elasticsearch_domain_policy.html

provider "aws" {
  region = "us-east-1"
}

# Helper
resource "aws_elasticsearch_domain" "example" {
  domain_name           = "tf-test"
  elasticsearch_version = "2.3"
}

# PASS: Version is 2012-10-17
resource "aws_elasticsearch_domain_policy" "policy_right_version" {
  domain_name = "${aws_elasticsearch_domain.example.domain_name}"

  access_policies = <<POLICIES
{
    "Version": "2012-10-17",
    "Statement": [
        {
            "Action": "es:ListDomainNames",
            "Principal": {
                "AWS": [
                    "arn:aws:iam::1234567890:user/foo"
                ]
            },
            "Effect": "Allow",
            "Resource": "${aws_elasticsearch_domain.example.arn}/*"
        }
    ]
}
POLICIES
}

# FAIL: Version is not 2012-10-17
resource "aws_elasticsearch_domain_policy" "policy_wrong_version" {
  domain_name = "${aws_elasticsearch_domain.example.domain_name}"

  access_policies = <<POLICIES
{
    "Version": "2008-10-17",
    "Statement": [
        {
            "Action": "es:ListDomainNames",
            "Principal": {
                "AWS": [
                    "arn:aws:iam::1234567890:user/foo"
                ]
            },
            "Effect": "Allow",
            "Resource": "${aws_elasticsearch_domain.example.arn}/*"
        }
    ]
}
POLICIES
}
