# Test that an OpsWorks application has SSL enabled
# https://www.terraform.io/docs/providers/aws/r/opsworks_application.html#enable_ssl

provider "aws" {
  region = "us-east-1"
}

# PASS: SSL is enabled
resource "aws_opsworks_application" "opsworks_app_ssl_enabled" {
  name        = "foobar application"
  short_name  = "foobar"
  stack_id    = "ExampleStackID"
  type        = "rails"
  description = "This is a Rails application"

  domains = [
    "example.com",
    "sub.example.com",
  ]

  environment {
    key    = "key"
    value  = "value"
    secure = false
  }

  app_source {
    type     = "git"
    revision = "master"
    url      = "https://github.com/example.git"
  }

  enable_ssl = true

  ssl_configuration {
    private_key = "example.key"
    certificate = "example.crt"
  }

  document_root         = "public"
  auto_bundle_on_deploy = true
  rails_env             = "staging"
}

# FAIL: SSL is not enabled
resource "aws_opsworks_application" "opsworks_app_ssl_disabled" {
  name        = "foobar application"
  short_name  = "foobar"
  stack_id    = "ExampleStackID"
  type        = "rails"
  description = "This is a Rails application"

  domains = [
    "example.com",
    "sub.example.com",
  ]

  environment {
    key    = "key"
    value  = "value"
    secure = false
  }

  app_source {
    type     = "git"
    revision = "master"
    url      = "https://github.com/example.git"
  }

  enable_ssl = false

  document_root         = "public"
  auto_bundle_on_deploy = true
  rails_env             = "staging"
}

# FAIL: enable_ssl is not defined
resource "aws_opsworks_application" "opsworks_app_enable_ssl_not_defined" {
  name        = "foobar application"
  short_name  = "foobar"
  stack_id    = "ExampleStackID"
  type        = "rails"
  description = "This is a Rails application"

  domains = [
    "example.com",
    "sub.example.com",
  ]

  environment {
    key    = "key"
    value  = "value"
    secure = false
  }

  app_source {
    type     = "git"
    revision = "master"
    url      = "https://github.com/example.git"
  }

  document_root         = "public"
  auto_bundle_on_deploy = true
  rails_env             = "staging"
}
