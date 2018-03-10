package main

import "fmt"

type LoggingFunction func(string)

func makeLogger(verbose bool) LoggingFunction {
	if verbose {
		return func(message string) {
			fmt.Println(message)
		}
	}
	return func(message string) {}
}
