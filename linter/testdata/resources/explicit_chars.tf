# Fail
resource "aws_s3_bucket_policy" "b" {
  bucket = "${aws_s3_bucket.b.id}"
  policy =<<POLICY
{
    "Version": "2018-08-09",
    "Statement": [
    {
        "Effect": "Deny",
        "Action": "s3:*",
        "Principal": {"AWS": [
            "*"
        ]},
        "Resource": [
            "arn:aws:s3:::BUCKETNAME",
            "arn:aws:s3:::BUCKETNAME/*"
        ],
        "Condition": { "Bool": { "aws:SecureTransport": "true" } }
    }]
}
POLICY
}

# Pass
resource "aws_s3_bucket_policy" "b" {
  bucket = "${aws_s3_bucket.b.id}"
  policy =<<POLICY
{
    "Version": "2018-08-09",
    "Statement": [
    {
        "Effect": "Deny",
        "Action": "s3:*",
        "Principal": {"AWS": [
            "*"
        ]},
        "Resource": [
            "arn:aws:s3:::BUCKETNAME",
            "arn:aws:s3:::BUCKETNAME/*"
        ],
        "Condition": { "Bool": { "aws:SecureTransport": "false" } }
    }]
}
POLICY
}
