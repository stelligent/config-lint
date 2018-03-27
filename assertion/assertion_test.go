package assertion

import (
	"encoding/json"
	"testing"
)

func testLogging(s string) {
}

func failTestIfError(err error, message string, t *testing.T) {
	if err != nil {
		t.Error(message + ":" + err.Error())
	}
}

type AssertionTestCase struct {
	Rule           Rule
	Resource       Resource
	ExpectedStatus string
}

func TestCheckAssertion(t *testing.T) {

	testCases := map[string]AssertionTestCase{
		"testEq": {
			Rule{
				ID:       "test1",
				Message:  "test rule",
				Severity: "failure",
				Resource: "aws_instance",
				Assertions: []Assertion{
					Assertion{
						Type:  "value",
						Key:   "instance_type",
						Op:    "eq",
						Value: "t2.micro",
					},
				},
			},
			Resource{
				ID:         "a_test_resource",
				Type:       "aws_instance",
				Properties: map[string]interface{}{"instance_type": "t2.micro"},
				Filename:   "test.tf",
			},
			"OK",
		},
		"testOr": {
			Rule{
				ID:       "TEST1",
				Message:  "Test Rule",
				Severity: "FAILURE",
				Resource: "aws_instance",
				Assertions: []Assertion{
					Assertion{
						Or: []Assertion{
							Assertion{
								Type:  "value",
								Key:   "instance_type",
								Op:    "eq",
								Value: "t2.micro",
							},
							Assertion{
								Type:  "value",
								Key:   "instance_type",
								Op:    "eq",
								Value: "m4.large",
							},
						},
					},
				},
			},
			Resource{
				ID:         "a_test_resource",
				Type:       "aws_instance",
				Properties: map[string]interface{}{"instance_type": "t2.micro"},
				Filename:   "test.tf",
			},
			"OK",
		},
		"testOrFails": {
			Rule{
				ID:       "TEST1",
				Message:  "Test Rule",
				Severity: "FAILURE",
				Resource: "aws_instance",
				Assertions: []Assertion{
					Assertion{
						Or: []Assertion{
							Assertion{
								Type:  "value",
								Key:   "instance_type",
								Op:    "eq",
								Value: "t2.nano",
							},
							Assertion{
								Type:  "value",
								Key:   "instance_type",
								Op:    "eq",
								Value: "m4.large",
							},
						},
					},
				},
			},
			Resource{
				ID:         "a_test_resource",
				Type:       "aws_instance",
				Properties: map[string]interface{}{"instance_type": "t2.micro"},
				Filename:   "test.tf",
			},
			"FAILURE",
		},
		"testAnd": {
			Rule{
				ID:       "TEST1",
				Message:  "Test Rule",
				Severity: "FAILURE",
				Resource: "aws_instance",
				Assertions: []Assertion{
					Assertion{
						And: []Assertion{
							Assertion{
								Type:  "value",
								Key:   "instance_type",
								Op:    "eq",
								Value: "t2.micro",
							},
							Assertion{
								Type:  "value",
								Key:   "ami",
								Op:    "eq",
								Value: "ami-f2d3638a",
							},
						},
					},
				},
			},
			Resource{
				ID:   "a_test_resource",
				Type: "aws_instance",
				Properties: map[string]interface{}{
					"instance_type": "t2.micro",
					"ami":           "ami-f2d3638a",
				},
				Filename: "test.tf",
			},
			"OK",
		},
		"testAndFails": {
			Rule{
				ID:       "TEST1",
				Message:  "Test Rule",
				Severity: "FAILURE",
				Resource: "aws_instance",
				Assertions: []Assertion{
					Assertion{
						And: []Assertion{
							Assertion{
								Type:  "value",
								Key:   "instance_type",
								Op:    "eq",
								Value: "t2.micro",
							},
							Assertion{
								Type:  "value",
								Key:   "ami",
								Op:    "eq",
								Value: "ami-f2d3638a",
							},
						},
					},
				},
			},
			Resource{
				ID:   "a_test_resource",
				Type: "aws_instance",
				Properties: map[string]interface{}{
					"instance_type": "m3.medium",
					"ami":           "ami-f2d3638a",
				},
				Filename: "test.tf",
			},
			"FAILURE",
		},
		"testNot": {
			Rule{
				ID:       "TEST1",
				Message:  "Test Rule",
				Severity: "FAILURE",
				Resource: "aws_instance",
				Assertions: []Assertion{
					Assertion{
						Not: []Assertion{
							Assertion{
								Type:  "value",
								Key:   "instance_type",
								Op:    "eq",
								Value: "t2.micro",
							},
						},
					},
				},
			},
			Resource{
				ID:   "a_test_resource",
				Type: "aws_instance",
				Properties: map[string]interface{}{
					"instance_type": "c4.large",
				},
				Filename: "test.tf",
			},
			"OK",
		},
		"testNotFails": {
			Rule{
				ID:       "TEST1",
				Message:  "Test Rule",
				Severity: "FAILURE",
				Resource: "aws_instance",
				Assertions: []Assertion{
					Assertion{
						Not: []Assertion{
							Assertion{
								Type:  "value",
								Key:   "instance_type",
								Op:    "eq",
								Value: "t2.micro",
							},
						},
					},
				},
			},
			Resource{
				ID:   "a_test_resource",
				Type: "aws_instance",
				Properties: map[string]interface{}{
					"instance_type": "t2.micro",
				},
				Filename: "test.tf",
			},
			"FAILURE",
		},
		"testNestedNot": {
			Rule{
				ID:       "TEST1",
				Message:  "Test Rule",
				Severity: "FAILURE",
				Resource: "aws_instance",
				Assertions: []Assertion{
					Assertion{
						Not: []Assertion{
							Assertion{
								Or: []Assertion{
									Assertion{
										Type:  "value",
										Key:   "instance_type",
										Op:    "eq",
										Value: "t2.micro",
									},
									Assertion{
										Type:  "value",
										Key:   "instance_type",
										Op:    "eq",
										Value: "m3.medium",
									},
								},
							},
						},
					},
				},
			},
			Resource{
				ID:   "a_test_resource",
				Type: "aws_instance",
				Properties: map[string]interface{}{
					"instance_type": "m3.medium",
				},
				Filename: "test.tf",
			},
			"FAILURE",
		},
	}

	for k, tc := range testCases {
		status, err := CheckAssertion(tc.Rule, tc.Rule.Assertions[0], tc.Resource, testLogging)
		failTestIfError(err, "TestSimple", t)
		if status != tc.ExpectedStatus {
			t.Error("%s Failed Expected '%s' to be '%s'", k, status, tc.ExpectedStatus)
		}
	}
}

