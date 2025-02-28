// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package main

import (
	"os"

	"github.com/hashicorp/terraform-plugin-docs/internal/cmd"
)

func main() {
	os.Exit(cmd.Main())
}
