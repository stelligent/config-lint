variable "service_ports" {
  default = [22, 80, 1433, 6379]
}

resource "aws_security_group" "example" {
  name = "example"

  dynamic "ingress" {
    for_each = var.service_ports
    content {
      from_port = ingress.value
      to_port   = ingress.value
      protocol  = "tcp"
    }
  }

  egress {
    from_port = 443
    to_port = 443
    protocol = "tcp"
  }
}