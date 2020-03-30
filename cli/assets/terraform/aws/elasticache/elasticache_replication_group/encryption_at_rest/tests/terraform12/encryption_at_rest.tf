# Pass
resource "aws_elasticache_replication_group" "at_rest_encryption_enabled_is_set_to_true" {
  replication_group_id          = "foo"
  replication_group_description = "test description"
  node_type                     = "cache.m4.large"
  number_cache_clusters         = 2
  at_rest_encryption_enabled    = true
}

# Fail
resource "aws_elasticache_replication_group" "at_rest_encryption_enabled_is_set_to_false" {
  replication_group_id          = "foo"
  replication_group_description = "test description"
  node_type                     = "cache.m4.large"
  number_cache_clusters         = 2
  at_rest_encryption_enabled    = false
}

# Fail
resource "aws_elasticache_replication_group" "at_rest_encryption_enabled_is_not_set" {
  replication_group_id          = "foo"
  replication_group_description = "test description"
  node_type                     = "cache.m4.large"
  number_cache_clusters         = 2
}
