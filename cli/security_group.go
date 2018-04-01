package main

import (
	"encoding/json"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/stelligent/config-lint/assertion"
)

type (
	// SecurityGroupLoader calls the AWS SDK DescribeSecurityGroups
	SecurityGroupLoader struct{}
)

// Load gets security group information from AWS and generates Resources suitable for linting
func (sg SecurityGroupLoader) Load() ([]assertion.Resource, error) {
	resources := make([]assertion.Resource, 0)
	region := &aws.Config{Region: aws.String("us-east-1")}
	awsSession := session.New()
	ec2Client := ec2.New(awsSession, region)
	response, err := ec2Client.DescribeSecurityGroups(&ec2.DescribeSecurityGroupsInput{})
	if err != nil {
		return resources, err
	}
	for _, securityGroup := range response.SecurityGroups {

		// convert to JSON string
		jsonData, err := json.Marshal(securityGroup)
		if err != nil {
			return resources, err
		}

		// then convert to an interface{}
		// seem to need this for JMESPath to work properly
		var data interface{}
		err = json.Unmarshal(jsonData, &data)
		if err != nil {
			return resources, err
		}

		r := assertion.Resource{
			ID:         *securityGroup.GroupId,
			Type:       "AWS::EC2::SecurityGroup",
			Properties: data,
		}
		resources = append(resources, r)
	}
	return resources, nil
}
