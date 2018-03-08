resource "aws_ebs_volume" "vol1" {
    availability_zone = "us-west-2a"
    size = 40
    tags {
        Name = "HelloWorld"
    }
    encrypted = true
}
resource "aws_ebs_volume" "vol2" {
    availability_zone = "us-west-2a"
    size = 40
    tags {
        Name = "HelloWorld"
    }
    encrypted = false
}
