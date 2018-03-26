package main

import (
	"encoding/json"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/stelligent/config-lint/assertion"
)

// SecurityGroupLinter implements a Linter for data returned by the DescribeSecurityGroups SDK call
type SecurityGroupLinter struct {
	BaseLinter
	Log assertion.LoggingFunction
}

func loadSecurityGroupResources(log assertion.LoggingFunction) ([]assertion.Resource, error) {
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

// Validate applies a Ruleset to all SecurityGroups
func (l SecurityGroupLinter) Validate(filenames []string, ruleSet assertion.RuleSet, tags []string, ruleIDs []string) ([]string, []assertion.Violation, error) {
	noFilenames := []string{}
	rules := assertion.FilterRulesByID(ruleSet.Rules, ruleIDs)
	resources, err := loadSecurityGroupResources(l.Log)
	if err != nil {
		return noFilenames, []assertion.Violation{}, err
	}
	violations, err := l.ValidateResources(resources, rules, tags, l.Log)
	return noFilenames, violations, err
}

// Search applies a JMESPath to all SecurityGroups
func (l SecurityGroupLinter) Search(filenames []string, ruleSet assertion.RuleSet, searchExpression string) {
	resources, _ := loadSecurityGroupResources(l.Log) // FIXME what about error?
	for _, resource := range resources {
		v, err := assertion.SearchData(searchExpression, resource.Properties)
		if err != nil {
			fmt.Println(err)
		} else {
			s, err := assertion.JSONStringify(v)
			if err != nil {
				fmt.Println(err)
			} else {
				fmt.Printf("%s: %s\n", resource.ID, s)
			}
		}
	}
}
