## Setup Helper
variable "test_db_username" {
  default = "foo"
}

variable "test_db_password" {
  default = "foobarbaz"
}

# Pass
resource "aws_db_instance" "publicly_accessible_not_set" {
  allocated_storage = 20
  engine            = "mysql"
  engine_version    = "5.7"
  instance_class    = "db.t2.micro"
  username          = var.test_db_username
  password          = var.test_db_password
  storage_encrypted = true
}

# Pass
resource "aws_db_instance" "publicly_accessible_set_to_false" {
  allocated_storage   = 20
  engine              = "mysql"
  engine_version      = "5.7"
  instance_class      = "db.t2.micro"
  username            = var.test_db_username
  password            = var.test_db_password
  storage_encrypted   = true
  publicly_accessible = false
}

# Fail
resource "aws_db_instance" "publicly_accessible_set_to_true" {
  allocated_storage   = 20
  engine              = "mysql"
  engine_version      = "5.7"
  instance_class      = "db.t2.micro"
  username            = var.test_db_username
  password            = var.test_db_password
  storage_encrypted   = true
  publicly_accessible = true
}
