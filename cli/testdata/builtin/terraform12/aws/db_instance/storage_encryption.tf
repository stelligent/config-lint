## Setup Helper
variable "test_db_username" {
  default = "foo"
}

variable "test_db_password" {
  default = "foobarbaz"
}

resource "aws_kms_key" "test_key" {
  enable_key_rotation = true
}

# Pass
resource "aws_db_instance" "storage_encrypted_set_to_true" {
  allocated_storage = 20
  engine            = "mysql"
  engine_version    = "5.7"
  instance_class    = "db.t2.micro"
  username          = var.test_db_username
  password          = var.test_db_password
  storage_encrypted = true
}

# Fail
resource "aws_db_instance" "storage_encrypted_set_to_false" {
  allocated_storage = 20
  engine            = "mysql"
  engine_version    = "5.7"
  instance_class    = "db.t2.micro"
  username          = var.test_db_username
  password          = var.test_db_password
  storage_encrypted = false
}

# Fail
resource "aws_db_instance" "storage_encrypted_not_set" {
  allocated_storage = 20
  engine            = "mysql"
  engine_version    = "5.7"
  instance_class    = "db.t2.micro"
  username          = var.test_db_username
  password          = var.test_db_password
}

# Pass
resource "aws_db_instance" "replicate_source_db_is_set_with_kms_key_id" {
  instance_class      = "db.t2.micro"
  replicate_source_db = "foo"
  kms_key_id          = aws_kms_key.test_key.id
}

# Warn
resource "aws_db_instance" "replicate_source_db_is_set_without_kms_key_id" {
  instance_class      = "db.t2.micro"
  replicate_source_db = "foo"
}
