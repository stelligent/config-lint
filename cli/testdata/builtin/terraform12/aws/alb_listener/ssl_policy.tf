## Setup Helper
resource "aws_vpc" "test_vpc" {
  cidr_block = "10.0.0.0/16"
}

resource "aws_acm_certificate" "test_cert" {
  domain_name       = "foobar.com"
  validation_method = "DNS"

  lifecycle {
    create_before_destroy = true
  }
}

resource "aws_lb" "test_lb" {
}

resource "aws_lb_target_group" "test_lb_target_group" {
  vpc_id = aws_vpc.test_vpc.id
}

# Pass
resource "aws_alb_listener" "ssl_policy_set_to_ELBSecurityPolicy-2016-08" {
  load_balancer_arn = aws_lb.test_lb.arn
  port              = "443"
  protocol          = "HTTPS"
  ssl_policy        = "ELBSecurityPolicy-2016-08"
  certificate_arn   = aws_acm_certificate.test_cert.arn

  default_action {
    type             = "forward"
    target_group_arn = aws_lb_target_group.test_lb_target_group.arn
  }
}

# Pass
resource "aws_alb_listener" "ssl_policy_set_to_ELBSecurityPolicy-TLS-1-2-2017-01" {
  load_balancer_arn = aws_lb.test_lb.arn
  port              = "443"
  protocol          = "HTTPS"
  ssl_policy        = "ELBSecurityPolicy-TLS-1-2-2017-01"
  certificate_arn   = aws_acm_certificate.test_cert.arn

  default_action {
    type             = "forward"
    target_group_arn = aws_lb_target_group.test_lb_target_group.arn
  }
}

# Pass
resource "aws_alb_listener" "ssl_policy_set_to_ELBSecurityPolicy-TLS-1-1-2017-01" {
  load_balancer_arn = aws_lb.test_lb.arn
  port              = "443"
  protocol          = "HTTPS"
  ssl_policy        = "ELBSecurityPolicy-TLS-1-1-2017-01"
  certificate_arn   = aws_acm_certificate.test_cert.arn

  default_action {
    type             = "forward"
    target_group_arn = aws_lb_target_group.test_lb_target_group.arn
  }
}

# Fail
resource "aws_alb_listener" "ssl_policy_set_to_ELBSecurityPolicy-2015-05" {
  load_balancer_arn = aws_lb.test_lb.arn
  port              = "443"
  protocol          = "HTTPS"
  ssl_policy        = "ELBSecurityPolicy-2015-05"
  certificate_arn   = aws_acm_certificate.test_cert.arn

  default_action {
    type             = "forward"
    target_group_arn = aws_lb_target_group.test_lb_target_group.arn
  }
}

# Fail
resource "aws_alb_listener" "ssl_policy_set_to_ELBSecurityPolicy-2015-03" {
  load_balancer_arn = aws_lb.test_lb.arn
  port              = "443"
  protocol          = "HTTPS"
  ssl_policy        = "ELBSecurityPolicy-2015-03"
  certificate_arn   = aws_acm_certificate.test_cert.arn

  default_action {
    type             = "forward"
    target_group_arn = aws_lb_target_group.test_lb_target_group.arn
  }
}

# Fail
resource "aws_alb_listener" "ssl_policy_set_to_ELBSecurityPolicy-2015-02" {
  load_balancer_arn = aws_lb.test_lb.arn
  port              = "443"
  protocol          = "HTTPS"
  ssl_policy        = "ELBSecurityPolicy-2015-02"
  certificate_arn   = aws_acm_certificate.test_cert.arn

  default_action {
    type             = "forward"
    target_group_arn = aws_lb_target_group.test_lb_target_group.arn
  }
}

# Fail
resource "aws_alb_listener" "ssl_policy_set_to_ELBSecurityPolicy-2014-10" {
  load_balancer_arn = aws_lb.test_lb.arn
  port              = "443"
  protocol          = "HTTPS"
  ssl_policy        = "ELBSecurityPolicy-2014-10"
  certificate_arn   = aws_acm_certificate.test_cert.arn

  default_action {
    type             = "forward"
    target_group_arn = aws_lb_target_group.test_lb_target_group.arn
  }
}

# Fail
resource "aws_alb_listener" "ssl_policy_set_to_ELBSecurityPolicy-2014-01" {
  load_balancer_arn = aws_lb.test_lb.arn
  port              = "443"
  protocol          = "HTTPS"
  ssl_policy        = "ELBSecurityPolicy-2014-01"
  certificate_arn   = aws_acm_certificate.test_cert.arn

  default_action {
    type             = "forward"
    target_group_arn = aws_lb_target_group.test_lb_target_group.arn
  }
}

# Fail
resource "aws_alb_listener" "ssl_policy_set_to_ELBSecurityPolicy-2011-08" {
  load_balancer_arn = aws_lb.test_lb.arn
  port              = "443"
  protocol          = "HTTPS"
  ssl_policy        = "ELBSecurityPolicy-2011-08"
  certificate_arn   = aws_acm_certificate.test_cert.arn

  default_action {
    type             = "forward"
    target_group_arn = aws_lb_target_group.test_lb_target_group.arn
  }
}
