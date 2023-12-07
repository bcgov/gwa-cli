package pkg

import (
	"bytes"
	"fmt"
	"log"
)

var (
	buf         bytes.Buffer
	ErrorLogger *log.Logger
	InfoLogger  *log.Logger
)

func Error(msg string) {
	ErrorLogger.Println(msg)
}

func Info(msg string) {
	InfoLogger.Println(msg)
}

func PrintLog() {
	fmt.Print(buf.String())
}

func init() {
	InfoLogger = log.New(&buf, WarningStyle.Render("[INFO] "), log.Ltime)
	ErrorLogger = log.New(&buf, ErrorStyle.Render("[ERROR] "), log.Ltime)
}
