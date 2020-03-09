## Setup Helper
variable "test_domain_s3_location" {
  default = "http://foo.s3-website-us-east-1.amazonaws.com"
}

variable "test_origin_id" {
  default = "fooOrigin"
}

variable "test_logging_bucket" {
  default = "foologs.s3.amazonaws.com"
}

variable "test_logging_prefix" {
  default = "aws_cloudfront_distribution"
}

# Pass
resource "aws_cloudfront_distribution" "logging_enabled" {
  enabled = true

  origin {
    domain_name = var.test_domain_s3_location
    origin_id   = var.test_origin_id

    s3_origin_config {
      origin_access_identity = "origin-access-identity/cloudfront/ABCDEFG1234567"
    }
  }

  logging_config {
    include_cookies = false
    bucket          = var.test_logging_bucket
    prefix          = var.test_logging_prefix
  }

  default_cache_behavior {
    allowed_methods  = ["DELETE", "GET", "HEAD", "OPTIONS", "PATCH", "POST", "PUT"]
    cached_methods   = ["GET", "HEAD"]
    target_origin_id = "fooOrigin"

    forwarded_values {
      query_string = false

      cookies {
        forward = "none"
      }
    }

    viewer_protocol_policy = "allow-all"
    min_ttl                = 0
    default_ttl            = 3600
    max_ttl                = 86400
  }

  restrictions {
    geo_restriction {
      restriction_type = "whitelist"
      locations        = ["US", "CA", "GB", "DE"]
    }
  }

  viewer_certificate {
    cloudfront_default_certificate = true
  }
}

# Fail
resource "aws_cloudfront_distribution" "logging_disabled" {
  enabled = true

  origin {
    domain_name = var.test_domain_s3_location
    origin_id   = var.test_origin_id

    s3_origin_config {
      origin_access_identity = "origin-access-identity/cloudfront/ABCDEFG1234567"
    }
  }

  default_cache_behavior {
    allowed_methods  = ["DELETE", "GET", "HEAD", "OPTIONS", "PATCH", "POST", "PUT"]
    cached_methods   = ["GET", "HEAD"]
    target_origin_id = "fooOrigin"

    forwarded_values {
      query_string = false

      cookies {
        forward = "none"
      }
    }

    viewer_protocol_policy = "allow-all"
    min_ttl                = 0
    default_ttl            = 3600
    max_ttl                = 86400
  }

  restrictions {
    geo_restriction {
      restriction_type = "whitelist"
      locations        = ["US", "CA", "GB", "DE"]
    }
  }

  viewer_certificate {
    cloudfront_default_certificate = true
  }
}
