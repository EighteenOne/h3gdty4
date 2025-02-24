package logger

import (
	"fmt"
	"log"
)

var logger = log.Default()

func Errorf(format string, v ...any) {
	msg := fmt.Sprintf(format, v...)
	logger.Printf("[Error]: %s\n", msg)
}

func Info(format string, v ...any) {
	msg := fmt.Sprintf(format, v...)
	logger.Printf("[Info]: %s\n", msg)
}

func Fatalf(format string, v ...any) {
	msg := fmt.Sprintf(format, v...)
	logger.Fatalf("[Fatal]: %s\n", msg)
}
