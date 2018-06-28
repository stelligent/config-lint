package assertion

import (
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
