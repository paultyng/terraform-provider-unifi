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

package tracing

import (
	"strings"
	"time"

	"github.com/docker/compose/v2/pkg/utils"

	"github.com/compose-spec/compose-go/types"
	moby "github.com/docker/docker/api/types"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

// SpanOptions is a small helper type to make it easy to share the options helpers between
// downstream functions that accept slices of trace.SpanStartOption and trace.EventOption.
type SpanOptions []trace.SpanStartEventOption

func (s SpanOptions) SpanStartOptions() []trace.SpanStartOption {
	out := make([]trace.SpanStartOption, len(s))
	for i := range s {
		out[i] = s[i]
	}
	return out
}

func (s SpanOptions) EventOptions() []trace.EventOption {
	out := make([]trace.EventOption, len(s))
	for i := range s {
		out[i] = s[i]
	}
	return out
}

// ProjectOptions returns common attributes from a Compose project.
//
// For convenience, it's returned as a SpanOptions object to allow it to be
// passed directly to the wrapping helper methods in this package such as
// SpanWrapFunc.
func ProjectOptions(proj *types.Project) SpanOptions {
	if proj == nil {
		return nil
	}

	disabledServiceNames := make([]string, len(proj.DisabledServices))
	for i := range proj.DisabledServices {
		disabledServiceNames[i] = proj.DisabledServices[i].Name
	}

	attrs := []attribute.KeyValue{
		attribute.String("project.name", proj.Name),
		attribute.String("project.dir", proj.WorkingDir),
		attribute.StringSlice("project.compose_files", proj.ComposeFiles),
		attribute.StringSlice("project.services.active", proj.ServiceNames()),
		attribute.StringSlice("project.services.disabled", disabledServiceNames),
		attribute.StringSlice("project.profiles", proj.Profiles),
		attribute.StringSlice("project.volumes", proj.VolumeNames()),
		attribute.StringSlice("project.networks", proj.NetworkNames()),
		attribute.StringSlice("project.secrets", proj.SecretNames()),
		attribute.StringSlice("project.configs", proj.ConfigNames()),
		attribute.StringSlice("project.extensions", keys(proj.Extensions)),
		attribute.StringSlice("project.includes", flattenIncludeReferences(proj.IncludeReferences)),
	}
	return []trace.SpanStartEventOption{
		trace.WithAttributes(attrs...),
	}
}

// ServiceOptions returns common attributes from a Compose service.
//
// For convenience, it's returned as a SpanOptions object to allow it to be
// passed directly to the wrapping helper methods in this package such as
// SpanWrapFunc.
func ServiceOptions(service types.ServiceConfig) SpanOptions {
	attrs := []attribute.KeyValue{
		attribute.String("service.name", service.Name),
		attribute.String("service.image", service.Image),
		attribute.StringSlice("service.networks", keys(service.Networks)),
	}

	configNames := make([]string, len(service.Configs))
	for i := range service.Configs {
		configNames[i] = service.Configs[i].Source
	}
	attrs = append(attrs, attribute.StringSlice("service.configs", configNames))

	secretNames := make([]string, len(service.Secrets))
	for i := range service.Secrets {
		secretNames[i] = service.Secrets[i].Source
	}
	attrs = append(attrs, attribute.StringSlice("service.secrets", secretNames))

	volNames := make([]string, len(service.Volumes))
	for i := range service.Volumes {
		volNames[i] = service.Volumes[i].Source
	}
	attrs = append(attrs, attribute.StringSlice("service.volumes", volNames))

	return []trace.SpanStartEventOption{
		trace.WithAttributes(attrs...),
	}
}

// ContainerOptions returns common attributes from a Moby container.
//
// For convenience, it's returned as a SpanOptions object to allow it to be
// passed directly to the wrapping helper methods in this package such as
// SpanWrapFunc.
func ContainerOptions(container moby.Container) SpanOptions {
	attrs := []attribute.KeyValue{
		attribute.String("container.id", container.ID),
		attribute.String("container.image", container.Image),
		unixTimeAttr("container.created_at", container.Created),
	}

	if len(container.Names) != 0 {
		attrs = append(attrs, attribute.String("container.name", strings.TrimPrefix(container.Names[0], "/")))
	}

	return []trace.SpanStartEventOption{
		trace.WithAttributes(attrs...),
	}
}

func keys[T any](m map[string]T) []string {
	out := make([]string, 0, len(m))
	for k := range m {
		out = append(out, k)
	}
	return out
}

func timeAttr(key string, value time.Time) attribute.KeyValue {
	return attribute.String(key, value.Format(time.RFC3339))
}

func unixTimeAttr(key string, value int64) attribute.KeyValue {
	return timeAttr(key, time.Unix(value, 0).UTC())
}

func flattenIncludeReferences(includeRefs map[string][]types.IncludeConfig) []string {
	ret := utils.NewSet[string]()
	for _, included := range includeRefs {
		for i := range included {
			ret.AddAll(included[i].Path...)
		}
	}
	return ret.Elements()
}
