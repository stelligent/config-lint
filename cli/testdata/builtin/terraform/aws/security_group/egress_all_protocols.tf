# Pass
resource "aws_security_group" "egress_tcp_protocol" {
  name        = "allow_tcp"
  description = "Allow TCP traffic"
  egress {
    from_port = 10000
    to_port   = 10000
    protocol  = "tcp"
  }
}

# Pass
resource "aws_security_group" "egress_udp_protocol" {
  name        = "allow_udp"
  description = "Allow UDP traffic"
  egress {
    from_port = 10000
    to_port   = 10000
    protocol  = "udp"
  }
}

# Pass
resource "aws_security_group" "egress_icmp_protocol" {
  name        = "allow_icmp"
  description = "Allow ICMP traffic"
  egress {
    from_port = 10000
    to_port   = 10000
    protocol  = "icmp"
  }
}

# Warn
resource "aws_security_group" "egress_all_protocols" {
  name        = "allow_all"
  description = "Allow ALL traffic"
  egress {
    from_port = 0
    to_port   = 0
    protocol  = "-1"
  }
}
