package provider

import (
	"context"
	"fmt"
	"os"
	"sync"
	"testing"

	"github.com/hashicorp/go-version"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
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

func siteAndIDImportStateIDFunc(resourceName string) func(*terraform.State) (string, error) {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("not found: %s", resourceName)
		}
		networkID := rs.Primary.Attributes["id"]
		site := rs.Primary.Attributes["site"]
		return site + ":" + networkID, nil
	}
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

func preCheckMinVersion(t *testing.T, min *version.Version) {
	v, err := version.NewVersion(testClient.Version())
	if err != nil {
		t.Fatalf("error parsing version: %s", err)
	}
	if v.LessThan(min) {
		t.Skipf("skipping test on controller version %q (need at least %q)", v, min)
	}
}

func preCheckV6Only(t *testing.T) {
	preCheckMinVersion(t, controllerV6)
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
