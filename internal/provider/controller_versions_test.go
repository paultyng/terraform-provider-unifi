package provider

import (
	"testing"

	"github.com/hashicorp/go-version"
)

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
