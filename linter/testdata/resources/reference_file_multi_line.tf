terraform {
  required_version = ">= 0.12.0"
}

resource "test_resource" "test" {
  test_value = file("./testdata/data/multi_line_content")
}

resource "test_resource2" "test2" {
  test_value2 = <<EOT
multi
line
example
EOT
}
