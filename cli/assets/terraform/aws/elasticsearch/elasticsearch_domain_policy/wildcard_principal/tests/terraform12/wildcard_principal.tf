# Test that an elasticsearch domain policy is not using a wildcard principal
# https://www.terraform.io/docs/providers/aws/r/elasticsearch_domain_policy.html

provider "aws" {
  region = "us-east-1"
}

# Helper
resource "aws_elasticsearch_domain" "example" {
  domain_name           = "tf-test"
  elasticsearch_version = "2.3"
}

# PASS: Allow principal does not contain a wildcard
resource "aws_elasticsearch_domain_policy" "policy_allow_principal_no_wildcard" {
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

# PASS: Deny principal doesn't contain a wildcard
resource "aws_elasticsearch_domain_policy" "policy_allow_principal_no_wildcard" {
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
            "Effect": "Deny",
            "Resource": "${aws_elasticsearch_domain.example.arn}/*"
        }
    ]
}
POLICIES
}

# PASS: Deny principal contains a wildcard
resource "aws_elasticsearch_domain_policy" "policy_deny_principal_contains_wildcard" {
  domain_name = "${aws_elasticsearch_domain.example.domain_name}"

  access_policies = <<POLICIES
{
    "Version": "2012-10-17",
    "Statement": [
        {
            "Action": "es:ListDomainNames",
            "Principal": {
                "AWS": [
                    "arn:aws:iam::1234567890:user/foo*"
                ]
            },
            "Effect": "Allow",
            "Resource": "${aws_elasticsearch_domain.example.arn}/*"
        }
    ]
}
POLICIES
}

# FAIL: Allow principal contains a wildcard
resource "aws_elasticsearch_domain_policy" "policy_allow_principal_contains_wildcard" {
  domain_name = "${aws_elasticsearch_domain.example.domain_name}"

  access_policies = <<POLICIES
{
    "Version": "2012-10-17",
    "Statement": [
        {
            "Action": "es:ListDomainNames",
            "Principal": {
                "AWS": [
                    "arn:aws:iam::1234567890:user/foo*"
                ]
            },
            "Effect": "Allow",
            "Resource": "${aws_elasticsearch_domain.example.arn}/*"
        }
    ]
}
POLICIES
}

# FAIL: Principal is a wildcard
resource "aws_elasticsearch_domain_policy" "policy_allow_principal_is_wildcard" {
  domain_name = "${aws_elasticsearch_domain.example.domain_name}"

  access_policies = <<POLICIES
{
    "Version": "2012-10-17",
    "Statement": [
        {
            "Action": "es:ListDomainNames",
            "Principal": "*",
            "Effect": "Allow",
            "Resource": "${aws_elasticsearch_domain.example.arn}/*"
        }
    ]
}
POLICIES
}