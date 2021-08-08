package ui

import "github.com/buger/goterm"

type TerminalConfig struct {
	Rows    int
	Columns int
}

func NewTerminalConfig() *TerminalConfig {
	return &TerminalConfig{
		Rows:    goterm.Height(),
		Columns: goterm.Width(),
	}
}
