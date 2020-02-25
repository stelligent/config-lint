## Setup Helper
resource "aws_kms_key" "test_key" {
  enable_key_rotation = true
}

resource "aws_s3_bucket" "test_bucket" {
  acl = "private"
}

resource "aws_iam_role" "test_firehose_role" {
  name = "test_firehose_role"

  assume_role_policy = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Action": "sts:AssumeRole",
      "Principal": {
        "Service": "firehose.amazonaws.com"
      },
      "Effect": "Allow",
      "Sid": ""
    }
  ]
}
EOF
}

resource "aws_kinesis_stream" "test_stream" {
  name        = "test_stream"
  shard_count = 1
}

# Pass
resource "aws_kinesis_firehose_delivery_stream" "extended_s3_configuration_kms_key_arn_is_set" {
  name        = "foo"
  destination = "extended_s3"

  server_side_encryption {
    enabled = true
  }

  extended_s3_configuration {
    role_arn    = "${aws_iam_role.test_firehose_role.arn}"
    bucket_arn  = "${aws_s3_bucket.test_bucket.arn}"
    kms_key_arn = "${aws_kms_key.test_key.arn}"
  }
}

# Fail
resource "aws_kinesis_firehose_delivery_stream" "extended_s3_configuration_kms_key_arn_is_not_set" {
  name        = "foo"
  destination = "extended_s3"

  server_side_encryption {
    enabled = true
  }

  extended_s3_configuration {
    role_arn   = "${aws_iam_role.test_firehose_role.arn}"
    bucket_arn = "${aws_s3_bucket.test_bucket.arn}"
  }
}

# Pass
resource "aws_kinesis_firehose_delivery_stream" "s3_configuration_kms_key_arn_is_set" {
  name        = "terraform-kinesis-firehose-test-stream"
  destination = "s3"

  server_side_encryption {
    enabled = true
  }

  s3_configuration {
    role_arn    = "${aws_iam_role.test_firehose_role.arn}"
    bucket_arn  = "${aws_s3_bucket.test_bucket.arn}"
    kms_key_arn = "${aws_kms_key.test_key.arn}"
  }
}

# Fail
resource "aws_kinesis_firehose_delivery_stream" "s3_configuration_kms_key_arn_is_not_set" {
  name        = "terraform-kinesis-firehose-test-stream"
  destination = "s3"

  server_side_encryption {
    enabled = true
  }

  s3_configuration {
    role_arn   = "${aws_iam_role.test_firehose_role.arn}"
    bucket_arn = "${aws_s3_bucket.test_bucket.arn}"
  }
}

# Pass
resource "aws_kinesis_firehose_delivery_stream" "kinesis_source_configuration_is_set" {
  name        = "foo"
  destination = "extended_s3"

  kinesis_source_configuration {
    kinesis_stream_arn = "${aws_kinesis_stream.test_stream.arn}"
    role_arn           = "${aws_iam_role.test_firehose_role.arn}"
  }

  extended_s3_configuration {
    role_arn    = "${aws_iam_role.test_firehose_role.arn}"
    bucket_arn  = "${aws_s3_bucket.test_bucket.arn}"
    kms_key_arn = "${aws_kms_key.test_key.arn}"
  }
}

# Pass
resource "aws_kinesis_firehose_delivery_stream" "server_side_encryption_enabled_set_to_true" {
  name        = "foo"
  destination = "extended_s3"

  server_side_encryption {
    enabled = true
  }

  extended_s3_configuration {
    role_arn    = "${aws_iam_role.test_firehose_role.arn}"
    bucket_arn  = "${aws_s3_bucket.test_bucket.arn}"
    kms_key_arn = "${aws_kms_key.test_key.arn}"
  }
}

# Fail
resource "aws_kinesis_firehose_delivery_stream" "server_side_encryption_enabled_set_to_false" {
  name        = "foo"
  destination = "extended_s3"

  server_side_encryption {
    enabled = false
  }

  extended_s3_configuration {
    role_arn    = "${aws_iam_role.test_firehose_role.arn}"
    bucket_arn  = "${aws_s3_bucket.test_bucket.arn}"
    kms_key_arn = "${aws_kms_key.test_key.arn}"
  }
}

# Fail
resource "aws_kinesis_firehose_delivery_stream" "server_side_encryption_enabled_not_set" {
  name        = "foo"
  destination = "extended_s3"

  server_side_encryption {
  }

  extended_s3_configuration {
    role_arn    = "${aws_iam_role.test_firehose_role.arn}"
    bucket_arn  = "${aws_s3_bucket.test_bucket.arn}"
    kms_key_arn = "${aws_kms_key.test_key.arn}"
  }
}
