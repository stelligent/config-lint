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

# Pass
resource "aws_security_group" "ingress_cidr_blocks_set_to_ip_with_32_subnet" {
  name        = "allow_http"
  description = "Allow HTTP traffic"
  ingress {
    from_port   = 80
    to_port     = 80
    protocol    = "tcp"
    cidr_blocks = ["10.0.0.100/32"]
  }
}

# Warn
resource "aws_security_group" "ingress_cidr_blocks_set_to_ip_with_24_subnet" {
  name        = "allow_http"
  description = "Allow HTTP traffic"
  ingress {
    from_port   = 80
    to_port     = 80
    protocol    = "tcp"
    cidr_blocks = ["10.0.0.0/24"]
  }
}

# Warn
resource "aws_security_group" "ingress_cidr_blocks_set_to_ip_with_multiple_subnets" {
  name        = "allow_http"
  description = "Allow HTTP traffic"
  ingress {
    from_port = 80
    to_port   = 80
    protocol  = "tcp"
    cidr_blocks = [
      "10.0.0.100/32",
      "10.0.0.110/32",
      "10.0.0.120/32",
      "10.0.0.130/32",
      "10.1.0.0/24"
    ]
  }
}
