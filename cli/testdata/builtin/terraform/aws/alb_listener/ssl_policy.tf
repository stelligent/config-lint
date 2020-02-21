# Pass
resource "aws_alb_listener" "ssl_policy_set_to_ELBSecurityPolicy-2016-08" {
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
resource "aws_alb_listener" "ssl_policy_set_to_ELBSecurityPolicy-TLS-1-2-2017-01" {
  load_balancer_arn = "arn:aws:elasticloadbalancing:us-east-1:1234567890:loadbalancer/app/foo"
  port              = "443"
  protocol          = "HTTPS"
  ssl_policy        = "ELBSecurityPolicy-TLS-1-2-2017-01"
  certificate_arn   = "arn:aws:iam::1234567890:server-certificate/foo"

  default_action {
    type             = "forward"
    target_group_arn = "arn:aws:elasticloadbalancing:us-east-1:1234567890:targetgroup/foo"
  }
}

# Pass
resource "aws_alb_listener" "ssl_policy_set_to_ELBSecurityPolicy-TLS-1-1-2017-01" {
  load_balancer_arn = "arn:aws:elasticloadbalancing:us-east-1:1234567890:loadbalancer/app/foo"
  port              = "443"
  protocol          = "HTTPS"
  ssl_policy        = "ELBSecurityPolicy-TLS-1-1-2017-01"
  certificate_arn   = "arn:aws:iam::1234567890:server-certificate/foo"

  default_action {
    type             = "forward"
    target_group_arn = "arn:aws:elasticloadbalancing:us-east-1:1234567890:targetgroup/foo"
  }
}

# Fail
resource "aws_alb_listener" "ssl_policy_set_to_ELBSecurityPolicy-2015-05" {
  load_balancer_arn = "arn:aws:elasticloadbalancing:us-east-1:1234567890:loadbalancer/app/foo"
  port              = "443"
  protocol          = "HTTPS"
  ssl_policy        = "ELBSecurityPolicy-2015-05"
  certificate_arn   = "arn:aws:iam::1234567890:server-certificate/foo"

  default_action {
    type             = "forward"
    target_group_arn = "arn:aws:elasticloadbalancing:us-east-1:1234567890:targetgroup/foo"
  }
}

# Fail
resource "aws_alb_listener" "ssl_policy_set_to_ELBSecurityPolicy-2015-03" {
  load_balancer_arn = "arn:aws:elasticloadbalancing:us-east-1:1234567890:loadbalancer/app/foo"
  port              = "443"
  protocol          = "HTTPS"
  ssl_policy        = "ELBSecurityPolicy-2015-03"
  certificate_arn   = "arn:aws:iam::1234567890:server-certificate/foo"

  default_action {
    type             = "forward"
    target_group_arn = "arn:aws:elasticloadbalancing:us-east-1:1234567890:targetgroup/foo"
  }
}

# Fail
resource "aws_alb_listener" "ssl_policy_set_to_ELBSecurityPolicy-2015-02" {
  load_balancer_arn = "arn:aws:elasticloadbalancing:us-east-1:1234567890:loadbalancer/app/foo"
  port              = "443"
  protocol          = "HTTPS"
  ssl_policy        = "ELBSecurityPolicy-2015-02"
  certificate_arn   = "arn:aws:iam::1234567890:server-certificate/foo"

  default_action {
    type             = "forward"
    target_group_arn = "arn:aws:elasticloadbalancing:us-east-1:1234567890:targetgroup/foo"
  }
}

# Fail
resource "aws_alb_listener" "ssl_policy_set_to_ELBSecurityPolicy-2014-10" {
  load_balancer_arn = "arn:aws:elasticloadbalancing:us-east-1:1234567890:loadbalancer/app/foo"
  port              = "443"
  protocol          = "HTTPS"
  ssl_policy        = "ELBSecurityPolicy-2014-10"
  certificate_arn   = "arn:aws:iam::1234567890:server-certificate/foo"

  default_action {
    type             = "forward"
    target_group_arn = "arn:aws:elasticloadbalancing:us-east-1:1234567890:targetgroup/foo"
  }
}

# Fail
resource "aws_alb_listener" "ssl_policy_set_to_ELBSecurityPolicy-2014-01" {
  load_balancer_arn = "arn:aws:elasticloadbalancing:us-east-1:1234567890:loadbalancer/app/foo"
  port              = "443"
  protocol          = "HTTPS"
  ssl_policy        = "ELBSecurityPolicy-2014-01"
  certificate_arn   = "arn:aws:iam::1234567890:server-certificate/foo"

  default_action {
    type             = "forward"
    target_group_arn = "arn:aws:elasticloadbalancing:us-east-1:1234567890:targetgroup/foo"
  }
}

# Fail
resource "aws_alb_listener" "ssl_policy_set_to_ELBSecurityPolicy-2011-08" {
  load_balancer_arn = "arn:aws:elasticloadbalancing:us-east-1:1234567890:loadbalancer/app/foo"
  port              = "443"
  protocol          = "HTTPS"
  ssl_policy        = "ELBSecurityPolicy-2011-08"
  certificate_arn   = "arn:aws:iam::1234567890:server-certificate/foo"

  default_action {
    type             = "forward"
    target_group_arn = "arn:aws:elasticloadbalancing:us-east-1:1234567890:targetgroup/foo"
  }
}
