package main

import (
	"encoding/json"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/configservice"

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

func noLog(string) {
}

func handler(configEvent events.ConfigEvent) (string, error) {
	fmt.Printf("AWS Config rule: %s\n", configEvent.ConfigRuleName)
	fmt.Printf("Invoking event JSON: %s\n", configEvent.InvokingEvent)
	fmt.Printf("Event version: %s\n", configEvent.Version)

	region := &aws.Config{Region: aws.String("us-east-1")}
	config := configservice.New(session.New(), region)

	var invokingEvent InvokingEvent
	err := json.Unmarshal([]byte(configEvent.InvokingEvent), &invokingEvent)
	if err == nil {
		fmt.Println("invokingEvent:", invokingEvent)
		configurationItem := invokingEvent.ConfigurationItem
		fmt.Println("configurationItem:", configurationItem)
		fmt.Println("configuration:", configurationItem.Configuration)

		//printValue("@", configurationItem.Configuration)
		printValue("provisionedThroughput", configurationItem.Configuration)
		printValue("provisionedThroughput.writeCapacityUnits", configurationItem.Configuration)
		printValue("provisionedThroughput.readCapacityUnits", configurationItem.Configuration)
		printValue("tableName", configurationItem.Configuration)
		printValue("tableStatus", configurationItem.Configuration)

		ruleSet := MustParseRules(loadRules("./example-files/rules/aws-config.yml")) // FIXME remove Terraform from function name
		for _, rule := range ruleSet.Rules {
			for _, filter := range rule.Filters {
				resource := KubernetesResource{
					Id:         configurationItem.ResourceId,
					Type:       configurationItem.ResourceType,
					Properties: configurationItem.Configuration,
				}
				status := applyFilter(rule, filter, resource, noLog)
				fmt.Println(status, resource)
			}
		}

		params := &configservice.PutEvaluationsInput{
			Evaluations: []*configservice.Evaluation{
				&configservice.Evaluation{
					ComplianceResourceType: aws.String(configurationItem.ResourceType),
					ComplianceResourceId:   aws.String(configurationItem.ResourceId),
					ComplianceType:         aws.String("COMPLIANT"),
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
	}
	return "Done", nil
}

func handler2(configEvent events.ConfigEvent) (string, error) {
	fmt.Println(configEvent)
	return "Done", nil
}

func main() {
	lambda.Start(handler)
}
