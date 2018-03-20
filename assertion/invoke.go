package assertion

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type InvokeViolation struct {
	Message string
}

type InvokeResponse struct {
	Violations []InvokeViolation
}

func invoke(rule Rule, resource Resource, log LoggingFunction) (string, []Violation) {
	status := "OK"
	violations := make([]Violation, 0)
	payload := resource.Properties
	if rule.Invoke.Payload != "" {
		p, err := SearchData(rule.Invoke.Payload, resource.Properties)
		if err != nil {
			panic(err)
		}
		payload = p
	}
	payloadJSON, err := JSONStringify(payload)
	log(fmt.Sprintf("Invoke %s on %s\n", rule.Invoke.Url, payloadJSON))
	httpResponse, err := http.Get(rule.Invoke.Url)
	if err != nil {
		return rule.Severity, violations // TODO set violation to HTTP call failed
	}
	defer httpResponse.Body.Close()
	body, err := ioutil.ReadAll(httpResponse.Body)
	if err != nil {
		return rule.Severity, violations // TODO set violation to read body failed
	}
	log(string(body))
	var invokeResponse InvokeResponse
	err = json.Unmarshal(body, &invokeResponse)
	if err != nil {
		return rule.Severity, violations // TODO cannot parse response
	}
	for _, violation := range invokeResponse.Violations {
		status = rule.Severity
		violations = append(violations, Violation{
			RuleId:       rule.Id,
			Status:       status,
			ResourceId:   resource.Id,
			ResourceType: resource.Type,
			Filename:     resource.Filename,
			Message:      violation.Message,
		})
	}
	return status, violations
}
