package assertion

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestInvokeOK(t *testing.T) {

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "{}")
	}))
	defer ts.Close()

	i := StandardExternalRuleInvoker{}
	rule := Rule{
		Invoke: InvokeRuleAPI{
			URL: ts.URL,
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
	response := InvokeResponse{
		Violations: []InvokeViolation{
			InvokeViolation{Message: "Something is not right"},
		},
	}
	jsonData, err := json.Marshal(response)
	if err != nil {
		t.Errorf("Failed to marshal test response: %v\n", err.Error())
	}
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, string(jsonData))
	}))
	defer ts.Close()

	i := StandardExternalRuleInvoker{}
	rule := Rule{
		Severity: "FAILURE",
		Invoke: InvokeRuleAPI{
			URL: ts.URL,
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

func TestInvokeSendsMetadata(t *testing.T) {

	var invokedResource Resource
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := ioutil.ReadAll(r.Body)
		_ = json.Unmarshal(body, &invokedResource)
		fmt.Fprintln(w, "{}")
	}))
	defer ts.Close()

	i := StandardExternalRuleInvoker{}
	rule := Rule{
		Invoke: InvokeRuleAPI{
			URL: ts.URL,
		},
	}
	resource := Resource{
		Filename: "example.tf",
	}
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
	if invokedResource.Filename != resource.Filename {
		t.Errorf("Expecting filename metadata to be passed to external endpoint: %s != %s\n", resource.Filename, invokedResource.Filename)
	}
}
