# Test that a cloudfront_distribution resource is using OAI or origin_protocol_policy is using https-only
# https://www.terraform.io/docs/providers/aws/r/cloudfront_distribution.html#origin_access_identity
# https://www.terraform.io/docs/providers/aws/r/cloudfront_distribution.html#origin_protocol_policy

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

# PASS: OIA is used
resource "aws_cloudfront_distribution" "custom_origin_config_not_set" {
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

# PASS: origin_protocol_policy is https-only
resource "aws_cloudfront_distribution" "custom_origin_config_set_to_https-only" {
  enabled = true

  origin {
    domain_name = var.test_domain_s3_location
    origin_id   = var.test_origin_id

    custom_origin_config {
      http_port              = 80
      https_port             = 443
      origin_protocol_policy = "https-only"
      origin_ssl_protocols   = ["TLSv1.2"]
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

# FAIL: origin_protocol_policy is not https-only
resource "aws_cloudfront_distribution" "custom_origin_config_set_to_http-only" {
  enabled = true

  origin {
    domain_name = var.test_domain_s3_location
    origin_id   = var.test_origin_id

    custom_origin_config {
      http_port              = 80
      https_port             = 443
      origin_protocol_policy = "http-only"
      origin_ssl_protocols   = ["TLSv1.2"]
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

# FAIL: origin_protocol_policy is not https-only
resource "aws_cloudfront_distribution" "custom_origin_config_set_to_match-viewer" {
  enabled = true

  origin {
    domain_name = var.test_domain_s3_location
    origin_id   = var.test_origin_id

    custom_origin_config {
      http_port              = 80
      https_port             = 443
      origin_protocol_policy = "match-viewer"
      origin_ssl_protocols   = ["TLSv1.2"]
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
