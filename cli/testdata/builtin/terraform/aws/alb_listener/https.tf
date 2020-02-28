# Pass
resource "aws_alb_listener" "listener_secure_https_set" {
  load_balancer_arn = "arn:aws:elasticloadbalancing:us-east-1:1234567890:loadbalancer/app/foo"
  port              = "443"
  protocol          = "HTTPS"
  ssl_policy        = "ELBSecurityPolicy-2016-08"
  certificate_arn   = "arn:aws:iam::1234567890:server-certificate/foo"

  default_action {
    type             = "forward"
    target_group_arn = "arn:aws:elasticloadbalancing:us-east-1:1234567890:targetgroup/foo"
  }
}

# Pass
resource "aws_alb_listener" "listener_secure_https_set_lowercase" {
  load_balancer_arn = "arn:aws:elasticloadbalancing:us-east-1:1234567890:loadbalancer/app/foo"
  port              = "443"
  protocol          = "https"
  ssl_policy        = "ELBSecurityPolicy-2016-08"
  certificate_arn   = "arn:aws:iam::1234567890:server-certificate/foo"

  default_action {
    type             = "forward"
    target_group_arn = "arn:aws:elasticloadbalancing:us-east-1:1234567890:targetgroup/foo"
  }
}

# Fail
resource "aws_alb_listener" "port_set_to_80" {
  load_balancer_arn = "arn:aws:elasticloadbalancing:us-east-1:1234567890:loadbalancer/app/foo"
  port              = "80"
  protocol          = "HTTPS"
  ssl_policy        = "ELBSecurityPolicy-2016-08"
  certificate_arn   = "arn:aws:iam::1234567890:server-certificate/foo"

  default_action {
    type             = "forward"
    target_group_arn = "arn:aws:elasticloadbalancing:us-east-1:1234567890:targetgroup/foo"
  }
}

# Fail
resource "aws_alb_listener" "protocol_set_to_http" {
  load_balancer_arn = "arn:aws:elasticloadbalancing:us-east-1:1234567890:loadbalancer/app/foo"
  port              = "443"
  protocol          = "HTTP"
  ssl_policy        = "ELBSecurityPolicy-2016-08"
  certificate_arn   = "arn:aws:iam::1234567890:server-certificate/foo"

  default_action {
    type             = "forward"
    target_group_arn = "arn:aws:elasticloadbalancing:us-east-1:1234567890:targetgroup/foo"
  }
}

# Fail
resource "aws_alb_listener" "ssl_policy_not_set" {
  load_balancer_arn = "arn:aws:elasticloadbalancing:us-east-1:1234567890:loadbalancer/app/foo"
  port              = "443"
  protocol          = "HTTPS"
  certificate_arn   = "arn:aws:iam::1234567890:server-certificate/foo"

  default_action {
    type             = "forward"
    target_group_arn = "arn:aws:elasticloadbalancing:us-east-1:1234567890:targetgroup/foo"
  }
}

# Fail
resource "aws_alb_listener" "certificate_arn_not_set" {
  load_balancer_arn = "arn:aws:elasticloadbalancing:us-east-1:1234567890:loadbalancer/app/foo"
  port              = "443"
  protocol          = "HTTPS"
  ssl_policy        = "ELBSecurityPolicy-2016-08"

  default_action {
    type             = "forward"
    target_group_arn = "arn:aws:elasticloadbalancing:us-east-1:1234567890:targetgroup/foo"
  }
}
