package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"time"
)

func LogSetupAndDestruct() func() {
	var now string
	now = time.Now().Format("2006-01-02TT150405")
	fn := "./log-" + now + ".txt"
	logFile, err := os.OpenFile(fn, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0666)
	if err != nil {
		log.Panicln(err)
	}

	log.SetOutput(io.MultiWriter(os.Stderr, logFile))

	return func() {
		e := logFile.Close()
		if e != nil {
			fmt.Fprintf(os.Stderr, "Problem closing the log file: %s\n", e)
		}
	}
}
