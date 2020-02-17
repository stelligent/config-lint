resource "aws_s3_bucket" "content" {}

resource "aws_cloudfront_origin_access_identity" "example" {
  comment = "Testing"
}

resource "aws_s3_bucket_policy" "policy" {
    bucket = "${aws_s3_bucket.content.bucket}"
    policy = <<POLICY
{
    "Version": "2012-10-17",
    "Statement": [
        {
            "Sid": " CloudFront Origin Identity",
            "Effect": "Allow",
            "Principal": {
                "AWS": "arn:aws:iam::cloudfront:user/CloudFront Origin Access Identity ${aws_cloudfront_origin_access_identity.example.id}"
            },
            "Action": "s3:GetObject",
            "Resource": "arn:aws:s3:::${aws_s3_bucket.content.bucket}/*"
        }
    ]
}
POLICY
}

resource "aws_s3_bucket_object" "index" {
    bucket = "${aws_s3_bucket.content.bucket}"
    key = "index.html"
    content = "Hello, world!"
    content_type = "text/html"
}

variable "origin_id" {
    default = "website_origin"
}

resource "aws_cloudfront_distribution" "s3_distribution" {
    origin {
        domain_name = "${aws_s3_bucket.content.bucket_regional_domain_name}"
        origin_id = "${var.origin_id}"
        s3_origin_config {
            origin_access_identity = "origin-access-identity/cloudfront/${aws_cloudfront_origin_access_identity.example.id}"
        }
    }
    enabled = true
    default_root_object = "index.html"
    default_cache_behavior {
        allowed_methods  = ["DELETE", "GET", "HEAD", "OPTIONS", "PATCH", "POST", "PUT"]
        cached_methods   = ["GET", "HEAD"]
        target_origin_id = "${var.origin_id}"

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
            locations        = ["US"]
        }
    }
    viewer_certificate {
        cloudfront_default_certificate = true
    }
}
