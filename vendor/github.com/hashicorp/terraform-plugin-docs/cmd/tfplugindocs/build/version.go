// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package build

var (
	// These vars will be set by goreleaser.
	version string = `dev`
	commit  string = ``
)

func GetVersion() string {
	version := "tfplugindocs" + " Version " + version
	if commit != "" {
		version += " from commit " + commit
	}
	return version
}
