package provider

import (
	"os"
	"io/ioutil"
	"path/filepath"
	"testing"
	"strings"

	
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func directoryConfig(t *testing.T, path string) string {
	t.Helper()

	files, err := ioutil.ReadDir(path)
	if err != nil {
		t.Fatal(err)
	}

	config := ""
	for _, f := range files {
		ext := filepath.Ext(f.Name())
		if ext != ".tf" {
			continue
		}

		data, err := ioutil.ReadFile(filepath.Join(path, f.Name()))
		if err != nil {
			t.Fatal(err)
		}

		config += "\n\n" + string(data) + "\n\n"	
	}

	config = strings.TrimSpace(config)
	return config
}

const examplesPath = "../../examples"

func TestAccExamples(t *testing.T) {
	err := filepath.Walk(examplesPath, func(path string, fi os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		
		if !fi.IsDir() {
			return nil
		}

		config := directoryConfig(t, path)
		if config == "" {
			return nil
		}

		name, err := filepath.Rel(examplesPath, path)
		if err != nil {
			t.Fatal(err)
		}
		t.Run(name, func(t *testing.T) {
			resource.ParallelTest(t, resource.TestCase{
				PreCheck:          wlanPreCheck(t),
				ProviderFactories: providerFactories,
				CheckDestroy: func(*terraform.State) error {
					// TODO: actual CheckDestroy
		
					<-wlanConcurrency
					return nil
				},
				Steps: []resource.TestStep{
					{
						Config: config,
						Check:  resource.ComposeTestCheckFunc(
						// testCheckNetworkExists(t, "name"),
						),
					},
				},
			})
		})

		// use error here to stop recursing?
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
}
