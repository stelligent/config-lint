package main

import (
	"bytes"
	"github.com/stelligent/config-lint/assertion"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestReportWriter(t *testing.T) {
	var b bytes.Buffer
	w := DefaultReportWriter{Writer: &b}
	r := assertion.ValidationReport{
		Violations: []assertion.Violation{
			assertion.Violation{
				RuleID: "RULE_1",
			},
		},
	}
	w.WriteReport(r, LinterOptions{})
	assert.Contains(t, b.String(), "RULE_1")
}
