# Pass
resource "aws_ecs_task_definition" "container_definitions_environment_not_set" {
  family                = "foo"
  container_definitions = <<EOF
[
  {
    "name": "bar",
    "image": "foobar",
    "cpu": 10,
    "memory": 512,
    "essential": true,
    "portMappings": [
      {
        "containerPort": 80,
        "hostPort": 80
      }
    ]
  }
]
EOF
}

# Pass
resource "aws_ecs_task_definition" "container_definitions_environment_aws_secrets_not_set" {
  family                = "foo"
  container_definitions = <<EOF
[
  {
    "name": "bar",
    "image": "foobar",
    "cpu": 10,
    "memory": 512,
    "essential": true,
    "portMappings": [
      {
        "containerPort": 80,
        "hostPort": 80
      }
    ],
    "environment": [
        {
            "name": "foo",
            "value": "bar"
        }
    ]
  }
]
EOF
}

# Pass
resource "aws_ecs_task_definition" "container_definitions_environment_aws_secrets_not_set_20_character_capital_string" {
  family                = "foo"
  container_definitions = <<EOF
[
  {
    "name": "bar",
    "image": "foobar",
    "cpu": 10,
    "memory": 512,
    "essential": true,
    "portMappings": [
      {
        "containerPort": 80,
        "hostPort": 80
      }
    ],
    "environment": [
        {
            "name": "foo",
            "value": "AXYZIOSFODNN7EXAMPLE"
        }
    ]
  }
]
EOF
}

# Pass
resource "aws_ecs_task_definition" "container_definitions_environment_aws_secrets_not_set_21_character_capital_string" {
  family                = "foo"
  container_definitions = <<EOF
[
  {
    "name": "bar",
    "image": "foobar",
    "cpu": 10,
    "memory": 512,
    "essential": true,
    "portMappings": [
      {
        "containerPort": 80,
        "hostPort": 80
      }
    ],
    "environment": [
        {
            "name": "foo",
            "value": "AKIAIOSFODNN7FEXAMPLE"
        }
    ]
  }
]
EOF
}

# Pass
resource "aws_ecs_task_definition" "container_definitions_environment_aws_secrets_not_set_40_character_string" {
  family                = "foo"
  container_definitions = <<EOF
[
  {
    "name": "bar",
    "image": "foobar",
    "cpu": 10,
    "memory": 512,
    "essential": true,
    "portMappings": [
      {
        "containerPort": 80,
        "hostPort": 80
      }
    ],
    "environment": [
        {
            "name": "foo",
            "value": "wJalrXUtnFEMI>K7MDENG^bPxRfiCYEXAMPLEKEY"
        }
    ]
  }
]
EOF
}

# Pass
resource "aws_ecs_task_definition" "container_definitions_environment_aws_secrets_not_set_41_character_string" {
  family                = "foo"
  container_definitions = <<EOF
[
  {
    "name": "bar",
    "image": "foobar",
    "cpu": 10,
    "memory": 512,
    "essential": true,
    "portMappings": [
      {
        "containerPort": 80,
        "hostPort": 80
      }
    ],
    "environment": [
        {
            "name": "foo",
            "value": "wJalrXUtnFOOMI/K7MDENG/bPxRfiCYEXAMPLEKEY"
        }
    ]
  }
]
EOF
}

# Fail
resource "aws_ecs_task_definition" "container_definitions_environment_aws_access_key_set" {
  family                = "foo"
  container_definitions = <<EOF
[
  {
    "name": "bar",
    "image": "foobar",
    "cpu": 10,
    "memory": 512,
    "essential": true,
    "portMappings": [
      {
        "containerPort": 80,
        "hostPort": 80
      }
    ],
    "environment": [
        {
            "name": "foo",
            "value": "AKIAIOSFODNN7EXAMPLE"
        }
    ]
  }
]
EOF
}

# Fail
resource "aws_ecs_task_definition" "container_definitions_environment_aws_secret_access_key_set" {
  family                = "foo"
  container_definitions = <<EOF
[
  {
    "name": "bar",
    "image": "foobar",
    "cpu": 10,
    "memory": 512,
    "essential": true,
    "portMappings": [
      {
        "containerPort": 80,
        "hostPort": 80
      }
    ],
    "environment": [
        {
            "name": "foo",
            "value": "wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY"
        }
    ]
  }
]
EOF
}

# Fail
resource "aws_ecs_task_definition" "container_definitions_environment_aws_access_key_and_secret_access_key_set" {
  family                = "foo"
  container_definitions = <<EOF
[
  {
    "name": "bar",
    "image": "foobar",
    "cpu": 10,
    "memory": 512,
    "essential": true,
    "portMappings": [
      {
        "containerPort": 80,
        "hostPort": 80
      }
    ],
    "environment": [
        {
            "name": "foo",
            "value": "AIPAIOSFODNN7EXAMPLE"
        },
        {
            "name": "bar",
            "value": "wJalrXUtnFEMI/K7MDENG+bPxRfiCYEXAMPLEKEY"
        }
    ]
  }
]
EOF
}
