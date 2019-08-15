resource "aws_security_group" "sg_ok" {
  name        = "allow_http"
  description = "Allow HTTP traffic"
  ingress {
    from_port   = 80
    to_port     = 80
    protocol    = "tcp"
    cidr_blocks = ["1.2.3.4/32"]
  }
  egress {
    from_port       = 0
    to_port         = 0
    protocol        = "-1"
    cidr_blocks     = ["0.0.0.0/0"]
  }
}

resource "aws_security_group" "sg_ssh" {
  name        = "allow_ssh"
  description = "Allow SSH traffic"
  ingress {
    from_port   = 22
    to_port     = 22
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }
  egress {
    from_port       = 0
    to_port         = 0
    protocol        = "-1"
    cidr_blocks     = ["0.0.0.0/0"]
  }
}

resource "aws_security_group" "sg_all_protocols" {
  name = "all_all_protocols"
  description = "Allow all protocols and ports"

  ingress {
    protocol    = "-1"
    cidr_blocks = ["1.2.3.4/32"]
  }
  egress {
    protocol    = "-1"
    cidr_blocks = ["1.2.3.4/32"]
  }
}
