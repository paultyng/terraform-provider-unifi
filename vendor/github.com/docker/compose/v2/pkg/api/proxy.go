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

package api

import (
	"context"

	"github.com/compose-spec/compose-go/types"
)

var _ Service = &ServiceProxy{}

// ServiceProxy implements Service by delegating to implementation functions. This allows lazy init and per-method overrides
type ServiceProxy struct {
	BuildFn              func(ctx context.Context, project *types.Project, options BuildOptions) error
	PushFn               func(ctx context.Context, project *types.Project, options PushOptions) error
	PullFn               func(ctx context.Context, project *types.Project, opts PullOptions) error
	CreateFn             func(ctx context.Context, project *types.Project, opts CreateOptions) error
	StartFn              func(ctx context.Context, projectName string, options StartOptions) error
	RestartFn            func(ctx context.Context, projectName string, options RestartOptions) error
	StopFn               func(ctx context.Context, projectName string, options StopOptions) error
	UpFn                 func(ctx context.Context, project *types.Project, options UpOptions) error
	DownFn               func(ctx context.Context, projectName string, options DownOptions) error
	LogsFn               func(ctx context.Context, projectName string, consumer LogConsumer, options LogOptions) error
	PsFn                 func(ctx context.Context, projectName string, options PsOptions) ([]ContainerSummary, error)
	ListFn               func(ctx context.Context, options ListOptions) ([]Stack, error)
	ConfigFn             func(ctx context.Context, project *types.Project, options ConfigOptions) ([]byte, error)
	KillFn               func(ctx context.Context, project string, options KillOptions) error
	RunOneOffContainerFn func(ctx context.Context, project *types.Project, opts RunOptions) (int, error)
	RemoveFn             func(ctx context.Context, project string, options RemoveOptions) error
	ExecFn               func(ctx context.Context, project string, opts RunOptions) (int, error)
	CopyFn               func(ctx context.Context, project string, options CopyOptions) error
	PauseFn              func(ctx context.Context, project string, options PauseOptions) error
	UnPauseFn            func(ctx context.Context, project string, options PauseOptions) error
	TopFn                func(ctx context.Context, projectName string, services []string) ([]ContainerProcSummary, error)
	EventsFn             func(ctx context.Context, project string, options EventsOptions) error
	PortFn               func(ctx context.Context, project string, service string, port uint16, options PortOptions) (string, int, error)
	ImagesFn             func(ctx context.Context, projectName string, options ImagesOptions) ([]ImageSummary, error)
	WatchFn              func(ctx context.Context, project *types.Project, services []string, options WatchOptions) error
	MaxConcurrencyFn     func(parallel int)
	DryRunModeFn         func(ctx context.Context, dryRun bool) (context.Context, error)
	VizFn                func(ctx context.Context, project *types.Project, options VizOptions) (string, error)
	WaitFn               func(ctx context.Context, projectName string, options WaitOptions) (int64, error)
	PublishFn            func(ctx context.Context, project *types.Project, repository string, options PublishOptions) error
	ScaleFn              func(ctx context.Context, project *types.Project, options ScaleOptions) error
	interceptors         []Interceptor
}

// NewServiceProxy produces a ServiceProxy
func NewServiceProxy() *ServiceProxy {
	return &ServiceProxy{}
}

// Interceptor allow to customize the compose types.Project before the actual Service method is executed
type Interceptor func(ctx context.Context, project *types.Project)

// WithService configure proxy to use specified Service as delegate
func (s *ServiceProxy) WithService(service Service) *ServiceProxy {
	s.BuildFn = service.Build
	s.PushFn = service.Push
	s.PullFn = service.Pull
	s.CreateFn = service.Create
	s.StartFn = service.Start
	s.RestartFn = service.Restart
	s.StopFn = service.Stop
	s.UpFn = service.Up
	s.DownFn = service.Down
	s.LogsFn = service.Logs
	s.PsFn = service.Ps
	s.ListFn = service.List
	s.ConfigFn = service.Config
	s.KillFn = service.Kill
	s.RunOneOffContainerFn = service.RunOneOffContainer
	s.RemoveFn = service.Remove
	s.ExecFn = service.Exec
	s.CopyFn = service.Copy
	s.PauseFn = service.Pause
	s.UnPauseFn = service.UnPause
	s.TopFn = service.Top
	s.EventsFn = service.Events
	s.PortFn = service.Port
	s.PublishFn = service.Publish
	s.ImagesFn = service.Images
	s.WatchFn = service.Watch
	s.MaxConcurrencyFn = service.MaxConcurrency
	s.DryRunModeFn = service.DryRunMode
	s.VizFn = service.Viz
	s.WaitFn = service.Wait
	s.ScaleFn = service.Scale
	return s
}

