package assertion

import (
	"encoding/json"
	"testing"
)

type ExpressionTestCase struct {
	Rule           Rule
	Resource       Resource
	ExpectedStatus string
}

func TestCheckExpression(t *testing.T) {

	simpleTestResource := Resource{
		ID:   "a_test_resource",
		Type: "aws_instance",
		Properties: map[string]interface{}{
			"instance_type": "t2.micro",
			"ami":           "ami-f2d3638a",
		},
		Filename: "test.tf",
	}
	resourceWithTags := Resource{
		ID:   "another_test_resource",
		Type: "aws_instance",
		Properties: map[string]interface{}{
			"instance_type": "t2.micro",
			"ami":           "ami-f2d3638a",
			"tags": map[string]interface{}{
				"Environment": "Development",
				"Project":     "Web",
			},
		},
		Filename: "test.tf",
	}
	resourceWithRootVolume := Resource{
		ID:   "another_test_resource",
		Type: "aws_instance",
		Properties: map[string]interface{}{
			"instance_type": "t2.micro",
			"ami":           "ami-f2d3638a",
			"root_block_device": map[string]interface{}{
				"volume_size": "1000",
			},
		},
		Filename: "test.tf",
	}

	testCases := map[string]ExpressionTestCase{
		"testEq": {
			Rule{
				ID:       "test1",
				Message:  "test rule",
				Severity: "FAILURE",
				Resource: "aws_instance",
				Assertions: []Expression{
					Expression{
						Key:   "instance_type",
						Op:    "eq",
						Value: "t2.micro",
					},
				},
			},
			simpleTestResource,
			"OK",
		},
		"testOr": {
			Rule{
				ID:       "TEST1",
				Message:  "Test Rule",
				Severity: "FAILURE",
				Resource: "aws_instance",
				Assertions: []Expression{
					Expression{
						Or: []Expression{
							Expression{
								Key:   "instance_type",
								Op:    "eq",
								Value: "t2.micro",
							},
							Expression{
								Key:   "instance_type",
								Op:    "eq",
								Value: "m4.large",
							},
						},
					},
				},
			},
			simpleTestResource,
			"OK",
		},
		"testOrFails": {
			Rule{
				ID:       "TEST1",
				Message:  "Test Rule",
				Severity: "FAILURE",
				Resource: "aws_instance",
				Assertions: []Expression{
					Expression{
						Or: []Expression{
							Expression{
								Key:   "instance_type",
								Op:    "eq",
								Value: "t2.nano",
							},
							Expression{
								Key:   "instance_type",
								Op:    "eq",
								Value: "m4.large",
							},
						},
					},
				},
			},
			simpleTestResource,
			"FAILURE",
		},
		"testXor": {
			Rule{
				ID:       "TEST1",
				Message:  "Test Rule",
				Severity: "FAILURE",
				Resource: "aws_instance",
				Assertions: []Expression{
					Expression{
						Xor: []Expression{
							Expression{
								Key:   "instance_type",
								Op:    "eq",
								Value: "t2.micro",
							},
							Expression{
								Key:   "instance_type",
								Op:    "eq",
								Value: "m4.large",
							},
						},
					},
				},
			},
			simpleTestResource,
			"OK",
		},
		"testXorFails": {
			Rule{
				ID:       "TEST1",
				Message:  "Test Rule",
				Severity: "FAILURE",
				Resource: "aws_instance",
				Assertions: []Expression{
					Expression{
						Xor: []Expression{
							Expression{
								Key:   "instance_type",
								Op:    "eq",
								Value: "t2.micro",
							},
							Expression{
								Key:   "instance_type",
								Op:    "eq",
								Value: "t2.micro",
							},
						},
					},
				},
			},
			simpleTestResource,
			"FAILURE",
		},
		"testXorFailsAgain": {
			Rule{
				ID:       "TEST1",
				Message:  "Test Rule",
				Severity: "FAILURE",
				Resource: "aws_instance",
				Assertions: []Expression{
					Expression{
						Xor: []Expression{
							Expression{
								Key:   "instance_type",
								Op:    "eq",
								Value: "m3.large",
							},
							Expression{
								Key:   "instance_type",
								Op:    "eq",
								Value: "c4.large",
							},
						},
					},
				},
			},
			simpleTestResource,
			"FAILURE",
		},
		"testAnd": {
			Rule{
				ID:       "TEST1",
				Message:  "Test Rule",
				Severity: "FAILURE",
				Resource: "aws_instance",
				Assertions: []Expression{
					Expression{
						And: []Expression{
							Expression{
								Key:   "instance_type",
								Op:    "eq",
								Value: "t2.micro",
							},
							Expression{
								Key:   "ami",
								Op:    "eq",
								Value: "ami-f2d3638a",
							},
						},
					},
				},
			},
			simpleTestResource,
			"OK",
		},
		"testAndFails": {
			Rule{
				ID:       "TEST1",
				Message:  "Test Rule",
				Severity: "FAILURE",
				Resource: "aws_instance",
				Assertions: []Expression{
					Expression{
						And: []Expression{
							Expression{
								Key:   "instance_type",
								Op:    "eq",
								Value: "m3.medium",
							},
							Expression{
								Key:   "ami",
								Op:    "eq",
								Value: "ami-f2d3638a",
							},
						},
					},
				},
			},
			simpleTestResource,
			"FAILURE",
		},
		"testNot": {
			Rule{
				ID:       "TEST1",
				Message:  "Test Rule",
				Severity: "FAILURE",
				Resource: "aws_instance",
				Assertions: []Expression{
					Expression{
						Not: []Expression{
							Expression{
								Key:   "instance_type",
								Op:    "eq",
								Value: "c4.large",
							},
						},
					},
				},
			},
			simpleTestResource,
			"OK",
		},
		"testNotFails": {
			Rule{
				ID:       "TEST1",
				Message:  "Test Rule",
				Severity: "FAILURE",
				Resource: "aws_instance",
				Assertions: []Expression{
					Expression{
						Not: []Expression{
							Expression{
								Key:   "instance_type",
								Op:    "eq",
								Value: "t2.micro",
							},
						},
					},
				},
			},
			simpleTestResource,
			"FAILURE",
		},
		"testNestedNot": {
			Rule{
				ID:       "TEST1",
				Message:  "Test Rule",
				Severity: "FAILURE",
				Resource: "aws_instance",
				Assertions: []Expression{
					Expression{
						Not: []Expression{
							Expression{
								Or: []Expression{
									Expression{
										Key:   "instance_type",
										Op:    "eq",
										Value: "t2.micro",
									},
									Expression{
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
			simpleTestResource,
			"FAILURE",
		},
		"testSizeFails": {
			Rule{
				ID:       "TESTCOUNT",
				Message:  "Test Resource Count Fails",
				Severity: "FAILURE",
				Resource: "aws_instance",
				Assertions: []Expression{
					Expression{
						Key:       "tags",
						ValueType: "size",
						Op:        "eq",
						Value:     "3",
					},
				},
			},
			resourceWithTags,
			"FAILURE",
		},
		"testSizeOK": {
			Rule{
				ID:       "TESTCOUNT",
				Message:  "Test Resource Count OK",
				Severity: "FAILURE",
				Resource: "aws_instance",
				Assertions: []Expression{
					Expression{
						Key:       "tags",
						ValueType: "size",
						Op:        "eq",
						Value:     "2",
					},
				},
			},
			resourceWithTags,
			"OK",
		},
		"testIntegerFails": {
			Rule{
				ID:       "TESTCOUNT",
				Message:  "Test integer Fails",
				Severity: "FAILURE",
				Resource: "aws_instance",
				Assertions: []Expression{
					Expression{
						Key:       "root_block_device.volume_size",
						ValueType: "integer",
						Op:        "le",
						Value:     "500",
					},
				},
			},
			resourceWithRootVolume,
			"FAILURE",
		},
		"testIntegerOK": {
			Rule{
				ID:       "TESTCOUNT",
				Message:  "Test integer OK",
				Severity: "FAILURE",
				Resource: "aws_instance",
				Assertions: []Expression{
					Expression{
						Key:       "root_block_device.volume_size",
						ValueType: "integer",
						Op:        "le",
						Value:     "2000",
					},
				},
			},
			resourceWithRootVolume,
			"OK",
		},
	}

	for k, tc := range testCases {
		expressionResult, err := CheckExpression(tc.Rule, tc.Rule.Assertions[0], tc.Resource)
		FailTestIfError(err, "TestSimple", t)
		if expressionResult.Status != tc.ExpectedStatus {
			t.Errorf("%s Failed Expected '%s' to be '%s'", k, expressionResult.Status, tc.ExpectedStatus)
		}
	}
}

func TestNestedBooleans(t *testing.T) {
	rule := Rule{
		ID:       "TEST1",
		Message:  "Do not allow access to port 22 from 0.0.0.0/0",
		Severity: "NOT_COMPLIANT",
		Resource: "aws_instance",
		Assertions: []Expression{
			Expression{
				Not: []Expression{
					Expression{
						And: []Expression{
							Expression{
								Key:   "ipPermissions[].fromPort[]",
								Op:    "contains",
								Value: "22",
							},
							Expression{
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
	resourceJSON := `{
            "description": "2017-12-03T03:14:29.856Z",
            "groupName": "test-8246",
            "ipPermissions": [
                {
                    "fromPort": "22",
                    "ipProtocol": "tcp",
                    "toPort": "22",
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
	err := json.Unmarshal([]byte(resourceJSON), &resource.Properties)
	if err != nil {
		t.Error("Error parsing resource JSON")
	}
	expressionResult, err := CheckExpression(rule, rule.Assertions[0], resource)
	FailTestIfError(err, "TestNestedBoolean", t)
	if expressionResult.Status != "NOT_COMPLIANT" {
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

func TestUsingFixtures(t *testing.T) {
	fixtureFilenames := []string{
		"./testdata/collection-assertions.yaml",
		"./testdata/has-properties.yaml",
	}

	for _, filename := range fixtureFilenames {
		RunTestCasesFromFixture(filename, t)
	}
}
