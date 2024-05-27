package logger

import (
	"log"
	"os"
)

var Logger *log.Logger
var LogFile *os.File

// LoggerInit initializes the logger by opening a file named "logger.txt" with the necessary permissions.
//
// Returns a pointer to the opened file.
func LoggerInit() *os.File {
	file, err := os.OpenFile("logger.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal(err)
	}
	return file
}