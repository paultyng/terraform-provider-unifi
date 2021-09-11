package provider

import (
	"github.com/hashicorp/go-version"
)

var (
	controllerV5 = version.Must(version.NewVersion("5.0.0"))
	controllerV6 = version.Must(version.NewVersion("6.0.0"))

	// https://community.ui.com/releases/UniFi-Network-Controller-6-1-61/62f1ad38-1ac5-430c-94b0-becbb8f71d7d
	controllerVersionWPA3 = version.Must(version.NewVersion("6.1.61"))
)

func (c *client) ControllerVersion() *version.Version {
	return version.Must(version.NewVersion(c.c.Version()))
}
