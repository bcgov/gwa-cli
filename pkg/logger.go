// Simple logging wrapper around Go's std.log module.
//
// Broken into standalone module so implementation can be easily swapped out with another
// library, if needed.
package pkg

import (
	"bytes"
	"fmt"
	"log"
)

var (
	buf           bytes.Buffer
	ErrorLogger   *log.Logger
	InfoLogger    *log.Logger
	WarningLogger *log.Logger
)

// Prints `[ERROR] 00:00:00 message`
func Error(msg string) {
	ErrorLogger.Println(msg)
}

// Prints `[INFO] 00:00:00 message`
func Info(msg string) {
	InfoLogger.Println(msg)
}

// Prints `[WARN] 00:00:00 message`
func Warning(msg string) {
	WarningLogger.Println(msg)
}

// Prints logger buffer to stderr once rest of the program has finished running
// NOTE: for this to display when a Cobra command throws an error, be sure to wrap
// the `RunE` command in the `WrapError` higher order function in [pkg/context]
func PrintLog() {
	fmt.Print(buf.String())
}

func init() {
	InfoLogger = log.New(&buf, InfoStyle.Render("[INFO] "), log.Ltime)
	WarningLogger = log.New(&buf, WarningStyle.Render("[WARN] "), log.Ltime)
	ErrorLogger = log.New(&buf, ErrorStyle.Render("[ERROR] "), log.Ltime)
}
