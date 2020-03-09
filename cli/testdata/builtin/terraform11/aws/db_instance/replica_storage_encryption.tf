## Setup Helper
provider "aws" {
  alias  = "west"
  region = "us-west-1"
}

resource "aws_kms_key" "test_key" {
  enable_key_rotation = true
}

# Pass
resource "aws_db_instance" "replicate_source_db_not_set" {
  allocated_storage = 20
  engine            = "mysql"
  engine_version    = "5.7"
  instance_class    = "db.t2.micro"
  username          = "foo"
  password          = "foobarbaz"
  storage_encrypted = true
}

# Pass
resource "aws_db_instance" "replicate_source_db_with_kms_key_id_set" {
  instance_class      = "db.t2.micro"
  replicate_source_db = "foo"
  kms_key_id          = "${aws_kms_key.test_key.id}"
}

# Warn
resource "aws_db_instance" "replicate_source_db_without_kms_key_id_set" {
  instance_class      = "db.t2.micro"
  replicate_source_db = "foo"
}
