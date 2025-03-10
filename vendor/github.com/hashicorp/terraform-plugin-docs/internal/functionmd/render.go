// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package functionmd

import (
	"bytes"
	"fmt"
	"strings"

	tfjson "github.com/hashicorp/terraform-json"

	"github.com/hashicorp/terraform-plugin-docs/internal/schemamd"
)

// RenderArguments returns a Markdown formatted string of the function arguments.
func RenderArguments(signature *tfjson.FunctionSignature) (string, error) {
	argBuffer := bytes.NewBuffer(nil)
	for i, p := range signature.Parameters {
		name := p.Name
		desc := strings.TrimSpace(p.Description)

		typeBuffer := bytes.NewBuffer(nil)
		err := schemamd.WriteType(typeBuffer, p.Type)
		if err != nil {
			return "", err
		}

		if p.IsNullable {
			argBuffer.WriteString(fmt.Sprintf("1. `%s` (%s, Nullable) %s", name, typeBuffer.String(), desc))
		} else {
			argBuffer.WriteString(fmt.Sprintf("1. `%s` (%s) %s", name, typeBuffer.String(), desc))
		}

		if i != len(signature.Parameters)-1 {
			argBuffer.WriteString("\n")
		}

	}
	return argBuffer.String(), nil

}

// RenderSignature returns a Markdown formatted string of the function signature.
func RenderSignature(funcName string, signature *tfjson.FunctionSignature) (string, error) {

	returnType := signature.ReturnType.FriendlyName()

	paramBuffer := bytes.NewBuffer(nil)
	for i, p := range signature.Parameters {
		if i != 0 {
			paramBuffer.WriteString(", ")
		}

		paramBuffer.WriteString(fmt.Sprintf("%s %s", p.Name, p.Type.FriendlyName()))
	}

	if signature.VariadicParameter != nil {
		if signature.Parameters != nil {
			paramBuffer.WriteString(", ")
		}

		paramBuffer.WriteString(fmt.Sprintf("%s %s...", signature.VariadicParameter.Name,
			signature.VariadicParameter.Type.FriendlyName()))

	}

	return fmt.Sprintf("```text\n"+
		"%s(%s) %s\n"+
		"```",
		funcName, paramBuffer.String(), returnType), nil
}

// RenderVariadicArg returns a Markdown formatted string of the variadic argument if it exists,
// otherwise an empty string.
func RenderVariadicArg(signature *tfjson.FunctionSignature) (string, error) {
	if signature.VariadicParameter == nil {
		return "", nil
	}

	name := signature.VariadicParameter.Name
	desc := strings.TrimSpace(signature.VariadicParameter.Description)

	typeBuffer := bytes.NewBuffer(nil)
	err := schemamd.WriteType(typeBuffer, signature.VariadicParameter.Type)
	if err != nil {
		return "", err
	}

	if signature.VariadicParameter.IsNullable {
		return fmt.Sprintf("1. `%s` (Variadic, %s, Nullable) %s", name, typeBuffer.String(), desc), nil
	} else {
		return fmt.Sprintf("1. `%s` (Variadic, %s) %s", name, typeBuffer.String(), desc), nil
	}

}
