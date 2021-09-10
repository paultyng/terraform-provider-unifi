package provider

import (
	"fmt"

	"github.com/hashicorp/go-version"
)

var (
	controllerV6 = version.Must(version.NewVersion("6.0.0"))

	// https://community.ui.com/releases/UniFi-Network-Controller-6-1-61/62f1ad38-1ac5-430c-94b0-becbb8f71d7d
	controllerVersionWPA3 = version.Must(version.NewVersion("6.1.61"))
)

func (c *client) ControllerVersion() *version.Version {
	return version.Must(version.NewVersion(c.c.Version()))
}

func checkMinimumControllerVersion(versionString string) error {
	v, err := version.NewVersion(versionString)
	if err != nil {
		return err
	}
	if v.LessThan(controllerV6) {
		return fmt.Errorf("Controller version %q or greater is required to use the provider, found %q.", controllerV6, v)
	}
	return nil
}