// WithInterceptor configures Interceptor to be applied to Service method execution
func (s *ServiceProxy) WithInterceptor(interceptors ...Interceptor) *ServiceProxy {
	s.interceptors = append(s.interceptors, interceptors...)
	return s
}

// Build implements Service interface
func (s *ServiceProxy) Build(ctx context.Context, project *types.Project, options BuildOptions) error {
	if s.BuildFn == nil {
		return ErrNotImplemented
	}
	for _, i := range s.interceptors {
		i(ctx, project)
	}
	return s.BuildFn(ctx, project, options)
}

// Push implements Service interface
func (s *ServiceProxy) Push(ctx context.Context, project *types.Project, options PushOptions) error {
	if s.PushFn == nil {
		return ErrNotImplemented
	}
	for _, i := range s.interceptors {
		i(ctx, project)
	}
	return s.PushFn(ctx, project, options)
}

// Pull implements Service interface
func (s *ServiceProxy) Pull(ctx context.Context, project *types.Project, options PullOptions) error {
	if s.PullFn == nil {
		return ErrNotImplemented
	}
	for _, i := range s.interceptors {
		i(ctx, project)
	}
	return s.PullFn(ctx, project, options)
}

// Create implements Service interface
func (s *ServiceProxy) Create(ctx context.Context, project *types.Project, options CreateOptions) error {
	if s.CreateFn == nil {
		return ErrNotImplemented
	}
	for _, i := range s.interceptors {
		i(ctx, project)
	}
	return s.CreateFn(ctx, project, options)
}

// Start implements Service interface
func (s *ServiceProxy) Start(ctx context.Context, projectName string, options StartOptions) error {
	if s.StartFn == nil {
		return ErrNotImplemented
	}
	return s.StartFn(ctx, projectName, options)
}

// Restart implements Service interface
func (s *ServiceProxy) Restart(ctx context.Context, projectName string, options RestartOptions) error {
	if s.RestartFn == nil {
		return ErrNotImplemented
	}
	return s.RestartFn(ctx, projectName, options)
}

// Stop implements Service interface
func (s *ServiceProxy) Stop(ctx context.Context, projectName string, options StopOptions) error {
	if s.StopFn == nil {
		return ErrNotImplemented
	}
	return s.StopFn(ctx, projectName, options)
}

// Up implements Service interface
func (s *ServiceProxy) Up(ctx context.Context, project *types.Project, options UpOptions) error {
	if s.UpFn == nil {
		return ErrNotImplemented
	}
	for _, i := range s.interceptors {
		i(ctx, project)
	}
	return s.UpFn(ctx, project, options)
}

// Down implements Service interface
func (s *ServiceProxy) Down(ctx context.Context, project string, options DownOptions) error {
	if s.DownFn == nil {
		return ErrNotImplemented
	}
	return s.DownFn(ctx, project, options)
}

// Logs implements Service interface
func (s *ServiceProxy) Logs(ctx context.Context, projectName string, consumer LogConsumer, options LogOptions) error {
	if s.LogsFn == nil {
		return ErrNotImplemented
	}
	return s.LogsFn(ctx, projectName, consumer, options)
}

// Ps implements Service interface
func (s *ServiceProxy) Ps(ctx context.Context, project string, options PsOptions) ([]ContainerSummary, error) {
	if s.PsFn == nil {
		return nil, ErrNotImplemented
	}
	return s.PsFn(ctx, project, options)
}

// List implements Service interface
func (s *ServiceProxy) List(ctx context.Context, options ListOptions) ([]Stack, error) {
	if s.ListFn == nil {
		return nil, ErrNotImplemented
	}
	return s.ListFn(ctx, options)
}

// Config implements Service interface
func (s *ServiceProxy) Config(ctx context.Context, project *types.Project, options ConfigOptions) ([]byte, error) {
	if s.ConfigFn == nil {
		return nil, ErrNotImplemented
	}
	for _, i := range s.interceptors {
		i(ctx, project)
	}
	return s.ConfigFn(ctx, project, options)
}

// Kill implements Service interface
func (s *ServiceProxy) Kill(ctx context.Context, projectName string, options KillOptions) error {
	if s.KillFn == nil {
		return ErrNotImplemented
	}
	return s.KillFn(ctx, projectName, options)
}

