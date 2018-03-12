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
	value, err := searchData(expression, data)
	if err != nil {
		fmt.Println("err:", err)
	}
	fmt.Println("value:", value)
}

func log(string) {
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
	fmt.Println("Bucket:", ruleParameters.Bucket)
	fmt.Println("Key:", ruleParameters.Key)

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

	printValue("@", configurationItem.Configuration)

	complianceType := "NOT_APPLICABLE"
	ruleSet := MustParseRules(rulesString)
	for _, rule := range ruleSet.Rules {
		if rule.Resource == configurationItem.ResourceType {
			complianceType = "COMPLIANT"
			for _, filter := range rule.Filters {
				resource := KubernetesResource{
					Id:         configurationItem.ResourceId,
					Type:       configurationItem.ResourceType,
					Properties: configurationItem.Configuration,
				}
				status := applyFilter(rule, filter, resource, log)
				fmt.Println(status, resource)
				if status != "OK" {
					complianceType = status
				}
			}
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

func handler2(configEvent events.ConfigEvent) (string, error) {
	fmt.Println(configEvent)
	return "Done", nil
}

func main() {
	lambda.Start(handler)
}
