resource "aws_sqs_queue" "q1" {}

resource "aws_sqs_queue_policy" "policy1" {
    queue_url = "${aws_sqs_queue.q1.id}"
    policy = <<POLICY
{   
   "Version": "2012-10-17",
   "Id": "IPAllow",
   "Statement" : [{
      "Sid": "1", 
      "Effect": "Allow",           
      "Principal": {
         "AWS": [
            "*"
         ]
      },
      "Action": [
         "sqs:SendMessage",
         "sqs:ReceiveMessage"
      ], 
      "Resource": "${aws_sqs_queue.q1.arn}",
      "Condition": {
         "IpAddress": {"aws:SourceIp": "10.10.1.10/32"}
      } 
   }]
}
POLICY
}
