resource "aws_security_group" "alb_sg" {
  vpc_id = "${aws_vpc.website.id}"
  name = "load balancer"
  ingress {
    protocol = "tcp"
    from_port = 80
    to_port = 80
    cidr_blocks = [ "0.0.0.0/0" ]
  }
  egress {
    protocol = "tcp"
    from_port = 0
    to_port = 65535
    cidr_blocks = [ "0.0.0.0/0" ]
  }
}

resource "aws_alb" "alb" {
  name = "web-alb"
  internal = false
  load_balancer_type = "application"
  security_groups = [ "${aws_security_group.alb_sg.id}" ]
  subnets = [
    "${aws_subnet.public1.id}",
    "${aws_subnet.public2.id}"
  ]
  access_logs {
    bucket = "${aws_s3_bucket.alb_logs.bucket}"
    enabled = true
  }
}

resource "aws_s3_bucket" "alb_logs" {}

resource "aws_s3_bucket_policy" "bucket_policy" {
    bucket = "${aws_s3_bucket.alb_logs.bucket}"
    policy = <<POLICY
{
  "Version": "2012-10-17",
  "Id": "MYBUCKETPOLICY",
  "Statement": [
    {
      "Sid": "AllowAccessLogs",
      "Effect": "Allow",
      "Principal": {
        "AWS": "929392832123"
      },
      "Action": "s3:PutObject",
      "Resource": [
        "${aws_s3_bucket.alb_logs.arn}",
        "${aws_s3_bucket.alb_logs.arn}/*"
      ]
    } 
  ]
}
POLICY
}

resource "aws_security_group" "ec2_sg" {
  vpc_id = "${aws_vpc.website.id}"
  name = "instance"
  ingress {
    protocol = "tcp"
    from_port = 22
    to_port = 22
    cidr_blocks = [ "${var.ssh_cidr_block}" ]
  }
  ingress {
    protocol = "tcp"
    from_port = 80
    to_port = 80
    security_groups = [ "${aws_security_group.alb_sg.id}" ]
  }
  egress {
    protocol = "tcp"
    from_port = 0
    to_port = 65535
    cidr_blocks = [ "0.0.0.0/0" ]
  }
}

data "template_file" "user_data" {
  template = "${file("user_data.tpl")}"
  vars {
    title = "${var.title}"
  }
}

resource "aws_launch_configuration" "webserver" {
  image_id = "${var.ami}"
  instance_type = "t2.micro"
  lifecycle {
    create_before_destroy = true
  }
  key_name = "${var.key_name}"
  associate_public_ip_address = true
  security_groups = [ "${aws_security_group.ec2_sg.id}" ]
  user_data = "${data.template_file.user_data.rendered}"
}

resource "aws_autoscaling_group" "asg" {
  name = "terraform-test"
  min_size = 1
  max_size = 2
  health_check_grace_period = 300
  health_check_type = "ELB"
  desired_capacity = 2
  launch_configuration = "${aws_launch_configuration.webserver.name}"
  vpc_zone_identifier = [ "${aws_subnet.public1.id}", "${aws_subnet.public2.id}" ]
  target_group_arns = [ "${aws_alb_target_group.tg.arn}" ]
  tag {
    key = "Name"
    value = "${var.name_tag}"
    propagate_at_launch = true
  }
}

resource "aws_alb_target_group" "tg" {
  name = "webserver-target"
  port = 80
  protocol = "HTTPS"
  vpc_id = "${aws_vpc.website.id}"
}

resource "aws_alb_listener" "l" {
  load_balancer_arn = "${aws_alb.alb.id}"
  port = 80
  protocol = "HTTPS"
  default_action {
    type = "forward"
    target_group_arn = "${aws_alb_target_group.tg.arn}"
  }
}

