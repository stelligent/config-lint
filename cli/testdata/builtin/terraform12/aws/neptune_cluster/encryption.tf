# Pass
resource "aws_neptune_cluster" "storage_encrypted_set_to_true" {
  storage_encrypted = true
}

# Fail
resource "aws_neptune_cluster" "storage_encrypted_set_to_false" {
  storage_encrypted = false
}

# Fail
resource "aws_neptune_cluster" "storage_encrypted_not_set" {
}
