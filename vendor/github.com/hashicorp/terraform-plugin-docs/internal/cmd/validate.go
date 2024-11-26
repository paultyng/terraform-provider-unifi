// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package cmd

import (
	"errors"
	"flag"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-docs/internal/provider"
)

type validateCmd struct {
	commonCmd

	flagProviderName    string
	flagProviderDir     string
	flagProvidersSchema string
	tfVersion           string
}

func (cmd *validateCmd) Synopsis() string {
	return "validates a plugin website"
}

func (cmd *validateCmd) Help() string {
	strBuilder := &strings.Builder{}

	longestName := 0
	longestUsage := 0
	cmd.Flags().VisitAll(func(f *flag.Flag) {
		if len(f.Name) > longestName {
			longestName = len(f.Name)
		}
		if len(f.Usage) > longestUsage {
			longestUsage = len(f.Usage)
		}
	})

	strBuilder.WriteString("\nUsage: tfplugindocs validate [<args>]\n\n")
	cmd.Flags().VisitAll(func(f *flag.Flag) {
		if f.DefValue != "" {
			strBuilder.WriteString(fmt.Sprintf("    --%s <ARG> %s%s%s  (default: %q)\n",
				f.Name,
				strings.Repeat(" ", longestName-len(f.Name)+2),
				f.Usage,
				strings.Repeat(" ", longestUsage-len(f.Usage)+2),
				f.DefValue,
			))
		} else {
			strBuilder.WriteString(fmt.Sprintf("    --%s <ARG> %s%s%s\n",
				f.Name,
				strings.Repeat(" ", longestName-len(f.Name)+2),
				f.Usage,
				strings.Repeat(" ", longestUsage-len(f.Usage)+2),
			))
		}
	})
	strBuilder.WriteString("\n")

	return strBuilder.String()
}

func (cmd *validateCmd) Flags() *flag.FlagSet {
	fs := flag.NewFlagSet("validate", flag.ExitOnError)
	fs.StringVar(&cmd.flagProviderName, "provider-name", "", "provider name, as used in Terraform configurations; defaults to the --provider-dir short name (after removing `terraform-provider-` prefix)")
	fs.StringVar(&cmd.flagProviderDir, "provider-dir", "", "relative or absolute path to the root provider code directory; this will default to the current working directory if not set")
	fs.StringVar(&cmd.flagProvidersSchema, "providers-schema", "", "path to the providers schema JSON file, which contains the output of the terraform providers schema -json command. Setting this flag will skip building the provider and calling Terraform CLI")
	fs.StringVar(&cmd.tfVersion, "tf-version", "", "terraform binary version to download. If not provided, will look for a terraform binary in the local environment. If not found in the environment, will download the latest version of Terraform")
	return fs
}

func (cmd *validateCmd) Run(args []string) int {
	fs := cmd.Flags()
	err := fs.Parse(args)
	if err != nil {
		cmd.ui.Error(fmt.Sprintf("unable to parse flags: %s", err))
		return 1
	}

	return cmd.run(cmd.runInternal)
}

func (cmd *validateCmd) runInternal() error {
	err := provider.Validate(cmd.ui,
		cmd.flagProviderDir,
		cmd.flagProviderName,
		cmd.flagProvidersSchema,
		cmd.tfVersion,
	)
	if err != nil {
		return errors.Join(errors.New("validation errors found: "), err)
	}

	return nil
}
