package assertion

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCommandLineVariable(t *testing.T) {
	s := StandardValueSource{
		Variables: map[string]string{"foo": "bar"},
	}
	e := Expression{
		ValueFrom: ValueFrom{Variable: "foo"},
	}
	v, err := s.GetValue(e)
	if err != nil {
		t.Errorf("Expected GetValue to return without error: %v\n", err.Error())
	}
	if v != "bar" {
		t.Errorf("Expected GetValue to find variable 'foo' with value 'bar', not '%s'\n", v)
	}
}

func TestValueFromHttp(t *testing.T) {
	cidrBlock := "0.0.0.0/0"
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, cidrBlock)
	}))
	defer ts.Close()
	s := StandardValueSource{}
	e := Expression{
		ValueFrom: ValueFrom{URL: ts.URL},
	}
	v, err := s.GetValue(e)
	assert.Nil(t, err, "Expecting GetValue to not return an error")
	assert.Equal(t, cidrBlock, v, "Expecting CIDR value to be returned")
}
