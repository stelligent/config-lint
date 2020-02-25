# Pass
resource "aws_waf_web_acl" "default_action_type_set_to_block" {
  name        = "foo"
  metric_name = "foo"

  default_action {
    type = "BLOCK"
  }
}

# Pass
resource "aws_waf_web_acl" "default_action_type_set_to_count" {
  name        = "foo"
  metric_name = "foo"

  default_action {
    type = "BLOCK"
  }
}

# Fail
resource "aws_waf_web_acl" "default_action_type_set_to_allow" {
  name        = "foo"
  metric_name = "foo"

  default_action {
    type = "ALLOW"
  }
}
