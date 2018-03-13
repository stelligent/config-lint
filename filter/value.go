package filter

import (
	"bytes"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

func getValue(filter Filter) string {
	if filter.ValueFrom.Bucket != "" {
		content, err := getValueFromS3(filter.ValueFrom.Bucket, filter.ValueFrom.Key)
		if err != nil {
			return "Error" // FIXME
		}
		return content
	}
	return filter.Value
}

func getValueFromS3(bucket string, key string) (string, error) {

	region := &aws.Config{Region: aws.String("us-east-1")}
	awsSession := session.New()
	s3Client := s3.New(awsSession, region)
	response, err := s3Client.GetObject(&s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		return "", err
	}
	buf := new(bytes.Buffer)
	buf.ReadFrom(response.Body)
	return buf.String(), nil
}
