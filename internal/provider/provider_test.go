package provider

import (
	"context"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/paultyng/go-unifi/unifi"
)

var providerFactories = map[string]func() (*schema.Provider, error){
	"unifi": func() (*schema.Provider, error) {
		return New("acctest")(), nil
	},
}

var testClient *unifi.Client

func TestMain(m *testing.M) {
	if os.Getenv("TF_ACC") == "" {
		// short circuit non acceptance test runs
		os.Exit(m.Run())
	}

	user := os.Getenv("UNIFI_USERNAME")
	pass := os.Getenv("UNIFI_PASSWORD")
	baseURL := os.Getenv("UNIFI_API")
	insecure := os.Getenv("UNIFI_INSECURE") == "true"

	testClient = &unifi.Client{}
	setHTTPClient(testClient, insecure)
	testClient.SetBaseURL(baseURL)
	err := testClient.Login(context.Background(), user, pass)
	if err != nil {
		panic(err)
	}

	resource.TestMain(m)
}

func importStep(name string, ignore ...string) resource.TestStep {
	step := resource.TestStep{
		ResourceName:      name,
		ImportState:       true,
		ImportStateVerify: true,
	}

	if len(ignore) > 0 {
		step.ImportStateVerifyIgnore = ignore
	}

	return step
}

func preCheck(t *testing.T) {
	variables := []string{
		"UNIFI_USERNAME",
		"UNIFI_PASSWORD",
		"UNIFI_API",
	}

	for _, variable := range variables {
		value := os.Getenv(variable)
		if value == "" {
			t.Fatalf("`%s` must be set for acceptance tests!", variable)
		}
	}
}
