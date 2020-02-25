# Pass
resource "aws_kinesis_stream" "kms_key_id_is_set" {
  name            = "foo"
  shard_count     = 1
  encryption_type = "KMS"
  kms_key_id      = "alias/aws/kinesis"
}

# Warn
resource "aws_kinesis_stream" "kms_key_id_is_not_set" {
  name            = "foo"
  shard_count     = 1
  encryption_type = "KMS"
}
