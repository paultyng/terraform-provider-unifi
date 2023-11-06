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

package compose

import (
	"fmt"
	"sync/atomic"

	"github.com/docker/compose/v2/pkg/api"
)

// logPrinter watch application containers an collect their logs
type logPrinter interface {
	HandleEvent(event api.ContainerEvent)
	Run(cascadeStop bool, exitCodeFrom string, stopFn func() error) (int, error)
	Cancel()
	Stop()
}

type printer struct {
	queue    chan api.ContainerEvent
	consumer api.LogConsumer
	stopped  atomic.Bool
}

// newLogPrinter builds a LogPrinter passing containers logs to LogConsumer
func newLogPrinter(consumer api.LogConsumer) logPrinter {
	queue := make(chan api.ContainerEvent)
	printer := printer{
		consumer: consumer,
		queue:    queue,
	}
	return &printer
}

func (p *printer) Cancel() {
	// note: HandleEvent is used to ensure this doesn't deadlock
	p.HandleEvent(api.ContainerEvent{Type: api.UserCancel})
}

func (p *printer) Stop() {
	if p.stopped.CompareAndSwap(false, true) {
		// only close if this is the first call to stop
		close(p.queue)
	}
}

func (p *printer) HandleEvent(event api.ContainerEvent) {
	// prevent deadlocking, if the printer is done, there's no reader for
	// queue, so this write could block indefinitely
	if p.stopped.Load() {
		return
	}
	p.queue <- event
}

//nolint:gocyclo
func (p *printer) Run(cascadeStop bool, exitCodeFrom string, stopFn func() error) (int, error) {
	var (
		aborting bool
		exitCode int
	)
	containers := map[string]struct{}{}
	for event := range p.queue {
		container, id := event.Container, event.ID
		switch event.Type {
		case api.UserCancel:
			aborting = true
		case api.ContainerEventAttach:
			if _, ok := containers[id]; ok {
				continue
			}
			containers[id] = struct{}{}
			p.consumer.Register(container)
		case api.ContainerEventExit, api.ContainerEventStopped, api.ContainerEventRecreated:
			if !event.Restarting {
				delete(containers, id)
			}
			if !aborting {
				p.consumer.Status(container, fmt.Sprintf("exited with code %d", event.ExitCode))
				if event.Type == api.ContainerEventRecreated {
					p.consumer.Status(container, "has been recreated")
				}
			}
			if cascadeStop {
				if !aborting {
					aborting = true
					err := stopFn()
					if err != nil {
						return 0, err
					}
				}
				if event.Type == api.ContainerEventExit {
					if exitCodeFrom == "" {
						exitCodeFrom = event.Service
					}
					if exitCodeFrom == event.Service {
						exitCode = event.ExitCode
					}
				}
			}
			if len(containers) == 0 {
				// Last container terminated, done
				return exitCode, nil
			}
		case api.ContainerEventLog:
			if !aborting {
				p.consumer.Log(container, event.Line)
			}
		case api.ContainerEventErr:
			if !aborting {
				p.consumer.Err(container, event.Line)
			}
		}
	}
	return exitCode, nil
}
