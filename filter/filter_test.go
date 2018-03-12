package filter

import (
	"encoding/json"
	"testing"
)

func testLogging(s string) {
}

func TestSimple(t *testing.T) {
	rule := Rule{
		Id:       "TEST1",
		Message:  "Test Rule",
		Severity: "FAILURE",
		Resource: "aws_instance",
		Filters: []Filter{
			Filter{
				Type:  "value",
				Key:   "instance_type",
				Op:    "eq",
				Value: "t2.micro",
			},
		},
	}
	resource := Resource{
		Id:         "a_test_resource",
		Type:       "aws_instance",
		Properties: map[string]interface{}{"instance_type": "t2.micro"},
		Filename:   "test.tf",
	}
	status := ApplyFilter(rule, rule.Filters[0], resource, testLogging)
	if status != "OK" {
		t.Error("Expecting simple rule to match")
	}
}

func TestOrToMatch(t *testing.T) {
	rule := Rule{
		Id:       "TEST1",
		Message:  "Test Rule",
		Severity: "FAILURE",
		Resource: "aws_instance",
		Filters: []Filter{
			Filter{
				Or: []Filter{
					Filter{
						Type:  "value",
						Key:   "instance_type",
						Op:    "eq",
						Value: "t2.micro",
					},
					Filter{
						Type:  "value",
						Key:   "instance_type",
						Op:    "eq",
						Value: "m4.large",
					},
				},
			},
		},
	}
	resource := Resource{
		Id:         "a_test_resource",
		Type:       "aws_instance",
		Properties: map[string]interface{}{"instance_type": "t2.micro"},
		Filename:   "test.tf",
	}
	status := ApplyFilter(rule, rule.Filters[0], resource, testLogging)
	if status != "OK" {
		t.Error("Expecting or to return OK")
	}
}

func TestOrToNotMatch(t *testing.T) {
	rule := Rule{
		Id:       "TEST1",
		Message:  "Test Rule",
		Severity: "FAILURE",
		Resource: "aws_instance",
		Filters: []Filter{
			Filter{
				Or: []Filter{
					Filter{
						Type:  "value",
						Key:   "instance_type",
						Op:    "eq",
						Value: "t2.micro",
					},
					Filter{
						Type:  "value",
						Key:   "instance_type",
						Op:    "eq",
						Value: "m4.large",
					},
				},
			},
		},
	}
	resource := Resource{
		Id:         "a_test_resource",
		Type:       "aws_instance",
		Properties: map[string]interface{}{"instance_type": "m3.medium"},
		Filename:   "test.tf",
	}
	status := ApplyFilter(rule, rule.Filters[0], resource, testLogging)
	if status != "FAILURE" {
		t.Error("Expecting or to return FAILURE")
	}
}

func TestAndToMatch(t *testing.T) {
	rule := Rule{
		Id:       "TEST1",
		Message:  "Test Rule",
		Severity: "FAILURE",
		Resource: "aws_instance",
		Filters: []Filter{
			Filter{
				And: []Filter{
					Filter{
						Type:  "value",
						Key:   "instance_type",
						Op:    "eq",
						Value: "t2.micro",
					},
					Filter{
						Type:  "value",
						Key:   "ami",
						Op:    "eq",
						Value: "ami-f2d3638a",
					},
				},
			},
		},
	}
	resource := Resource{
		Id:   "a_test_resource",
		Type: "aws_instance",
		Properties: map[string]interface{}{
			"instance_type": "t2.micro",
			"ami":           "ami-f2d3638a",
		},
		Filename: "test.tf",
	}
	status := ApplyFilter(rule, rule.Filters[0], resource, testLogging)
	if status != "OK" {
		t.Error("Expecting and to return OK")
	}
}

func TestAndToNotMatch(t *testing.T) {
	rule := Rule{
		Id:       "TEST1",
		Message:  "Test Rule",
		Severity: "FAILURE",
		Resource: "aws_instance",
		Filters: []Filter{
			Filter{
				And: []Filter{
					Filter{
						Type:  "value",
						Key:   "instance_type",
						Op:    "eq",
						Value: "t2.micro",
					},
					Filter{
						Type:  "value",
						Key:   "ami",
						Op:    "eq",
						Value: "ami-f2d3638a",
					},
				},
			},
		},
	}
	resource := Resource{
		Id:   "a_test_resource",
		Type: "aws_instance",
		Properties: map[string]interface{}{
			"instance_type": "m3.medium",
			"ami":           "ami-f2d3638a",
		},
		Filename: "test.tf",
	}
	status := ApplyFilter(rule, rule.Filters[0], resource, testLogging)
	if status != "FAILURE" {
		t.Error("Expecting and to return FAILURE")
	}
}

