package main

import (
	"fmt"
	"github.com/stelligent/config-lint/assertion"
)

func (w DefaultReportWriter) WriteReport(report assertion.ValidationReport, options LinterOptions) {
	if options.SearchExpression == "" {
		err := printReport(w.Writer, report, options.QueryExpression)
		if err != nil {
			fmt.Println(err.Error())
		}
	}
}
