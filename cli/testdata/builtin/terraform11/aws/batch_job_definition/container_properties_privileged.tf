# Pass
resource "aws_batch_job_definition" "container_properties_privileged_not_set" {
  name = "foo"
  type = "container"

  container_properties = <<EOF
{
    "command": ["ls", "-la"],
    "image": "busybox",
    "memory": 1024,
    "vcpus": 1,
    "volumes": [
      {
        "host": {
          "sourcePath": "/tmp"
        },
        "name": "tmp"
      }
    ],
    "environment": [
        {"name": "VARNAME", "value": "VARVAL"}
    ],
    "mountPoints": [
        {
          "sourceVolume": "tmp",
          "containerPath": "/tmp",
          "readOnly": false
        }
    ],
    "ulimits": [
      {
        "hardLimit": 1024,
        "name": "nofile",
        "softLimit": 1024
      }
    ]
}
EOF
}

# Pass
resource "aws_batch_job_definition" "container_properties_privileged_set_to_false" {
  name = "foo"
  type = "container"

  container_properties = <<EOF
{
    "command": ["ls", "-la"],
    "image": "busybox",
    "memory": 1024,
    "vcpus": 1,
    "privileged": "false",
    "volumes": [
      {
        "host": {
          "sourcePath": "/tmp"
        },
        "name": "tmp"
      }
    ],
    "environment": [
        {"name": "VARNAME", "value": "VARVAL"}
    ],
    "mountPoints": [
        {
          "sourceVolume": "tmp",
          "containerPath": "/tmp",
          "readOnly": false
        }
    ],
    "ulimits": [
      {
        "hardLimit": 1024,
        "name": "nofile",
        "softLimit": 1024
      }
    ]
}
EOF
}

# Warn
resource "aws_batch_job_definition" "container_properties_privileged_set_to_true" {
  name = "foo"
  type = "container"

  container_properties = <<EOF
{
    "command": ["ls", "-la"],
    "image": "busybox",
    "memory": 1024,
    "vcpus": 1,
    "privileged": "true",
    "volumes": [
      {
        "host": {
          "sourcePath": "/tmp"
        },
        "name": "tmp"
      }
    ],
    "environment": [
        {"name": "VARNAME", "value": "VARVAL"}
    ],
    "mountPoints": [
        {
          "sourceVolume": "tmp",
          "containerPath": "/tmp",
          "readOnly": false
        }
    ],
    "ulimits": [
      {
        "hardLimit": 1024,
        "name": "nofile",
        "softLimit": 1024
      }
    ]
}
EOF
}
