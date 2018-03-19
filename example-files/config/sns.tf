resource "aws_sns_topic" "test" {
  name = "my-topic-with-policy"
}

resource "aws_sns_topic_policy" "default" {
  arn = "${aws_sns_topic.test.arn}"

  policy =<<POLICY
{
  "Version": "2012-10-17",
  "Id": "MYSNSTOPICPOLICY",
  "Statement": [
    {
      "Effect": "Allow",
      "Principal": "*",
      "Action": "*"
    } 
  ]
}
POLICY
}

resource "aws_sns_topic_policy" "sns_topic_policy_with_not" {
  arn = "${aws_sns_topic.test.arn}"

  policy =<<POLICY
{
  "Version": "2012-10-17",
  "Id": "MYSNSTOPICPOLICY",
  "Statement": [
    {
      "Effect": "Allow",
      "NotPrincipal": "*",
      "NotAction": "*"
    } 
  ]
}
POLICY
}

