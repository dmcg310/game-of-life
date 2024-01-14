package main

import (
	"fmt"
	"os"
)

type AppError struct {
	Error    error
	Message  string
	Solution string // potential fix for the error
}

type AppWarning struct {
	Message  string
	Solution string // potential fix for the warning
}

func newAppError(error error, message string, solution string) *AppError {
	return &AppError{
		Error:    error,
		Message:  message,
		Solution: solution,
	}
}

func newAppWarning(message string, solution string) *AppWarning {
	return &AppWarning{
		Message:  message,
		Solution: solution,
	}
}

func (e *AppError) showAppErrorFatal() {
	msg := "\x1b[31m[ERROR]\x1b[0m\n%s\n%s\n%s"
	err := fmt.Sprintf(msg, e.Error, e.Message, e.Solution)
	fmt.Println(err)
	os.Exit(1)
}

func (w *AppWarning) showAppWarning() {
	msg := "\x1b[93m[WARNING]\x1b[0m\n%s\n%s"
	warning := fmt.Sprintf(msg, w.Message, w.Solution)
	fmt.Println(warning)
}
