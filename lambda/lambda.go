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

type ConfigurationItem struct {
	ResourceType                 string
	ResourceId                   string
	ConfigurationItemCaptureTime *time.Time
	Configuration                interface{}
}

type InvokingEvent struct {
	ConfigurationItem ConfigurationItem
}

func printValue(expression string, data interface{}) {
	fmt.Println("expression:", expression)
	value, err := assertion.SearchData(expression, data)
	if err != nil {
		fmt.Println("err:", err)
	}
	s, err := assertion.JSONStringify(value)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Printf("Value: %s\n", s)
	}
}

func log(s string) {
	fmt.Println(s)
}

type RuleParameters struct {
	Bucket string
	Key    string
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
		fmt.Println(err)
		return "Error parsing JSON for invokingEvent", nil
	}

	var ruleParameters RuleParameters
	err = json.Unmarshal([]byte(configEvent.RuleParameters), &ruleParameters)
	if err != nil {
		fmt.Println(err)
		return "Error parsing JSON for RuleParameters", nil
	}
	fmt.Printf("Rules object in bucket: %s key: %s\n", ruleParameters.Bucket, ruleParameters.Key)

	rulesString, err := loadRulesFromS3(s3Client, ruleParameters.Bucket, ruleParameters.Key)
	if err != nil {
		fmt.Println(err)
		return "Cannot GetObject", nil
	}
	fmt.Println("rulesString:", rulesString)

	fmt.Println("invokingEvent:", invokingEvent)
	configurationItem := invokingEvent.ConfigurationItem
	fmt.Println("configurationItem:", configurationItem)
	fmt.Println("configuration:", configurationItem.Configuration)

	complianceType := "NOT_APPLICABLE"
	ruleSet := assertion.MustParseRules(rulesString)
	valueSource := assertion.StandardValueSource{Log: log}
	resolvedRules := assertion.ResolveRules(ruleSet.Rules, valueSource, log)
	externalRules := assertion.StandardExternalRuleInvoker{Log: log}
	for _, rule := range resolvedRules {
		if rule.Resource == configurationItem.ResourceType {
			resource := assertion.Resource{
				Id:         configurationItem.ResourceId,
				Type:       configurationItem.ResourceType,
				Properties: configurationItem.Configuration,
			}
			_, violations := assertion.CheckRule(rule, resource, externalRules, log)
			if len(violations) > 0 {
				fmt.Println("Resource in NON_COMPLIANT")
				complianceType = "NON_COMPLIANT"
				for _, violation := range violations {
					fmt.Println(violation)
				}
			} else {
				fmt.Println("Resource in COMPLIANT")
				complianceType = "COMPLIANT"
			}
		} else {
			fmt.Println("Ignoring Resource")
		}
	}

	params := &configservice.PutEvaluationsInput{
		Evaluations: []*configservice.Evaluation{
			&configservice.Evaluation{
				ComplianceResourceType: aws.String(configurationItem.ResourceType),
				ComplianceResourceId:   aws.String(configurationItem.ResourceId),
				ComplianceType:         aws.String(complianceType),
				OrderingTimestamp:      aws.Time(time.Now()),
			},
		},
		ResultToken: aws.String(configEvent.ResultToken),
	}
	fmt.Println("params:", params)
	response, err := config.PutEvaluations(params)
	if err != nil {
		fmt.Println("err:", err)
	}
	fmt.Println("response:", response)
	return "Done", nil
}

func main() {
	lambda.Start(handler)
}
