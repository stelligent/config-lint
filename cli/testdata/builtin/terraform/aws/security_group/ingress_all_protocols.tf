# Pass
resource "aws_security_group" "ingress_tcp_protocol" {
  name        = "allow_tcp"
  description = "Allow TCP traffic"
  ingress {
    from_port = 10000
    to_port   = 10000
    protocol  = "tcp"
  }
}

# Pass
resource "aws_security_group" "ingress_udp_protocol" {
  name        = "allow_udp"
  description = "Allow UDP traffic"
  ingress {
    from_port = 10000
    to_port   = 10000
    protocol  = "udp"
  }
}

# Pass
resource "aws_security_group" "ingress_icmp_protocol" {
  name        = "allow_icmp"
  description = "Allow ICMP traffic"
  ingress {
    from_port = 10000
    to_port   = 10000
    protocol  = "icmp"
  }
}

# Warn
resource "aws_security_group" "ingress_all_protocols" {
  name        = "allow_all"
  description = "Allow ALL traffic"
  ingress {
    from_port = 0
    to_port   = 0
    protocol  = "-1"
  }
}
