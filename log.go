package main

import (
	"fmt"
	"time"
)

var (
	name = "HELIX-RELAYER-RUNNER"
)

func out(message string, a ...interface{}) {
	if len(a) != 0 {
		message = fmt.Sprintf(message, a...)
	}
	fmt.Printf("%s %s", time.Now().Format("2006-01-02 15:04:05"), message)
}

func Debug(message string) {
	out(fmt.Sprintf("[%s] DEBUG: %s", name, message))
}

func Debugf(format string, a ...interface{}) {
	out(fmt.Sprintf("[%s] DEBUG: %s", name, format), a...)
}

func Info(message string) {
	fmt.Println(fmt.Sprintf("[%s] INFO: %s", name, message))
}

func Infof(format string, a ...interface{}) {
	out(fmt.Sprintf("[%s] INFO: %s", name, format), a...)
}

func Warn(message string) {
	fmt.Println(fmt.Sprintf("[%s] WARN: %s", name, message))
}

func Warnf(format string, a ...interface{}) {
	out(fmt.Sprintf("[%s] WARN: %s", name, format), a...)
}

func Error(message string) {
	fmt.Println(fmt.Sprintf("[%s] ERROR: %s", name, message))
}

func Errorf(format string, a ...interface{}) {
	out(fmt.Sprintf("[%s] ERROR: %s", name, format), a...)
}

func Panic(message string) {
	fmt.Println(fmt.Sprintf("[%s] PANIC: %s", name, message))
	panic(message)
}

func Panicf(format string, a ...interface{}) {
	message := fmt.Sprintf(format, a...)
	fmt.Println(fmt.Sprintf("[%s] PANIC: %s", name, message))
	panic(message)
}
