data "template_file" "user_data" {
  template = "${file("user_data.tpl")}"
  vars {
    title = "${var.title}"
  }
}

resource "aws_instance" "web" {
    availability_zone = "${var.availability_zone}"
    ami = "${var.ami}"
    instance_type = "t2.micro"
    security_groups = [ "${aws_security_group.sg1.name}" ]
    user_data = "${data.template_file.user_data.rendered}"
    tags {
        name = "website"
    }
    tags {
        environment = "development"
    }
}

resource "aws_ebs_volume" "volume1" {
    availability_zone = "${var.availability_zone}"
    size = 10
}

resource "aws_security_group" "sg1" {
  name        = "allow_http"
  description = "Allow HTTP traffic"

  ingress {
    from_port   = 80
    to_port     = 80
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
