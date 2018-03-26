package assertion

import (
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
	Log LoggingFunction
}

// Invoke an external API to validate a Resource
func (e StandardExternalRuleInvoker) Invoke(rule Rule, resource Resource) (string, []Violation, error) {
	status := "OK"
	violations := make([]Violation, 0)
	payload := resource.Properties
	if rule.Invoke.Payload != "" {
		p, err := SearchData(rule.Invoke.Payload, resource.Properties)
		if err != nil {
			return status, violations, err
		}
		payload = p
	}
	payloadJSON, err := JSONStringify(payload)
	e.Log(fmt.Sprintf("Invoke %s on %s\n", rule.Invoke.URL, payloadJSON))
	httpResponse, err := http.Get(rule.Invoke.URL)
	if err != nil {
		violations := []Violation{
			Violation{
				RuleID:       rule.ID,
				Status:       rule.Severity,
				ResourceID:   resource.ID,
				ResourceType: resource.Type,
				Filename:     resource.Filename,
				Message:      fmt.Sprintf("Invoke failed: %s", err.Error()),
			},
		}
		return rule.Severity, violations, err
	}
	if httpResponse.StatusCode != 200 {
		violations := []Violation{
			Violation{
				RuleID:       rule.ID,
				Status:       rule.Severity,
				ResourceID:   resource.ID,
				ResourceType: resource.Type,
				Filename:     resource.Filename,
				Message:      fmt.Sprintf("Invoke failed, StatusCode: %d", httpResponse.StatusCode),
			},
		}
		return rule.Severity, violations, nil
	}
	defer httpResponse.Body.Close()
	body, err := ioutil.ReadAll(httpResponse.Body)
	if err != nil {
		violations := []Violation{
			Violation{
				RuleID:       rule.ID,
				Status:       rule.Severity,
				ResourceID:   resource.ID,
				ResourceType: resource.Type,
				Filename:     resource.Filename,
				Message:      "Invoke response cannot be read",
			},
		}
		return rule.Severity, violations, nil
	}
	e.Log(string(body))
	var invokeResponse InvokeResponse
	err = json.Unmarshal(body, &invokeResponse)
	if err != nil {
		violations := []Violation{
			Violation{
				RuleID:       rule.ID,
				Status:       rule.Severity,
				ResourceID:   resource.ID,
				ResourceType: resource.Type,
				Filename:     resource.Filename,
				Message:      "Invoke response cannot be parsed",
			},
		}
		return rule.Severity, violations, nil
	}
	for _, violation := range invokeResponse.Violations {
		status = rule.Severity
		violations = append(violations, Violation{
			RuleID:       rule.ID,
			Status:       status,
			ResourceID:   resource.ID,
			ResourceType: resource.Type,
			Filename:     resource.Filename,
			Message:      violation.Message,
		})
	}
	return status, violations, nil
}
