package sdkparse

import (
	"fmt"
	"go/types"

	"github.com/paultyng/tfplugingen/config"
)

func schemaFromType(ty types.Type) (*config.Schema, error) {
	switch ty := ty.(type) {
	case *types.Basic:
		switch ty.Kind() {
		case types.Bool:
			sty := "boolean"
			return &config.Schema{
				Type: &sty,
			}, nil
		case types.Int, types.Int8, types.Int16, types.Int32, types.Int64,
			types.Uint, types.Uint8, types.Uint16, types.Uint32, types.Uint64:
			sty := "integer"
			return &config.Schema{
				Type: &sty,
			}, nil
		case types.Float32, types.Float64:
			sty := "number"
			return &config.Schema{
				Type: &sty,
			}, nil
		case types.String:
			sty := "string"
			return &config.Schema{
				Type: &sty,
			}, nil
		}
	case *types.Named:
		// discard name information
		// TODO: if this is an alias for a basic, can we make it an enum?
		return schemaFromType(ty.Underlying())
	case *types.Struct:
		sty := "object"
		props := map[string]config.Schema{}
		for i := 0; i < ty.NumFields(); i++ {
			f := ty.Field(i)

			fs, err := schemaFromType(f.Type())
			if err != nil {
				return nil, fmt.Errorf("unable to determine schema of field %q: %w", f.Name(), err)
			}
			if fs == nil {
				return nil, fmt.Errorf("unable to determine schema of field %q", f.Name())
			}

			props[f.Name()] = *fs
		}
		return &config.Schema{
			Type: &sty,
			SchemaProperties: &config.SchemaProperties{
				Properties: props,
			},
		}, nil
	}

	return nil, nil
}
