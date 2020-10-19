package provider

import (
	"context"
	"os"
	"sync"
	"testing"

	"github.com/hashicorp/go-version"
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

func preCheckV6Only(t *testing.T) {
	v, err := version.NewVersion(testClient.Version())
	if err != nil {
		t.Fatalf("error parsing version: %s", err)
	}
	if v.LessThan(controllerV6) {
		t.Skipf("skipping test on controller version %q", v)
	}
}

func preCheckV5Only(t *testing.T) {
	v, err := version.NewVersion(testClient.Version())
	if err != nil {
		t.Fatalf("error parsing version: %s", err)
	}
	if v.GreaterThanOrEqual(controllerV6) {
		t.Skipf("skipping test on controller version %q", v)
	}
}

const (
	vlanMin = 2
	vlanMax = 4095
)

var (
	vlanLock sync.Mutex
	vlanNext = vlanMin
)

func getTestVLAN(t *testing.T) int {
	vlanLock.Lock()
	defer vlanLock.Unlock()

	vl := vlanNext
	vlanNext++

	return vl
}
