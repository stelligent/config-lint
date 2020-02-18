# Pass
resource "aws_security_group" "egress_port_range_matches" {
  name        = "allow_http"
  description = "Allow HTTP traffic"
  egress {
    from_port = 80
    to_port   = 80
    protocol  = "tcp"
  }
}

# Warn
resource "aws_security_group" "egress_port_range_does_not_match" {
  name        = "allow_foo"
  description = "Allow FOO traffic"
  egress {
    from_port = 10000
    to_port   = 10200
    protocol  = "tcp"
  }
}
