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
resource "aws_alb_listener" "listener_secure_https_set" {
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

# # Pass
# resource "aws_alb_listener" "listener_secure_https_set_lowercase" {
#   load_balancer_arn = aws_lb.test_lb.arn
#   port              = "443"
#   protocol          = "https"
#   ssl_policy        = "ELBSecurityPolicy-2016-08"
#   certificate_arn   = aws_acm_certificate.test_cert.arn

#   default_action {
#     type             = "forward"
#     target_group_arn = aws_lb_target_group.test_lb_target_group.arn
#   }
# }

# # Fail
# resource "aws_alb_listener" "port_set_to_80" {
#   load_balancer_arn = aws_lb.test_lb.arn
#   port              = "80"
#   protocol          = "HTTPS"
#   ssl_policy        = "ELBSecurityPolicy-2016-08"
#   certificate_arn   = aws_acm_certificate.test_cert.arn

#   default_action {
#     type             = "forward"
#     target_group_arn = aws_lb_target_group.test_lb_target_group.arn
#   }
# }

# # Fail
# resource "aws_alb_listener" "protocol_set_to_http" {
#   load_balancer_arn = aws_lb.test_lb.arn
#   port              = "443"
#   protocol          = "HTTP"
#   ssl_policy        = "ELBSecurityPolicy-2016-08"
#   certificate_arn   = aws_acm_certificate.test_cert.arn

#   default_action {
#     type             = "forward"
#     target_group_arn = aws_lb_target_group.test_lb_target_group.arn
#   }
# }

# # Fail
# resource "aws_alb_listener" "ssl_policy_not_set" {
#   load_balancer_arn = aws_lb.test_lb.arn
#   port              = "443"
#   protocol          = "HTTPS"
#   certificate_arn   = aws_acm_certificate.test_cert.arn

#   default_action {
#     type             = "forward"
#     target_group_arn = aws_lb_target_group.test_lb_target_group.arn
#   }
# }

# # Fail
# resource "aws_alb_listener" "certificate_arn_not_set" {
#   load_balancer_arn = aws_lb.test_lb.arn
#   port              = "443"
#   protocol          = "HTTPS"
#   ssl_policy        = "ELBSecurityPolicy-2016-08"

#   default_action {
#     type             = "forward"
#     target_group_arn = aws_lb_target_group.test_lb_target_group.arn
#   }
# }
