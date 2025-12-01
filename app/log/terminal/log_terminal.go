package terminal

import "fmt"

type TerminalLog struct{}

func NewTerminalLog() *TerminalLog {
	return &TerminalLog{}
}

func (t *TerminalLog) Info(message string, data any) {
	fmt.Printf("INFO: %s - %v\n", message, data)
}

func (t *TerminalLog) Warning(message string, data any) {
	fmt.Printf("WARNING: %s - %v\n", message, data)
}

func (t *TerminalLog) Error(message string, data any) {
	fmt.Printf("ERROR: %s - %v\n", message, data)
}
