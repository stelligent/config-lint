# Pass
resource "aws_security_group" "egress_block_exists" {
  name        = "allow_http"
  description = "Allow HTTP traffic"
  egress {
    from_port = 80
    to_port   = 80
    protocol  = "tcp"
  }
}

# Warn
resource "aws_security_group" "egress_block_missing" {
  name        = "allow_http"
  description = "Allow HTTP traffic"
  ingress {
    from_port = 80
    to_port   = 80
    protocol  = "tcp"
  }
}
