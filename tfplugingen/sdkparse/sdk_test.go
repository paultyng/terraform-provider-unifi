package sdkparse

import (
	"context"
	"testing"

	"github.com/paultyng/tfplugingen/config"
)

func TestCollectResourceInfo(t *testing.T) {
	ctx := context.Background()

	pkgMap, err := loadPackages(ctx, "../testdata/testprovider", t.Logf, "example.com/testprovider/testsdk")
	if err != nil {
		t.Fatal(err)
	}

	resConf := config.Resource{
		Name: "test_simple",
	}

	provConf := config.Provider{
		DefaultSDKPackage: "example.com/testprovider/testsdk",
	}

	ri, err := collectResourceInfo(ctx, resConf, provConf, pkgMap)
	if err != nil {
		t.Fatal(err)
	}

	if expected, actual := "Simple", ri.typeName.Name(); expected != actual {
		t.Fatalf("expected %q, got %q", expected, actual)
	}

}
