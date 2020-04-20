# Test that an api_gateway_domain_name is using TLS 1.2
# https://www.terraform.io/docs/providers/aws/r/api_gateway_domain_name.html#security_policy

provider "aws" {
  region = "us-east-1"
}

# PASS: security_policy is set to TLS 1.2
resource "aws_api_gateway_domain_name" "api_gw_domain_using_tls1_2" {
  domain_name = "api.example.com"

  endpoint_configuration {
    types = ["REGIONAL"]
  }

  security_policy = "TLS_1_2"
}

# FAIL: security_policy is not defined
resource "aws_api_gateway_domain_name" "api_gw_security_policy_not_set" {
  domain_name = "api.example.com"

  endpoint_configuration {
    types = ["REGIONAL"]
  }
}

# FAIL: security_policy is set to TLS 1.0 
resource "aws_api_gateway_domain_name" "api_gw_domain_using_tls1_0" {
  domain_name = "api.example.com"

  endpoint_configuration {
    types = ["REGIONAL"]
  }

  security_policy = "TLS_1_0"
}
