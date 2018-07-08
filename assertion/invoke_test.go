package assertion

import (
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/assert"
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
	assert.Equal(t, "OK", status, "Expecting Invoke to return 'OK'")
	assert.Equal(t, 0, len(violations), "Expecting Invoke to return no violations")
	assert.Nil(t, err, "Expecting Invoke to not return an error")
}

func TestInvokeWithViolations(t *testing.T) {
	response := InvokeResponse{
		Violations: []InvokeViolation{
			InvokeViolation{Message: "Something is not right"},
		},
	}
	jsonData, err := json.Marshal(response)
	assert.Nil(t, err, "Failed to marshal test response")
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
	assert.Equal(t, "FAILURE", status, "Expecting Invoke to return 'FAILURE'")
	assert.Equal(t, 1, len(violations), "Expecting Invoke to return 1 violation")
	assert.Nil(t, err, "Expecting Invoke to not return an error")
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
	i.Invoke(rule, resource)
	assert.Equal(t, resource.Filename, invokedResource.Filename, "Expecting filename metadata in request body")
}
