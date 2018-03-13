package filter

import (
	"bytes"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

type StandardValueSource struct {
	Log LoggingFunction
}

func (v StandardValueSource) GetValue(filter Filter) string {
	if filter.ValueFrom.Bucket != "" {
		v.Log(fmt.Sprintf("Getting value_from s3://%s/%s", filter.ValueFrom.Bucket, filter.ValueFrom.Key))
		content, err := v.GetValueFromS3(filter.ValueFrom.Bucket, filter.ValueFrom.Key)
		if err != nil {
			return "Error" // FIXME
		}
		return content
	}
	return filter.Value
}

func (v StandardValueSource) GetValueFromS3(bucket string, key string) (string, error) {
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
