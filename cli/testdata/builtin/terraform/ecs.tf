resource "aws_ecs_task_definition" "task1" {
  family = "application"

  container_definitions = <<DEFINITION
[
  {
    "cpu": 128,
    "environment": [
        {
            "name": "AWS_ACCESS_KEY_ID",
            "value": "AKIAIOSFODNN7EXAMPLE"
        },
        {
            "name": "AWS_SECRET_ACCESS_KEY",
            "value": "wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY"
        }
    ],
    "essential": true,
    "image": "application:latest",
    "memory": 128,
    "memoryReservation": 64,
    "name": "application"
  }
]
DEFINITION
}
