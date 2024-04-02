package main

import (
	"fmt"
	"os"
	"sync"
	"time"
)

var (
	name = "HELIX-RELAYER-RUNNER"
	m    sync.Mutex
)

func out(message string, a ...interface{}) {
	m.Lock()
	defer m.Unlock()
	if len(a) != 0 {
		message = fmt.Sprintf(message, a...)
	}
	_, _ = fmt.Fprint(os.Stderr, fmt.Sprintf("[%s] %s %s", name, time.Now().Format("2006-01-02 15:04:05"), message))
}

func Debug(message string) {
	out(fmt.Sprintf("DEBUG: %s\n", message))
}

func Debugf(format string, a ...interface{}) {
	out(fmt.Sprintf("DEBUG: %s", format), a...)
}

func Info(message string) {
	fmt.Println(fmt.Sprintf("INFO: %s\n", message))
}

func Infof(format string, a ...interface{}) {
	out(fmt.Sprintf("INFO: %s", format), a...)
}

func Warn(message string) {
	fmt.Println(fmt.Sprintf("WARN: %s\n", message))
}

func Warnf(format string, a ...interface{}) {
	out(fmt.Sprintf("WARN: %s", format), a...)
}

func Error(message string) {
	fmt.Println(fmt.Sprintf("ERROR: %s\n", message))
}

func Errorf(format string, a ...interface{}) {
	out(fmt.Sprintf("ERROR: %s", format), a...)
}

func Panic(message string) {
	fmt.Println(fmt.Sprintf("PANIC: %s\n", name, message))
	panic(message)
}

func Panicf(format string, a ...interface{}) {
	message := fmt.Sprintf(format, a...)
	fmt.Println(fmt.Sprintf("PANIC: %s", name, message))
	panic(message)
}
