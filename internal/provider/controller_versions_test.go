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
