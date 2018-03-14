package assertion

import "fmt"

type LoggingFunction func(string)

func MakeLogger(verbose bool) LoggingFunction {
	if verbose {
		return func(message string) {
			fmt.Println(message)
		}
	}
	return func(message string) {}
}
