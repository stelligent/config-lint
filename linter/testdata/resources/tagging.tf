locals {
  TAGS_WITH_DATA_CLASS = {
    managed_by       = "Terraform Process"
    data_class       = "internal"
  }
  TAGS_WITH_INVALID_DATA_CLASS = {
    managed_by       = "Terraform Process"
    data_class       = "invalid"
  }
  TAGS_WITHOUT_DATA_CLASS = {
    managed_by       = "Terraform Process"
  }
}

# Pass
resource "aws_db_instance" "pass_main_db" {
  count                     = 1
  allocated_storage         = 100
  max_allocated_storage     = 150
  storage_type              = "gp2"
  engine                    = "mysql"
  engine_version            = "5.7"
  instance_class            = "db.t2.micro"
  name                      = "test-data_class-tag"
  username                  = "testUser"
  password                  = "temppw-test1234%"
  tags = {"data_class" = "internal"}

  storage_encrypted = true
}

# Fail
resource "aws_db_instance" "fail_main_db" {
  count                     = 1
  allocated_storage         = 100
  max_allocated_storage     = 150
  storage_type              = "gp2"
  engine                    = "mysql"
  engine_version            = "5.7"
  instance_class            = "db.t2.micro"
  name                      = "test-data_class-tag"
  username                  = "testUser"
  password                  = "temppw-test1234%"
  tags = {"data_class" = "somethingelse"}

  storage_encrypted = true
}

# Pass
resource "aws_db_instance" "pass_with_merge_main_db" {
  count                     = 1
  allocated_storage         = 100
  max_allocated_storage     = 150
  storage_type              = "gp2"
  engine                    = "mysql"
  engine_version            = "5.7"
  instance_class            = "db.t2.micro"
  name                      = "test-data_class-tag"
  username                  = "testUser"
  password                  = "temppw-test1234%"
  tags = merge(
    local.TAGS_WITHOUT_DATA_CLASS,
    {"data_class" = "internal"}
  )

  storage_encrypted = true
}

# Fail
resource "aws_db_instance" "missing_tags_main_db" {
  count                     = 1
  allocated_storage         = 100
  max_allocated_storage     = 150
  storage_type              = "gp2"
  engine                    = "mysql"
  engine_version            = "5.7"
  instance_class            = "db.t2.micro"
  name                      = "test-data_class-tag"
  username                  = "testUser"
  password                  = "temppw-test1234%"

  storage_encrypted = true
}

# Fail
resource "aws_db_instance" "missing_data_class_db" {
  count                     = 1
  allocated_storage         = 100
  max_allocated_storage     = 150
  storage_type              = "gp2"
  engine                    = "mysql"
  engine_version            = "5.7"
  instance_class            = "db.t2.micro"
  name                      = "test-data_class-tag"
  username                  = "testUser"
  password                  = "temppw-test1234%"
  tags = merge(
    local.TAGS_WITHOUT_DATA_CLASS,
    {"somethingelse" = "example"}
  )
  storage_encrypted = true
}

# Pass
resource "aws_db_instance" "inherit_data_class_db" {
  count                     = 1
  allocated_storage         = 100
  max_allocated_storage     = 150
  storage_type              = "gp2"
  engine                    = "mysql"
  engine_version            = "5.7"
  instance_class            = "db.t2.micro"
  name                      = "test-data_class-tag"
  username                  = "testUser"
  password                  = "temppw-test1234%"
  tags = merge(
    local.TAGS_WITH_DATA_CLASS,
    {"somethingelse" = "example"}
  )
  storage_encrypted = true
}

# Fail
resource "aws_db_instance" "inherit_invalid_data_class_db" {
  count                     = 1
  allocated_storage         = 100
  max_allocated_storage     = 150
  storage_type              = "gp2"
  engine                    = "mysql"
  engine_version            = "5.7"
  instance_class            = "db.t2.micro"
  name                      = "test-data_class-tag"
  username                  = "testUser"
  password                  = "temppw-test1234%"
  tags = merge(
    local.TAGS_WITH_INVALID_DATA_CLASS,
    {"somethingelse" = "example"}
  )
  storage_encrypted = true
}


# Fail
resource "aws_db_instance" "invalid_data_class_db" {
  count                     = 1
  allocated_storage         = 100
  max_allocated_storage     = 150
  storage_type              = "gp2"
  engine                    = "mysql"
  engine_version            = "5.7"
  instance_class            = "db.t2.micro"
  name                      = "test-data_class-tag"
  username                  = "testUser"
  password                  = "temppw-test1234%"
  tags = merge(
    local.TAGS_WITHOUT_DATA_CLASS,
    {"data_class" = "somethingelse"}
  )
  storage_encrypted = true
}

# Pass
resource "aws_db_instance" "inherit_data_class_db" {
  count                     = 1
  allocated_storage         = 100
  max_allocated_storage     = 150
  storage_type              = "gp2"
  engine                    = "mysql"
  engine_version            = "5.7"
  instance_class            = "db.t2.micro"
  name                      = "test-data_class-tag"
  username                  = "testUser"
  password                  = "temppw-test1234%"
  tags = merge(
    local.TAGS_WITH_DATA_CLASS,
    {"somethingelse" = "example"}
  )
  storage_encrypted = true
}
