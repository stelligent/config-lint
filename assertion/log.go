package assertion

import "fmt"

var (
	isVerbose = false
)

// SetVerbose turns verbose logging on or off
func SetVerbose(b bool) {
	isVerbose = b
}

// Debugf prints a formatted string when verbose logging is turned on
func Debugf(format string, args ...interface{}) {
	if isVerbose == false {
		return
	}
	fmt.Printf(format, args...)
}
