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
  role        = "${aws_iam_role.build.name}"

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

# project with encryption. 
# Should Pass
resource "aws_codebuild_project" "pass_encryption" {
  name          = "pass_encryption_project"
  description   = "pass_encryption_project"
  build_timeout = "5"
  service_role  = "${aws_iam_role.build.arn}"

  artifacts {
    type = "NO_ARTIFACTS"
  }

  environment {
    compute_type = "BUILD_GENERAL1_SMALL"
    image        = "aws/codebuild/nodejs:6.3.1"
    type         = "LINUX_CONTAINER"
  }

  source {
    type            = "GITHUB"
    location        = "https://gist.github.com/blahblahblah.git"
  }

  encryption_key    = "iamanencryptionkey"
}

# project without encryption. 
# Should fail
resource "aws_codebuild_project" "fail_encryption" {
  name          = "fail_encryption_project"
  description   = "fail_encryption_project"
  build_timeout = "5"
  service_role  = "${aws_iam_role.build.arn}"

  artifacts {
    type = "NO_ARTIFACTS"
  }

  environment {
    compute_type = "BUILD_GENERAL1_SMALL"
    image        = "aws/codebuild/nodejs:6.3.1"
    type         = "LINUX_CONTAINER"
  }

  source {
    type            = "GITHUB"
    location        = "https://gist.github.com/blahblahblah.git"
  }
}

