package assertion

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

// InvokeViolation has message describing a single validation error
type InvokeViolation struct {
	Message string
}

// InvokeResponse contains a collection of validation errors
type InvokeResponse struct {
	Violations []InvokeViolation
}

// StandardExternalRuleInvoker implements an external HTTP or HTTPS call
type StandardExternalRuleInvoker struct {
}

func makeViolation(rule Rule, resource Resource, message string) Violation {
	return Violation{
		RuleID:           rule.ID,
		Status:           rule.Severity,
		ResourceID:       resource.ID,
		ResourceType:     resource.Type,
		Category:         resource.Category,
		Filename:         resource.Filename,
		RuleMessage:      rule.Message,
		AssertionMessage: message,
		CreatedAt:        currentTime(),
	}
}

func makeViolations(rule Rule, resource Resource, message string) []Violation {
	v := makeViolation(rule, resource, message)
	return []Violation{v}
}

// Invoke an external API to validate a Resource
func (e StandardExternalRuleInvoker) Invoke(rule Rule, resource Resource) (string, []Violation, error) {
	status := "OK"
	violations := make([]Violation, 0)
	var payload interface{}
	payload = resource
	if rule.Invoke.Payload != "" {
		p, err := SearchData(rule.Invoke.Payload, resource.Properties)
		if err != nil {
			return status, violations, err
		}
		payload = p
	}
	payloadJSON, err := JSONStringify(payload)
	if err != nil {
		violations := makeViolations(rule, resource, fmt.Sprintf("Unable to create JSON payload: %s", err.Error()))
		return rule.Severity, violations, err
	}
	Debugf("Invoke %s on %s\n", rule.Invoke.URL, payloadJSON)
	httpResponse, err := http.Post(rule.Invoke.URL, "application/json", bytes.NewBuffer([]byte(payloadJSON)))
	if err != nil {
		violations := makeViolations(rule, resource, fmt.Sprintf("Invoke failed: %s", err.Error()))
		return rule.Severity, violations, err
	}
	if httpResponse.StatusCode != 200 {
		violations := makeViolations(rule, resource, fmt.Sprintf("Invoke failed, StatusCode: %d", httpResponse.StatusCode))
		return rule.Severity, violations, nil
	}
	defer httpResponse.Body.Close()
	body, err := ioutil.ReadAll(httpResponse.Body)
	if err != nil {
		violations := makeViolations(rule, resource, "Invoke response cannot be read")
		return rule.Severity, violations, nil
	}
	Debugf("Invoke body: %s\n", string(body))
	var invokeResponse InvokeResponse
	err = json.Unmarshal(body, &invokeResponse)
	if err != nil {
		violations := makeViolations(rule, resource, "Invoke response cannot be parsed")
		return rule.Severity, violations, nil
	}
	for _, violation := range invokeResponse.Violations {
		status = rule.Severity
		v := makeViolation(rule, resource, violation.Message)
		violations = append(violations, v)
	}
	return status, violations, nil
}
