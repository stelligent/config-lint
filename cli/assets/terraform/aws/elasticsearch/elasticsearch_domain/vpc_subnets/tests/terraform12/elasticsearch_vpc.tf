# Test that elasticsearch domain is in a VPC using vpc_options
# https://www.terraform.io/docs/providers/aws/r/elasticsearch_domain.html#vpc_options

provider "aws" {
  region = "us-east-1"
}

# FAIL: vpc_options is not provided
resource "aws_elasticsearch_domain" "es_domain_vpc_options_not_provided" {
  domain_name           = "example"
  elasticsearch_version = "1.5"

  cluster_config {
    instance_type = "r4.large.elasticsearch"
  }
}

# PASS: vpc_options is provided
resource "aws_elasticsearch_domain" "es_domain_vpc_options_provided" {
  domain_name           = "example"
  elasticsearch_version = "1.5"

  cluster_config {
    instance_type = "r4.large.elasticsearch"
  }

  vpc_options {
    subnet_ids = [
      "Subnet1",
      "Subnet2"
    ]
  }
}
