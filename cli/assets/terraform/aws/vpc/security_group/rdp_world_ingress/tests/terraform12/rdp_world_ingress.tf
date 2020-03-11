# Pass
resource "aws_security_group" "ingress_cidr_blocks_not_set" {
  name        = "allow_rdp"
  description = "Allow RDP traffic"
  ingress {
    from_port = 3389
    to_port   = 3389
    protocol  = "tcp"
  }
}

# Pass
resource "aws_security_group" "ingress_cidr_blocks_set_to_ip_and_rdp_not_set" {
  name        = "allow_http"
  description = "Allow HTTP traffic"
  ingress {
    from_port   = 80
    to_port     = 80
    protocol    = "tcp"
    cidr_blocks = ["10.0.0.100/32"]
  }
}

# Pass
resource "aws_security_group" "ingress_cidr_blocks_set_to_world_and_rdp_not_set" {
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
resource "aws_security_group" "ingress_cidr_blocks_set_to_ip_and_rdp_is_set" {
  name        = "allow_rdp"
  description = "Allow RDP traffic"
  ingress {
    from_port   = 3389
    to_port     = 3389
    protocol    = "tcp"
    cidr_blocks = ["10.0.0.100/32"]
  }
}

# Fail
resource "aws_security_group" "ingress_cidr_blocks_set_to_world_and_rdp_is_set" {
  name        = "allow_rdp"
  description = "Allow RDP traffic"
  ingress {
    from_port   = 3389
    to_port     = 3389
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }
}

# Pass
resource "aws_security_group" "ingress_ipv6_cidr_blocks_set_to_ip_and_rdp_not_set" {
  name        = "allow_http"
  description = "Allow HTTP traffic"
  ingress {
    from_port        = 80
    to_port          = 80
    protocol         = "tcp"
    ipv6_cidr_blocks = ["0:0:0:0:0:ffff:a00:64/32"]
  }
}

# Pass
resource "aws_security_group" "ingress_ipv6_cidr_blocks_set_to_world_and_rdp_not_set" {
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
resource "aws_security_group" "ingress_ipv6_cidr_blocks_set_to_ip_and_rdp_is_set" {
  name        = "allow_rdp"
  description = "Allow RDP traffic"
  ingress {
    from_port        = 3389
    to_port          = 3389
    protocol         = "tcp"
    ipv6_cidr_blocks = ["0:0:0:0:0:ffff:a00:64/32"]
  }
}

# Fail
resource "aws_security_group" "ingress_ipv6_cidr_blocks_set_to_world_and_rdp_is_set" {
  name        = "allow_rdp"
  description = "Allow RDP traffic"
  ingress {
    from_port        = 3389
    to_port          = 3389
    protocol         = "tcp"
    ipv6_cidr_blocks = ["::/0"]
  }
}
