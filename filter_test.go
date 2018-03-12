package main

import (
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
	resource := TerraformResource{
		Id:         "a_test_resource",
		Type:       "aws_instance",
		Properties: map[string]interface{}{"instance_type": "t2.micro"},
		Filename:   "test.tf",
	}
	status := applyFilter(rule, rule.Filters[0], resource, testLogging)
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
	resource := TerraformResource{
		Id:         "a_test_resource",
		Type:       "aws_instance",
		Properties: map[string]interface{}{"instance_type": "t2.micro"},
		Filename:   "test.tf",
	}
	status := applyFilter(rule, rule.Filters[0], resource, testLogging)
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
	resource := TerraformResource{
		Id:         "a_test_resource",
		Type:       "aws_instance",
		Properties: map[string]interface{}{"instance_type": "m3.medium"},
		Filename:   "test.tf",
	}
	status := applyFilter(rule, rule.Filters[0], resource, testLogging)
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
	resource := TerraformResource{
		Id:   "a_test_resource",
		Type: "aws_instance",
		Properties: map[string]interface{}{
			"instance_type": "t2.micro",
			"ami":           "ami-f2d3638a",
		},
		Filename: "test.tf",
	}
	status := applyFilter(rule, rule.Filters[0], resource, testLogging)
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
	resource := TerraformResource{
		Id:   "a_test_resource",
		Type: "aws_instance",
		Properties: map[string]interface{}{
			"instance_type": "m3.medium",
			"ami":           "ami-f2d3638a",
		},
		Filename: "test.tf",
	}
	status := applyFilter(rule, rule.Filters[0], resource, testLogging)
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
	resource := TerraformResource{
		Id:   "a_test_resource",
		Type: "aws_instance",
		Properties: map[string]interface{}{
			"instance_type": "c4.large",
		},
		Filename: "test.tf",
	}
	status := applyFilter(rule, rule.Filters[0], resource, testLogging)
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
	resource := TerraformResource{
		Id:   "a_test_resource",
		Type: "aws_instance",
		Properties: map[string]interface{}{
			"instance_type": "t2.micro",
		},
		Filename: "test.tf",
	}
	status := applyFilter(rule, rule.Filters[0], resource, testLogging)
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
	resource := TerraformResource{
		Id:   "a_test_resource",
		Type: "aws_instance",
		Properties: map[string]interface{}{
			"instance_type": "m3.medium",
		},
		Filename: "test.tf",
	}
	status := applyFilter(rule, rule.Filters[0], resource, testLogging)
	if status != "FAILURE" {
		t.Error("Expecting nested boolean to return FAILURE")
	}
}
