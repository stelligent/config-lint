package assertion

import (
	"gopkg.in/h2non/gock.v1"
	"testing"
)

func TestInvokeOK(t *testing.T) {
	defer gock.Off()

	response := InvokeResponse{}
	gock.New("http://config-lint.org").
		Post("/lint").
		Reply(200).
		JSON(response)

	i := StandardExternalRuleInvoker{}
	rule := Rule{
		Invoke: InvokeRuleAPI{
			URL: "http://config-lint.org/lint",
		},
	}
	resource := Resource{}
	status, violations, err := i.Invoke(rule, resource)
	if status != "OK" {
		t.Errorf("Expecting Invoke to return 'OK': %s\n", status)
	}
	if len(violations) != 0 {
		t.Errorf("Expecting Invoke to return no violations: %v\n", violations)
	}
	if err != nil {
		t.Errorf("Expecting Invoke to not return an error: %v\n", err.Error())
	}
}

func TestInvokeWithViolations(t *testing.T) {
	defer gock.Off()

	response := InvokeResponse{
		Violations: []InvokeViolation{
			InvokeViolation{Message: "Something is not right"},
		},
	}
	gock.New("http://config-lint.org").
		Post("/lint").
		Reply(200).
		JSON(response)

	i := StandardExternalRuleInvoker{}
	rule := Rule{
		Severity: "FAILURE",
		Invoke: InvokeRuleAPI{
			URL: "http://config-lint.org/lint",
		},
	}
	resource := Resource{}
	status, violations, err := i.Invoke(rule, resource)
	if status != "FAILURE" {
		t.Errorf("Expecting Invoke to return 'OK': %s\n", status)
	}
	if len(violations) != 1 {
		t.Errorf("Expecting Invoke to return 1 violations: %v\n", violations)
	}
	if err != nil {
		t.Errorf("Expecting Invoke to not return an error: %v\n", err.Error())
	}
}
