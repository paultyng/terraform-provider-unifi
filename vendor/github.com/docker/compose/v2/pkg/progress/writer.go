/*
   Copyright 2020 Docker Compose CLI authors

   Licensed under the Apache License, Version 2.0 (the "License");
   you may not use this file except in compliance with the License.
   You may obtain a copy of the License at

       http://www.apache.org/licenses/LICENSE-2.0

   Unless required by applicable law or agreed to in writing, software
   distributed under the License is distributed on an "AS IS" BASIS,
   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
   See the License for the specific language governing permissions and
   limitations under the License.
*/

package progress

import (
	"context"
	"io"
	"sync"

	"github.com/containerd/console"
	"github.com/moby/term"
	"github.com/sirupsen/logrus"
	"golang.org/x/sync/errgroup"

	"github.com/docker/compose/v2/pkg/api"
)

// Writer can write multiple progress events
type Writer interface {
	Start(context.Context) error
	Stop()
	Event(Event)
	Events([]Event)
	TailMsgf(string, ...interface{})
}

type writerKey struct{}

// WithContextWriter adds the writer to the context
func WithContextWriter(ctx context.Context, writer Writer) context.Context {
	return context.WithValue(ctx, writerKey{}, writer)
}

// ContextWriter returns the writer from the context
func ContextWriter(ctx context.Context) Writer {
	s, ok := ctx.Value(writerKey{}).(Writer)
	if !ok {
		return &noopWriter{}
	}
	return s
}

type progressFunc func(context.Context) error

type progressFuncWithStatus func(context.Context) (string, error)

// Run will run a writer and the progress function in parallel
func Run(ctx context.Context, pf progressFunc, out io.Writer) error {
	_, err := RunWithStatus(ctx, func(ctx context.Context) (string, error) {
		return "", pf(ctx)
	}, out, "Running")
	return err
}

func RunWithTitle(ctx context.Context, pf progressFunc, out io.Writer, progressTitle string) error {
	_, err := RunWithStatus(ctx, func(ctx context.Context) (string, error) {
		return "", pf(ctx)
	}, out, progressTitle)
	return err
}

// RunWithStatus will run a writer and the progress function in parallel and return a status
func RunWithStatus(ctx context.Context, pf progressFuncWithStatus, out io.Writer, progressTitle string) (string, error) {
	eg, _ := errgroup.WithContext(ctx)
	w, err := NewWriter(ctx, out, progressTitle)
	var result string
	if err != nil {
		return "", err
	}
	eg.Go(func() error {
		return w.Start(context.Background())
	})

	ctx = WithContextWriter(ctx, w)

	eg.Go(func() error {
		defer w.Stop()
		s, err := pf(ctx)
		if err == nil {
			result = s
		}
		return err
	})

	err = eg.Wait()
	return result, err
}

const (
	// ModeAuto detect console capabilities
	ModeAuto = "auto"
	// ModeTTY use terminal capability for advanced rendering
	ModeTTY = "tty"
	// ModePlain dump raw events to output
	ModePlain = "plain"
	// ModeQuiet don't display events
	ModeQuiet = "quiet"
)

// Mode define how progress should be rendered, either as ModePlain or ModeTTY
var Mode = ModeAuto

// NewWriter returns a new multi-progress writer
func NewWriter(ctx context.Context, out io.Writer, progressTitle string) (Writer, error) {
	_, isTerminal := term.GetFdInfo(out)
	dryRun, ok := ctx.Value(api.DryRunKey{}).(bool)
	if !ok {
		dryRun = false
	}
	if Mode == ModeQuiet {
		return quiet{}, nil
	}
	f, isConsole := out.(console.File) // see https://github.com/docker/compose/issues/10560
	if Mode == ModeAuto && isTerminal && isConsole {
		return newTTYWriter(f, dryRun, progressTitle)
	}
	if Mode == ModeTTY {
		if !isConsole {
			logrus.Warn("Terminal is not a POSIX console")
		} else {
			return newTTYWriter(f, dryRun, progressTitle)
		}
	}
	return &plainWriter{
		out:    out,
		done:   make(chan bool),
		dryRun: dryRun,
	}, nil
}

func newTTYWriter(out console.File, dryRun bool, progressTitle string) (Writer, error) {
	con, err := console.ConsoleFromFile(out)
	if err != nil {
		return nil, err
	}

	return &ttyWriter{
		out:           con,
		eventIDs:      []string{},
		events:        map[string]Event{},
		repeated:      false,
		done:          make(chan bool),
		mtx:           &sync.Mutex{},
		dryRun:        dryRun,
		progressTitle: progressTitle,
	}, nil
}
