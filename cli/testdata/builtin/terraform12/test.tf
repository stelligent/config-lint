resource "aws_cloudtrail" "object_logging_enabled" {
  name                          = "tf-trail-foobar"
  s3_bucket_name                = "nwm-cloudtrail-logs"
  s3_key_prefix                 = "prefix"
  include_global_service_events = false
  event_selector {
    read_write_type           = "All"
    include_management_events = true
    data_resource {
      type   = "AWS::S3::Object"
      values = ["arn:aws:s3:::"]
    }
  }
}

resource "aws_cloudtrail" "object_logging_enabled" {
  name                          = "tf-trail-foobar"
  s3_bucket_name                = "nwm-cloudtrail-logs"
  s3_key_prefix                 = "prefix"
  include_global_service_events = false
  event_selector {
    read_write_type           = "All"
    include_management_events = true
    data_resource {
      type   = "wrong"
      values = ["arn:aws:s3:::"]
    }
  }
}
resource "aws_cloudtrail" "object_logging_enabled" {
  name                          = "tf-trail-foobar"
  s3_bucket_name                = "nwm-cloudtrail-logs"
  s3_key_prefix                 = "prefix"
  include_global_service_events = false
  event_selector {
    read_write_type           = "All"
    include_management_events = true
  }
}
