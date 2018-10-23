package assertion

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

// The non error cases are covered in match_test

func TestContainsWithNonJSONType(t *testing.T) {
	var complexNumber complex128
	_, err := contains(complexNumber, "foo", "1")
	if err == nil {
		t.Errorf("Expecting contains to return an error for non JSON encodable data")
	}
}

func TestDoesNotContainWithNonJSONType(t *testing.T) {
	var complexNumber complex128
	_, err := doesNotContain(complexNumber, "foo", "1")
	if err == nil {
		t.Errorf("Expecting doesNotContain to return an error for non JSON encodable data")
	}
}

func TestContainsWithString(t *testing.T) {
	s := "s3:Get*"
	match, err := contains(s, "Action", "*")
	assert.Nil(t, err, "Expecting no error from contains")
	assert.True(t, match.Match, "Expecting match for string")
}

func TestContainsWithSliceOfStrings(t *testing.T) {
	s := []string{"s3:Get*"}
	match, err := contains(s, "Action", "*")
	assert.Nil(t, err, "Expecting no error from contains")
	assert.True(t, match.Match, "Expecting match for string")
}
