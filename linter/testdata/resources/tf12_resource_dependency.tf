resource "null_resource" "valid_engine_version_check" {
  count = 1
  #"ERROR: elasticache_version_options can only be: 3.2.6, 4.0.10 or 5.0.5, Recommended version is 5.0.5 as that is the current GA release of Redis, see: https://docs.aws.amazon.com/AmazonElastiCache/latest/red-ug/supported-engine-versions.html" = true
}

resource "aws_elasticache_replication_group" "redis" {
  depends_on                    = [null_resource.valid_engine_version_check]
}
