package assertion

import (
	"bytes"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"net/url"
)

type StandardValueSource struct {
	Log LoggingFunction
}

func (v StandardValueSource) GetValue(assertion Assertion) string {
	if assertion.ValueFrom.Url != "" {
		v.Log(fmt.Sprintf("Getting value_from %s", assertion.ValueFrom.Url))
		parsedURL, err := url.Parse(assertion.ValueFrom.Url)
		if err != nil {
			panic(err)
		}
		if parsedURL.Scheme != "s3" && parsedURL.Scheme != "S3" {
			panic(fmt.Sprintf("Unsupported protocol for value_from: %s", parsedURL.Scheme))
		}
		content, err := v.GetValueFromS3(parsedURL.Host, parsedURL.Path)
		if err != nil {
			return "Error" // FIXME
		}
		v.Log(content)
		return content
	}
	return assertion.Value
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