// RunOneOffContainer implements Service interface
func (s *ServiceProxy) RunOneOffContainer(ctx context.Context, project *types.Project, options RunOptions) (int, error) {
	if s.RunOneOffContainerFn == nil {
		return 0, ErrNotImplemented
	}
	for _, i := range s.interceptors {
		i(ctx, project)
	}
	return s.RunOneOffContainerFn(ctx, project, options)
}

// Remove implements Service interface
func (s *ServiceProxy) Remove(ctx context.Context, projectName string, options RemoveOptions) error {
	if s.RemoveFn == nil {
		return ErrNotImplemented
	}
	return s.RemoveFn(ctx, projectName, options)
}

// Exec implements Service interface
func (s *ServiceProxy) Exec(ctx context.Context, projectName string, options RunOptions) (int, error) {
	if s.ExecFn == nil {
		return 0, ErrNotImplemented
	}
	return s.ExecFn(ctx, projectName, options)
}

// Copy implements Service interface
func (s *ServiceProxy) Copy(ctx context.Context, projectName string, options CopyOptions) error {
	if s.CopyFn == nil {
		return ErrNotImplemented
	}
	return s.CopyFn(ctx, projectName, options)
}

// Pause implements Service interface
func (s *ServiceProxy) Pause(ctx context.Context, projectName string, options PauseOptions) error {
	if s.PauseFn == nil {
		return ErrNotImplemented
	}
	return s.PauseFn(ctx, projectName, options)
}

// UnPause implements Service interface
func (s *ServiceProxy) UnPause(ctx context.Context, projectName string, options PauseOptions) error {
	if s.UnPauseFn == nil {
		return ErrNotImplemented
	}
	return s.UnPauseFn(ctx, projectName, options)
}

// Top implements Service interface
func (s *ServiceProxy) Top(ctx context.Context, project string, services []string) ([]ContainerProcSummary, error) {
	if s.TopFn == nil {
		return nil, ErrNotImplemented
	}
	return s.TopFn(ctx, project, services)
}

// Events implements Service interface
func (s *ServiceProxy) Events(ctx context.Context, projectName string, options EventsOptions) error {
	if s.EventsFn == nil {
		return ErrNotImplemented
	}
	return s.EventsFn(ctx, projectName, options)
}

// Port implements Service interface
func (s *ServiceProxy) Port(ctx context.Context, projectName string, service string, port uint16, options PortOptions) (string, int, error) {
	if s.PortFn == nil {
		return "", 0, ErrNotImplemented
	}
	return s.PortFn(ctx, projectName, service, port, options)
}

func (s *ServiceProxy) Publish(ctx context.Context, project *types.Project, repository string, options PublishOptions) error {
	return s.PublishFn(ctx, project, repository, options)
}

// Images implements Service interface
func (s *ServiceProxy) Images(ctx context.Context, project string, options ImagesOptions) ([]ImageSummary, error) {
	if s.ImagesFn == nil {
		return nil, ErrNotImplemented
	}
	return s.ImagesFn(ctx, project, options)
}

// Watch implements Service interface
func (s *ServiceProxy) Watch(ctx context.Context, project *types.Project, services []string, options WatchOptions) error {
	if s.WatchFn == nil {
		return ErrNotImplemented
	}
	return s.WatchFn(ctx, project, services, options)
}

// Viz implements Service interface
func (s *ServiceProxy) Viz(ctx context.Context, project *types.Project, options VizOptions) (string, error) {
	if s.VizFn == nil {
		return "", ErrNotImplemented
	}
	return s.VizFn(ctx, project, options)
}

// Wait implements Service interface
func (s *ServiceProxy) Wait(ctx context.Context, projectName string, options WaitOptions) (int64, error) {
	if s.WaitFn == nil {
		return 0, ErrNotImplemented
	}
	return s.WaitFn(ctx, projectName, options)
}

func (s *ServiceProxy) Scale(ctx context.Context, project *types.Project, options ScaleOptions) error {
	if s.ScaleFn == nil {
		return ErrNotImplemented
	}
	return s.ScaleFn(ctx, project, options)
}

func (s *ServiceProxy) MaxConcurrency(i int) {
	s.MaxConcurrencyFn(i)
}

func (s *ServiceProxy) DryRunMode(ctx context.Context, dryRun bool) (context.Context, error) {
	return s.DryRunModeFn(ctx, dryRun)
}
