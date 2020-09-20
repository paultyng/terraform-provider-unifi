package provider

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func hasConfig(t *testing.T, path string) bool {
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
	return config != ""
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

		if !hasConfig(t, path) {
			return nil
		}

		name, err := filepath.Rel(examplesPath, path)
		if err != nil {
			t.Fatal(err)
		}

		t.Run(name, func(t *testing.T) {
			switch name {
			case "csv_users":
				t.Skipf("for_each is not yet supported by acc test framework")
			}

			resource.ParallelTest(t, resource.TestCase{
				ConfigDir:         path,
				PreCheck:          wlanPreCheck(t),
				ProviderFactories: providerFactories,
				CheckDestroy: func(*terraform.State) error {
					// TODO: actual CheckDestroy

					<-wlanConcurrency
					return nil
				},
				Steps: []resource.TestStep{
					{
						Check: resource.ComposeTestCheckFunc(
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
