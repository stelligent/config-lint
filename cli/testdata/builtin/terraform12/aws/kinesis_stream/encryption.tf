# Pass
resource "aws_kinesis_stream" "encryption_type_set_to_kms" {
  name            = "foo"
  shard_count     = 1
  encryption_type = "KMS"
}

# Fail
resource "aws_kinesis_stream" "encryption_type_set_to_none" {
  name            = "foo"
  shard_count     = 1
  encryption_type = "NONE"
}

# Fail
resource "aws_kinesis_stream" "encryption_type_not_set" {
  name        = "foo"
  shard_count = 1
}
