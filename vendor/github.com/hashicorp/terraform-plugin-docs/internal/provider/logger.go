// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"fmt"

	"github.com/hashicorp/cli"
)

type Logger struct {
	ui cli.Ui
}

func NewLogger(ui cli.Ui) *Logger {
	return &Logger{ui}
}

func (l *Logger) infof(format string, args ...interface{}) {
	l.ui.Info(fmt.Sprintf(format, args...))
}

//nolint:unused
func (l *Logger) warnf(format string, args ...interface{}) {
	l.ui.Warn(fmt.Sprintf(format, args...))
}
