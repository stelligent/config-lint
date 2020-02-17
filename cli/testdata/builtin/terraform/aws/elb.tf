resource "aws_instance" "web" {
    availability_zone = "us-west-2a"
    ami = "ami-0fcd5791ba781e98f"
    instance_type = "t2.micro"
    security_groups = [ "${aws_security_group.ec2_sg.name}" ]
    user_data = <<USER_DATA
#!/bin/bash
yum update -y
amazon-linux-extras install -y nginx1.12
service nginx start
USER_DATA
}

resource "aws_security_group" "ec2_sg" {
  name = "web"
  ingress {
    protocol = "tcp"
    from_port = 80
    to_port = 80
    cidr_blocks = [ "0.0.0.0/0" ]
    # can classic elb have a sg?
  }
  egress {
    protocol = "tcp"
    from_port = 0
    to_port = 65535
    cidr_blocks = [ "0.0.0.0/0" ]
  }
}

resource "aws_elb" "elb1" {
  name               = "test-terraform-elb"
  availability_zones = ["us-west-2a", "us-west-2b", "us-west-2c"]

  listener {
    instance_port     = 80
    instance_protocol = "http"
    lb_port           = 80
    lb_protocol       = "http"
  }

  health_check {
    healthy_threshold   = 2
    unhealthy_threshold = 2
    timeout             = 3
    target              = "HTTP:80/"
    interval            = 30
  }

  instances                   = ["${aws_instance.web.id}"]
  cross_zone_load_balancing   = true
  idle_timeout                = 400
  connection_draining         = true
  connection_draining_timeout = 400

}
