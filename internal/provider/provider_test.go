package provider

import (
	"bytes"
	"context"
	"fmt"
	"math"
	"net"
	"os"
	"sync"
	"testing"

	"github.com/apparentlymart/go-cidr/cidr"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/paultyng/go-unifi/unifi"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/compose"
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

	os.Exit(runAcceptanceTests(m))
}

func runAcceptanceTests(m *testing.M) int {
	dc, err := compose.NewDockerCompose("../../docker-compose.yaml")
	if err != nil {
		panic(err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	if err = dc.WithOsEnv().Up(ctx, compose.Wait(true)); err != nil {
		panic(err)
	}

	defer func() {
		if err := dc.Down(context.Background(), compose.RemoveOrphans(true), compose.RemoveImagesLocal); err != nil {
			panic(err)
		}
	}()

	container, err := dc.ServiceContainer(ctx, "unifi")
	if err != nil {
		panic(err)
	}

	// Dump the container logs on exit.
	//
	// TODO: Use https://pkg.go.dev/github.com/testcontainers/testcontainers-go#LogConsumer instead.
	defer func() {
		if os.Getenv("UNIFI_STDOUT") == "" {
			return
		}

		stream, err := container.Logs(ctx)
		if err != nil {
			panic(err)
		}

		buffer := new(bytes.Buffer)
		buffer.ReadFrom(stream)
		testcontainers.Logger.Printf("%s", buffer)
	}()

	endpoint, err := container.PortEndpoint(ctx, "8443/tcp", "https")
	if err != nil {
		panic(err)
	}

	const user = "admin"
	const password = "admin"

	if err = os.Setenv("UNIFI_USERNAME", user); err != nil {
		panic(err)
	}

	if err = os.Setenv("UNIFI_PASSWORD", password); err != nil {
		panic(err)
	}

	if err = os.Setenv("UNIFI_INSECURE", "true"); err != nil {
		panic(err)
	}

	if err = os.Setenv("UNIFI_API", endpoint); err != nil {
		panic(err)
	}

	testClient = &unifi.Client{}
	setHTTPClient(testClient, true, "unifi")
	testClient.SetBaseURL(endpoint)
	if err = testClient.Login(ctx, user, password); err != nil {
		panic(err)
	}

	return m.Run()
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

const (
	vlanMin = 2
	vlanMax = 4095
)

var (
	network = &net.IPNet{
		IP:   net.IPv4(10, 0, 0, 0).To4(),
		Mask: net.IPv4Mask(255, 0, 0, 0),
	}

	vlanLock sync.Mutex
	vlanNext = vlanMin
)

func getTestVLAN(t *testing.T) (*net.IPNet, int) {
	vlanLock.Lock()
	defer vlanLock.Unlock()

	vlan := vlanNext
	vlanNext++

	subnet, err := cidr.Subnet(network, int(math.Ceil(math.Log2(vlanMax))), vlan)
	if err != nil {
		t.Error(err)
	}

	return subnet, vlan
}
