// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package tmplfuncs

import (
	"fmt"
	"os"
	"strings"
)

func PrefixLines(prefix, text string) string {
	return prefix + strings.Join(strings.Split(text, "\n"), "\n"+prefix)
}

func CodeFile(format, file string) (string, error) {
	content, err := os.ReadFile(file)
	if err != nil {
		return "", fmt.Errorf("unable to read content from %q: %w", file, err)
	}

	sContent := strings.TrimSpace(string(content))
	if sContent == "" {
		return "", fmt.Errorf("no file content in %q", file)
	}

	md := &strings.Builder{}
	_, err = md.WriteString("```" + format + "\n")
	if err != nil {
		return "", err
	}

	_, err = md.WriteString(sContent)
	if err != nil {
		return "", err
	}

	_, err = md.WriteString("\n```")
	if err != nil {
		return "", err
	}

	return md.String(), nil
}