func TestNotToMatch(t *testing.T) {
	rule := Rule{
		Id:       "TEST1",
		Message:  "Test Rule",
		Severity: "FAILURE",
		Resource: "aws_instance",
		Filters: []Filter{
			Filter{
				Not: []Filter{
					Filter{
						Type:  "value",
						Key:   "instance_type",
						Op:    "eq",
						Value: "t2.micro",
					},
				},
			},
		},
	}
	resource := Resource{
		Id:   "a_test_resource",
		Type: "aws_instance",
		Properties: map[string]interface{}{
			"instance_type": "c4.large",
		},
		Filename: "test.tf",
	}
	status := ApplyFilter(rule, rule.Filters[0], resource, testLogging)
	if status != "OK" {
		t.Error("Expecting no to return OK")
	}
}

func TestNotToNotMatch(t *testing.T) {
	rule := Rule{
		Id:       "TEST1",
		Message:  "Test Rule",
		Severity: "FAILURE",
		Resource: "aws_instance",
		Filters: []Filter{
			Filter{
				Not: []Filter{
					Filter{
						Type:  "value",
						Key:   "instance_type",
						Op:    "eq",
						Value: "t2.micro",
					},
				},
			},
		},
	}
	resource := Resource{
		Id:   "a_test_resource",
		Type: "aws_instance",
		Properties: map[string]interface{}{
			"instance_type": "t2.micro",
		},
		Filename: "test.tf",
	}
	status := ApplyFilter(rule, rule.Filters[0], resource, testLogging)
	if status != "FAILURE" {
		t.Error("Expecting no to return FAILURE")
	}
}

func TestNestedNot(t *testing.T) {
	rule := Rule{
		Id:       "TEST1",
		Message:  "Test Rule",
		Severity: "FAILURE",
		Resource: "aws_instance",
		Filters: []Filter{
			Filter{
				Not: []Filter{
					Filter{
						Or: []Filter{
							Filter{
								Type:  "value",
								Key:   "instance_type",
								Op:    "eq",
								Value: "t2.micro",
							},
							Filter{
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
	}
	resource := Resource{
		Id:   "a_test_resource",
		Type: "aws_instance",
		Properties: map[string]interface{}{
			"instance_type": "m3.medium",
		},
		Filename: "test.tf",
	}
	status := ApplyFilter(rule, rule.Filters[0], resource, testLogging)
	if status != "FAILURE" {
		t.Error("Expecting nested boolean to return FAILURE")
	}
}

func TestNestedBooleans(t *testing.T) {
	rule := Rule{
		Id:       "TEST1",
		Message:  "Do not allow access to port 22 from 0.0.0.0/0",
		Severity: "NOT_COMPLIANT",
		Resource: "aws_instance",
		Filters: []Filter{
			Filter{
				Not: []Filter{
					Filter{
						And: []Filter{
							Filter{
								Type:  "value",
								Key:   "ipPermissions[].fromPort[]",
								Op:    "contains",
								Value: "22",
							},
							Filter{
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
		Id:         "a_test_resource",
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
	status := ApplyFilter(rule, rule.Filters[0], resource, testLogging)
	if status != "NOT_COMPLIANT" {
		t.Error("Expecting nested boolean to return NOT_COMPLIANT")
	}
}

func TestExceptions(t *testing.T) {
	rule := Rule{
		Id:     "EXCEPT",
		Except: []string{"200", "300"},
	}
	resources := []Resource{
		Resource{Id: "100"},
		Resource{Id: "200"},
		Resource{Id: "300"},
		Resource{Id: "400"},
	}
	filteredResources := FilterResourceExceptions(rule, resources)
	if len(filteredResources) != 2 {
		t.Error("Expecting exceptions to be removed from resource list")
	}
}

func TestNoExceptions(t *testing.T) {
	rule := Rule{
		Id:     "EXCEPT",
		Except: []string{},
	}
	resources := []Resource{
		Resource{Id: "100"},
		Resource{Id: "200"},
		Resource{Id: "300"},
		Resource{Id: "400"},
	}
	filteredResources := FilterResourceExceptions(rule, resources)
	if len(filteredResources) != 4 {
		t.Error("Expecting no exceptions to return all resources")
	}
}
