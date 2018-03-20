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

type StandardExternalRuleInvoker struct {
	Log LoggingFunction
}

func (e StandardExternalRuleInvoker) Invoke(rule Rule, resource Resource) (string, []Violation) {
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
	e.Log(fmt.Sprintf("Invoke %s on %s\n", rule.Invoke.Url, payloadJSON))
	httpResponse, err := http.Get(rule.Invoke.Url)
	if err != nil {
		return rule.Severity, violations // TODO set violation to HTTP call failed
	}
	defer httpResponse.Body.Close()
	body, err := ioutil.ReadAll(httpResponse.Body)
	if err != nil {
		return rule.Severity, violations // TODO set violation to read body failed
	}
	e.Log(string(body))
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
