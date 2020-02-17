# Pass
resource "aws_security_group" "ingress_cidr_blocks_not_set" {
  name        = "allow_http"
  description = "Allow HTTP traffic"
  ingress {
    from_port = 80
    to_port   = 80
    protocol  = "tcp"
  }
}

# Warn
resource "aws_security_group" "ingress_cidr_blocks_set_to_world" {
  name        = "allow_http"
  description = "Allow HTTP traffic"
  ingress {
    from_port   = 80
    to_port     = 80
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }
}

# Pass
resource "aws_security_group" "ingress_cidr_blocks_set_to_ip" {
  name        = "allow_http"
  description = "Allow HTTP traffic"
  ingress {
    from_port   = 80
    to_port     = 80
    protocol    = "tcp"
    cidr_blocks = ["1.2.3.4/32"]
  }
}

# Pass
resource "aws_security_group" "ingress_ipv6_cidr_blocks_not_set" {
  name        = "allow_http"
  description = "Allow HTTP traffic"
  ingress {
    from_port = 80
    to_port   = 80
    protocol  = "tcp"
  }
}

# Warn
resource "aws_security_group" "ingress_ipv6_cidr_blocks_set_to_world" {
  name        = "allow_http"
  description = "Allow HTTP traffic"
  ingress {
    from_port        = 80
    to_port          = 80
    protocol         = "tcp"
    ipv6_cidr_blocks = ["::/0"]
  }
}

# Pass
resource "aws_security_group" "ingress_ipv6_cidr_blocks_set_to_ip" {
  name        = "allow_http"
  description = "Allow HTTP traffic"
  ingress {
    from_port        = 80
    to_port          = 80
    protocol         = "tcp"
    ipv6_cidr_blocks = ["0:0:0:0:0:ffff:102:304/32"]
  }
}
