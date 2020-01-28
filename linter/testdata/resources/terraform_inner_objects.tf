data "aws_kms_alias" "s3kmskey" {
  name = "alias/myKmsKey"
}

resource "aws_codepipeline" "foo" {
  artifact_store = {
    encryption_key = {
      id   = "id"
      type = "KMS"
    }
  }
}