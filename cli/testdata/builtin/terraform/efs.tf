variable "aws_kms_alias" {}

resource "aws_efs_file_system" "fs" {
  creation_token = "my-efs"
  encrypted = true
  kms_key_id = "${var.aws_kms_alias}"
  tags {
    Name = "MyProduct"
  }
}
