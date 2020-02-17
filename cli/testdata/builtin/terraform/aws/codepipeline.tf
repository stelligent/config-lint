data "aws_kms_alias" "pipeline" {
  name = "alias/pipeline-key"
}

resource "aws_s3_bucket" "project" {
    versioning {
      enabled = true
    }
}

resource "aws_s3_bucket_object" "zip" {
    bucket = "${aws_s3_bucket.project.id}"
    key = "app.zip"
    source = "app.zip"
    kms_key_id = "${data.aws_kms_alias.pipeline.arn}"
}

resource "aws_iam_role" "build" {
  name = "build"

  assume_role_policy = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Effect": "Allow",
      "Principal": {
        "Service": "codebuild.amazonaws.com"
      },
      "Action": "sts:AssumeRole"
    }
  ]
}
EOF
}

resource "aws_iam_role_policy" "build" {
  role        = "${aws_iam_role.build.name}"

  policy = <<POLICY
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Effect": "Allow",
      "Action": [
        "logs:CreateLogGroup",
        "logs:CreateLogStream",
        "logs:PutLogEvents"
      ],
      "Resource": [
        "*"
      ]
    },
    {
      "Effect": "Allow",
      "Action": [
        "s3:*"
      ],
      "Resource": [
        "${aws_s3_bucket.artifact.arn}",
        "${aws_s3_bucket.artifact.arn}/*"
      ]
    },
    {
      "Effect": "Allow",
      "Action": [
        "kms:Encrypt",
        "kms:Decrypt",
        "kms:ReEncrypt*",
        "kms:GenerateDataKey*",
        "kms:DescribeKey"
      ],
      "Resource": "*"
    }
  ]
}
POLICY
}

resource "aws_codebuild_project" "example" {
  name         = "example_project"
  description  = "example_project"
  build_timeout      = "5"
  service_role = "${aws_iam_role.build.arn}"

  artifacts {
    type = "CODEPIPELINE"
  }

  environment {
    compute_type = "BUILD_GENERAL1_SMALL"
    image        = "aws/codebuild/golang:1.10"
    type         = "LINUX_CONTAINER"
  }

  source {
    type = "CODEPIPELINE"
  }

  encryption_key = "${data.aws_kms_alias.pipeline.arn}"
}

resource "aws_s3_bucket" "artifact" {}

resource "aws_iam_role" "pipeline" {
  name = "test-role"

  assume_role_policy = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Effect": "Allow",
      "Principal": {
        "Service": "codepipeline.amazonaws.com"
      },
      "Action": "sts:AssumeRole"
    }
  ]
}
EOF
}

resource "aws_iam_role_policy" "pipeline" {
  name = "codepipeline_policy"
  role = "${aws_iam_role.pipeline.id}"
  policy = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Effect":"Allow",
      "Action": [
        "s3:*"
      ],
      "Resource": [
        "${aws_s3_bucket.project.arn}",
        "${aws_s3_bucket.project.arn}/*",
        "${aws_s3_bucket.artifact.arn}",
        "${aws_s3_bucket.artifact.arn}/*"
      ]
    },
    {
      "Effect": "Allow",
      "Action": [
        "codebuild:BatchGetBuilds",
        "codebuild:StartBuild"
      ],
      "Resource": "*"
    },
    {
      "Effect": "Allow",
      "Action": [
        "kms:Encrypt",
        "kms:Decrypt",
        "kms:ReEncrypt*",
        "kms:GenerateDataKey*",
        "kms:DescribeKey"
      ],
      "Resource": "*"
    }
  ]
}
EOF
}

resource "aws_codepipeline" "demo" {
  name     = "tf-test-pipeline"
  role_arn = "${aws_iam_role.pipeline.arn}"

  artifact_store {
    location = "${aws_s3_bucket.artifact.bucket}"
    type     = "S3"
    encryption_key {
      id   = "${data.aws_kms_alias.pipeline.arn}"
      type = "KMS"
    }
  }

  stage {
    name = "Source"

    action {
      name             = "Source"
      category         = "Source"
      owner            = "AWS"
      provider         = "S3"
      version          = "1"
      output_artifacts = ["source"]

      configuration {
        S3Bucket      = "${aws_s3_bucket.project.bucket}"
        S3ObjectKey   = "app.zip"
      }
    }
  }

  stage {
    name = "Lint"

    action {
      name            = "Build"
      category        = "Build"
      owner           = "AWS"
      provider        = "CodeBuild"
      input_artifacts = ["source"]
      version         = "1"

      configuration {
        ProjectName = "${aws_codebuild_project.example.id}"
      }
    }
  }
}

