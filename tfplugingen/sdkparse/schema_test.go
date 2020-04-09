package sdkparse

import (
	"context"
	"fmt"
	"go/types"
	"reflect"
	"testing"

	"github.com/paultyng/tfplugingen/config"
)

func TestSchemaFromType(t *testing.T) {
	ctx := context.Background()

	pkgMap, err := loadPackages(ctx, "../testdata/testprovider", t.Logf, "example.com/testprovider/testsdk")
	if err != nil {
		t.Fatal(err)
	}

	sdkPkg := pkgMap["example.com/testprovider/testsdk"]

	sdkType := func(name string) types.Type {
		obj := sdkPkg.Types.Scope().Lookup(name)
		tn := obj.(*types.TypeName)
		return tn.Type()
	}

	ps := func(s string) *string {
		return &s
	}

	for _, c := range []struct {
		expected *config.Schema
		ty       types.Type
	}{
		{&config.Schema{Type: ps("boolean")}, types.Typ[types.Bool]},

		{&config.Schema{Type: ps("integer")}, types.Typ[types.Int]},
		{&config.Schema{Type: ps("integer")}, types.Typ[types.Int8]},
		{&config.Schema{Type: ps("integer")}, types.Typ[types.Int16]},
		{&config.Schema{Type: ps("integer")}, types.Typ[types.Int32]},
		{&config.Schema{Type: ps("integer")}, types.Typ[types.Int64]},
		{&config.Schema{Type: ps("integer")}, sdkType("AliasedInt64")},

		{&config.Schema{Type: ps("number")}, types.Typ[types.Float32]},
		{&config.Schema{Type: ps("number")}, types.Typ[types.Float64]},

		{&config.Schema{Type: ps("string")}, types.Typ[types.String]},

		{&config.Schema{
			Type: ps("object"),
			SchemaProperties: &config.SchemaProperties{
				Properties: map[string]config.Schema{
					"String":  config.Schema{Type: ps("string")},
					"Int":     config.Schema{Type: ps("integer")},
					"Int64":   config.Schema{Type: ps("integer")},
					"Bool":    config.Schema{Type: ps("boolean")},
					"Uint32":  config.Schema{Type: ps("integer")},
					"Float32": config.Schema{Type: ps("number")},
				},
			},
		}, sdkType("Simple")},
	} {
		t.Run(fmt.Sprintf("%#v", c.ty), func(t *testing.T) {
			actual, err := schemaFromType(c.ty)
			if err != nil {
				t.Fatal(err)
			}
			if !reflect.DeepEqual(c.expected, actual) {
				t.Fatalf("expected %#v, got %#v", c.expected, actual)
			}
		})
	}
}
