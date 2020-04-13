# Test for encryption at rest options on an elasticsearch domain
# https://www.terraform.io/docs/providers/aws/r/elasticsearch_domain.html#encrypt_at_rest

provider "aws" {
  region = "us-east-1"
}

# FAIL: encrypt_at_rest not provided
resource "aws_elasticsearch_domain" "es_domain_encryption_at_rest_not_provided" {
  domain_name           = "example"
  elasticsearch_version = "1.5"

  cluster_config {
    instance_type = "r4.large.elasticsearch"
  }
}

# FAIL: encrypt_at_rest provided but disabled
resource "aws_elasticsearch_domain" "es_domain_encryption_at_rest_not_enabled" {
  domain_name           = "example"
  elasticsearch_version = "1.5"

  cluster_config {
    instance_type = "r4.large.elasticsearch"
  }

  encrypt_at_rest {
    enabled = false
  }
}

# PASS: encrypt_at_rest provided and enabled
resource "aws_elasticsearch_domain" "es_domain_encryption_at_rest_enabled" {
  domain_name           = "example"
  elasticsearch_version = "1.5"

  cluster_config {
    instance_type = "r4.large.elasticsearch"
  }

  encrypt_at_rest {
    enabled = true
  }
}
