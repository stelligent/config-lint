# Test that at-rest encryption is enabled with a KMS key
# https://www.terraform.io/docs/providers/aws/r/elasticache_replication_group.html#at_rest_encryption_enabled
# https://www.terraform.io/docs/providers/aws/r/elasticache_replication_group.html#kms_key_id

provider "aws" {
  region = "us-east-1"
}

# PASS: Encryption at rest is enabled with a KMS key
resource "aws_elasticache_replication_group" "at_rest_encryption_enabled_is_set_to_true" {
  replication_group_id          = "foo"
  replication_group_description = "test description"
  node_type                     = "cache.m4.large"
  number_cache_clusters         = 2
  at_rest_encryption_enabled    = true
  kms_key_id                    = "arn:aws:kms:us-west-2:111122223333:key/1234abcd-12ab-34cd-56ef-1234567890ab"
}

# FAIL: Encryption at rest is disabled
# WARN: KMS key is not provided
resource "aws_elasticache_replication_group" "at_rest_encryption_enabled_is_set_to_false" {
  replication_group_id          = "foo"
  replication_group_description = "test description"
  node_type                     = "cache.m4.large"
  number_cache_clusters         = 2
  at_rest_encryption_enabled    = false
}

# FAIL: Encryption at rest is not specified
# WARN: KMS key is not provided
resource "aws_elasticache_replication_group" "at_rest_encryption_enabled_is_not_set" {
  replication_group_id          = "foo"
  replication_group_description = "test description"
  node_type                     = "cache.m4.large"
  number_cache_clusters         = 2
}
