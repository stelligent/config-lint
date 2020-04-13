# Test for node to node encryption for an elasticsearch domain
# https://www.terraform.io/docs/providers/aws/r/elasticsearch_domain.html#node_to_node_encryption

provider "aws" {
  region = "us-east-1"
}

# PASS: node_to_node_encryption provided and enabled
resource "aws_elasticsearch_domain" "es_domain_node_to_node_encryption_enabled" {
  domain_name           = "example"
  elasticsearch_version = "1.5"

  cluster_config {
    instance_type = "r4.large.elasticsearch"
  }

  node_to_node_encryption {
    enabled = true
  }
}

# FAIL: node_to_node_encryption option not provided
resource "aws_elasticsearch_domain" "es_domain_node_to_node_encryption_not_provided" {
  domain_name           = "example"
  elasticsearch_version = "1.5"

  cluster_config {
    instance_type = "r4.large.elasticsearch"
  }
}

# FAIL: node_to_node_encryption provided but disabled
resource "aws_elasticsearch_domain" "es_domain_node_to_node_encryption_disabled" {
  domain_name           = "example"
  elasticsearch_version = "1.5"

  cluster_config {
    instance_type = "r4.large.elasticsearch"
  }

  node_to_node_encryption {
    enabled = false
  }
}
