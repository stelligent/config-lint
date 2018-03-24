package assertion

import "fmt"

// LoggingFunction used to control verbosity of output
type LoggingFunction func(string)

// MakeLogger returns a logging function with appropriate level of logging
func MakeLogger(verbose bool) LoggingFunction {
	if verbose {
		return func(message string) {
			fmt.Println(message)
		}
	}
	return func(message string) {}
}