func TestNestedBooleans(t *testing.T) {
	rule := Rule{
		ID:       "TEST1",
		Message:  "Do not allow access to port 22 from 0.0.0.0/0",
		Severity: "NOT_COMPLIANT",
		Resource: "aws_instance",
		Assertions: []Assertion{
			Assertion{
				Not: []Assertion{
					Assertion{
						And: []Assertion{
							Assertion{
								Type:  "value",
								Key:   "ipPermissions[].fromPort[]",
								Op:    "contains",
								Value: "22",
							},
							Assertion{
								Type:  "value",
								Key:   "ipPermissions[].ipRanges[]",
								Op:    "contains",
								Value: "0.0.0.0/0",
							},
						},
					},
				},
			},
		},
	}
	resource := Resource{
		ID:         "a_test_resource",
		Type:       "aws_instance",
		Properties: map[string]interface{}{},
		Filename:   "test.tf",
	}
	rulesJSON := `{
            "description": "2017-12-03T03:14:29.856Z",
            "groupName": "test-8246",
            "ipPermissions": [
                {
                    "fromPort": 22,
                    "ipProtocol": "tcp",
                    "toPort": 22,
                    "ipv4Ranges": [
                        {
                            "cidrIp": "0.0.0.0/0"
                        }
                    ],
                    "ipRanges": [
                        "0.0.0.0/0"
                    ]
                }
            ]
        }`
	err := json.Unmarshal([]byte(rulesJSON), &resource.Properties)
	if err != nil {
		t.Error("Error parsing resource JSON")
	}
	status, err := CheckAssertion(rule, rule.Assertions[0], resource, testLogging)
	failTestIfError(err, "TestNestedBoolean", t)
	if status != "NOT_COMPLIANT" {
		t.Error("Expecting nested boolean to return NOT_COMPLIANT")
	}
}

func TestExceptions(t *testing.T) {
	rule := Rule{
		ID:     "EXCEPT",
		Except: []string{"200", "300"},
	}
	resources := []Resource{
		Resource{ID: "100"},
		Resource{ID: "200"},
		Resource{ID: "300"},
		Resource{ID: "400"},
	}
	filteredResources := FilterResourceExceptions(rule, resources)
	if len(filteredResources) != 2 {
		t.Error("Expecting exceptions to be removed from resource list")
	}
}

func TestNoExceptions(t *testing.T) {
	rule := Rule{
		ID:     "EXCEPT",
		Except: []string{},
	}
	resources := []Resource{
		Resource{ID: "100"},
		Resource{ID: "200"},
		Resource{ID: "300"},
		Resource{ID: "400"},
	}
	filteredResources := FilterResourceExceptions(rule, resources)
	if len(filteredResources) != 4 {
		t.Error("Expecting no exceptions to return all resources")
	}
}
