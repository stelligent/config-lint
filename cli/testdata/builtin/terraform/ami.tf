resource "aws_ami" "my-ami" {
    name = "nginx-example"
    virtualization_type = "hvm"
    root_device_name = "/dev/xvda"
    ebs_block_device {
        device_name = "/dev/xvda"
        snapshot_id = "${var.snapshot_id}"
        volume_size = 8
    }
}

resource "aws_ami_copy" "my-ami-copy" {
  name              = "nginx-example-copy"
  description       = "A copy of nginx-example"
  source_ami_id     = "${aws_ami.my-ami.id}"
  source_ami_region = "us-east-1"
}
