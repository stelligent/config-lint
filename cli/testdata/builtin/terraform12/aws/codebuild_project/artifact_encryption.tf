# Resource required for creating project
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

# Resource required for creating project
resource "aws_iam_role_policy" "build" {
  role = aws_iam_role.build.name

  policy = <<POLICY
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Effect": "Allow",
      "Resource": [
        "*"
      ],
      "Action": [
        "logs:CreateLogGroup",
        "logs:CreateLogStream",
        "logs:PutLogEvents"
      ]
    }
  ]
}
POLICY
}

# project with encryption key, but artifact encryption is disabled. 
# Should fail
resource "aws_codebuild_project" "fail_artifact_encryption" {
  name          = "fail_artifact_encryption_project"
  description   = "fail_artifact_encryption_project"
  build_timeout = "5"
  service_role  = aws_iam_role.build.arn

  artifacts {
    type                = "S3"
    encryption_disabled = true
  }

  environment {
    compute_type = "BUILD_GENERAL1_SMALL"
    image        = "aws/codebuild/nodejs:6.3.1"
    type         = "LINUX_CONTAINER"
  }

  source {
    type     = "GITHUB"
    location = "https://gist.github.com/blahblahblah.git"
  }
  encryption_key = "iamanencryptionkey"
}

# project with encryption key, and artifact encryption is not disabled. 
# Should pass
resource "aws_codebuild_project" "pass_artifact_encryption" {
  name          = "pass_artifact_encryption_project"
  description   = "pass_artifact_encryption_project"
  build_timeout = "5"
  service_role  = aws_iam_role.build.arn

  artifacts {
    type = "S3"
  }

  environment {
    compute_type = "BUILD_GENERAL1_SMALL"
    image        = "aws/codebuild/nodejs:6.3.1"
    type         = "LINUX_CONTAINER"
  }

  source {
    type     = "GITHUB"
    location = "https://gist.github.com/blahblahblah.git"
  }

  encryption_key = "iamanencryptionkey"
}

# project with encryption key, but secondary artifact encryption is disabled. 
# Should fail
resource "aws_codebuild_project" "fail_secondary_artifact_encryption" {
  name          = "fail_secondary_artifact_encryption_project"
  description   = "fail_secondary_artifact_encryption_project"
  build_timeout = "5"
  service_role  = aws_iam_role.build.arn

  artifacts {
    type = "S3"
  }

  secondary_artifacts {
    type                = "S3"
    artifact_identifier = "i_am_an_identifier"
    encryption_disabled = true
  }

  environment {
    compute_type = "BUILD_GENERAL1_SMALL"
    image        = "aws/codebuild/nodejs:6.3.1"
    type         = "LINUX_CONTAINER"
  }

  source {
    type     = "GITHUB"
    location = "https://gist.github.com/blahblahblah.git"
  }

  encryption_key = "iamanencryptionkey"
}

# project with encryption key, and secondary artifact encryption is not disabled. 
# Should pass
resource "aws_codebuild_project" "pass_secondary_artifact_encryption" {
  name          = "pass_secondary_artifact_encryption_project"
  description   = "pass_secondary_artifact_encryption_project"
  build_timeout = "5"
  service_role  = aws_iam_role.build.arn

  artifacts {
    type = "S3"
  }

  secondary_artifacts {
    type                = "S3"
    artifact_identifier = "i_am_an_identifier"
  }

  environment {
    compute_type = "BUILD_GENERAL1_SMALL"
    image        = "aws/codebuild/nodejs:6.3.1"
    type         = "LINUX_CONTAINER"
  }

  source {
    type     = "GITHUB"
    location = "https://gist.github.com/blahblahblah.git"
  }

  encryption_key = "iamanencryptionkey"
}

# project with encryption key, but S3 encryption is disabled. 
# Should fail
resource "aws_codebuild_project" "fail_s3_encryption" {
  name          = "fail_s3_encryption_project"
  description   = "fail_s3_encryption_project"
  build_timeout = "5"
  service_role  = aws_iam_role.build.arn

  artifacts {
    type = "S3"
  }

  s3_logs {
    status              = "ENABLED"
    location            = "iamabucket/path/to/a/location"
    encryption_disabled = true
  }

  environment {
    compute_type = "BUILD_GENERAL1_SMALL"
    image        = "aws/codebuild/nodejs:6.3.1"
    type         = "LINUX_CONTAINER"
  }

  source {
    type     = "GITHUB"
    location = "https://gist.github.com/blahblahblah.git"
  }

  encryption_key = "iamanencryptionkey"
}

# project with encryption key, and S3 encryption is not disabled. 
# Should pass
resource "aws_codebuild_project" "pass_s3_encryption" {
  name          = "pass_s3_encryption_project"
  description   = "pass_s3_encryption_project"
  build_timeout = "5"
  service_role  = aws_iam_role.build.arn

  artifacts {
    type = "S3"
  }

  s3_logs {
    status   = "ENABLED"
    location = "iamabucket/path/to/a/location"
  }

  environment {
    compute_type = "BUILD_GENERAL1_SMALL"
    image        = "aws/codebuild/nodejs:6.3.1"
    type         = "LINUX_CONTAINER"
  }

  source {
    type     = "GITHUB"
    location = "https://gist.github.com/blahblahblah.git"
  }

  encryption_key = "iamanencryptionkey"
}
