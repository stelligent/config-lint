package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/configservice"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/stelligent/config-lint/assertion"
	"time"
)

type (

	// ConfigurationItem to hold info from Lambda event
	ConfigurationItem struct {
		ResourceType                 string
		ResourceID                   string
		ConfigurationItemCaptureTime *time.Time
		Configuration                interface{}
	}

	// InvokingEvent for Lambda event info
	InvokingEvent struct {
		ConfigurationItem ConfigurationItem
	}

	// RuleParameters bucket and key for loading RuleSet
	RuleParameters struct {
		Bucket string
		Key    string
	}
)

func main() {
	lambda.Start(handler)
}

func handler(configEvent events.ConfigEvent) (string, error) {
	fmt.Printf("AWS Config rule: %s\n", configEvent.ConfigRuleName)
	fmt.Printf("Invoking event JSON: %s\n", configEvent.InvokingEvent)
	fmt.Printf("Event version: %s\n", configEvent.Version)
	fmt.Printf("Rule parameters: %s\n", configEvent.RuleParameters)

	region := &aws.Config{Region: aws.String("us-east-1")}
	awsSession := session.New()
	config := configservice.New(awsSession, region)
	s3Client := s3.New(awsSession)

	var invokingEvent InvokingEvent
	err := json.Unmarshal([]byte(configEvent.InvokingEvent), &invokingEvent)
	if err != nil {
		fmt.Printf("Parse InvokingEvent failed: %b\n", err)
		return err.Error(), nil
	}

	var ruleParameters RuleParameters
	err = json.Unmarshal([]byte(configEvent.RuleParameters), &ruleParameters)
	if err != nil {
		fmt.Printf("Parse RuleParameters failed: %v\n", err)
		return err.Error(), nil
	}
	fmt.Printf("Rules object in bucket: %s key: %s\n", ruleParameters.Bucket, ruleParameters.Key)

	rulesString, err := loadRulesFromS3(s3Client, ruleParameters.Bucket, ruleParameters.Key)
	if err != nil {
		fmt.Printf("loadRulesFromS3 failed: %v\n", err)
		return err.Error(), nil
	}
	fmt.Println("rulesString:", rulesString)

	fmt.Println("invokingEvent:", invokingEvent)
	configurationItem := invokingEvent.ConfigurationItem
	fmt.Println("configurationItem:", configurationItem)
	fmt.Println("configuration:", configurationItem.Configuration)

	ruleSet, err := assertion.ParseRules(rulesString)
	if err != nil {
		fmt.Printf("Unable to parse rules: %v\n", err.Error())
		return "checkCompliance failed", err
	}
	valueSource := assertion.StandardValueSource{}
	resolvedRules := assertion.ResolveRules(ruleSet.Rules, valueSource)
	externalRules := assertion.StandardExternalRuleInvoker{}

	complianceType, err := checkCompliance(resolvedRules, configurationItem, externalRules)
	if err != nil {
		fmt.Printf("checkCompliance failed: %v\n", err)
		return err.Error(), nil
	}

	params := &configservice.PutEvaluationsInput{
		Evaluations: []*configservice.Evaluation{
			&configservice.Evaluation{
				ComplianceResourceType: aws.String(configurationItem.ResourceType),
				ComplianceResourceId:   aws.String(configurationItem.ResourceID),
				ComplianceType:         aws.String(complianceType),
				OrderingTimestamp:      aws.Time(time.Now()),
			},
		},
		ResultToken: aws.String(configEvent.ResultToken),
	}
	fmt.Println("params:", params)
	response, err := config.PutEvaluations(params)
	if err != nil {
		fmt.Printf("PutEvaluations failed: %v\n", err)
		return err.Error(), nil
	}
	fmt.Printf("PutEvaluations response: %v\n", response)
	return "Done", nil
}

func checkCompliance(rules []assertion.Rule, configurationItem ConfigurationItem, invoker assertion.ExternalRuleInvoker) (string, error) {
	complianceType := "NOT_APPLICABLE"
	for _, rule := range rules {
		if rule.Resource == configurationItem.ResourceType {
			resource := assertion.Resource{
				ID:         configurationItem.ResourceID,
				Type:       configurationItem.ResourceType,
				Properties: configurationItem.Configuration,
			}
			_, violations, err := assertion.CheckRule(rule, resource, invoker)
			if err != nil {
				return "NOT_APPLICABLE", err
			}
			if len(violations) > 0 {
				fmt.Println("Resource is NON_COMPLIANT")
				complianceType = "NON_COMPLIANT"
				for _, violation := range violations {
					fmt.Println(violation)
				}
			} else {
				fmt.Println("Resource is COMPLIANT")
				complianceType = "COMPLIANT"
			}
		} else {
			fmt.Println("Ignoring Resource")
		}
	}
	return complianceType, nil
}

func loadRulesFromS3(s3Client *s3.S3, bucket string, key string) (string, error) {
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
