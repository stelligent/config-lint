## Setup Helper
resource "aws_iam_group" "test_group" {
  name = "test-group"
}

resource "aws_iam_user" "test_user" {
  name = "test-user"
}

# Pass
resource "aws_iam_group_membership" "group_and_users_set" {
  name = "tf-testing-group-membership"

  group = "${aws_iam_group.test_group.name}"

  users = [
    "${aws_iam_user.test_user.name}"
  ]
}

# Fail
resource "aws_iam_group_membership" "group_set_and_users_empty" {
  name = "tf-testing-group-membership"

  users = []

  group = "${aws_iam_group.test_group.name}"
}

# Fail
resource "aws_iam_group_membership" "group_empty_and_users_set" {
  name = "tf-testing-group-membership"

  users = [
    "${aws_iam_user.test_user.name}"
  ]

  group = ""
}

# Fail x 2
resource "aws_iam_group_membership" "group_empty_and_users_empty" {
  name = "tf-testing-group-membership"

  users = []

  group = ""
}
