package linter

import (
	"encoding/json"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/stelligent/config-lint/assertion"
	"net/url"
)

type (
	// IAMUserLoader calls the AWS SDK to get user information
	IAMUserLoader struct{}
	// IAMRoleLoader calls the AWS SDK to get user information
	IAMRoleLoader struct{}
	// IAMGroupLoader calls the AWS SDK to get user information
	IAMGroupLoader struct{}
)

// Load gets user information from AWS and generates Resources suitable for linting
func (u IAMUserLoader) Load() ([]assertion.Resource, error) {
	resources := make([]assertion.Resource, 0)
	region := &aws.Config{Region: aws.String("us-east-1")}
	awsSession := session.New()
	iamClient := iam.New(awsSession, region)

	response, err := iamClient.ListUsers(&iam.ListUsersInput{})
	if err != nil {
		return resources, err
	}
	for _, user := range response.Users {

		// convert to JSON string
		jsonData, err := json.Marshal(user)
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

		userPolicyResponse, err := iamClient.ListUserPolicies(&iam.ListUserPoliciesInput{
			UserName: aws.String(*user.UserName),
		})
		if err != nil {
			return resources, err
		}
		policies := make([]map[string]interface{}, 0)
		for _, policyName := range userPolicyResponse.PolicyNames {

			policyResponse, err := iamClient.GetUserPolicy(&iam.GetUserPolicyInput{
				UserName:   aws.String(*user.UserName),
				PolicyName: aws.String(*policyName),
			})
			if err != nil {
				return resources, err
			}
			decoded, err := url.QueryUnescape(*policyResponse.PolicyDocument)
			if err != nil {
				return resources, err
			}
			policies = append(policies, map[string]interface{}{
				"PolicyName":     *policyName,
				"PolicyDocument": decoded,
			})
		}
		m := data.(map[string]interface{})
		m["Policies"] = policies

		r := assertion.Resource{
			ID:         *user.UserName,
			Type:       "AWS::IAM::User",
			Properties: data,
		}
		resources = append(resources, r)

	}

	return resources, nil
}
