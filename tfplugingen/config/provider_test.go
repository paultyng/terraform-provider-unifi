package config

import (
	"testing"
)

func TestDefaultResourceType(t *testing.T) {
	for _, c := range []struct {
		name     string
		expected string
	}{
		{"test_simple", "Simple"},
		{"test_longer_name", "LongerName"},
	} {
		t.Run(c.name, func(t *testing.T) {
			p := Provider{}
			actual := p.DefaultResourceType(c.name)
			if c.expected != actual {
				t.Fatalf("expected %q, got %q", c.expected, actual)
			}

		})
	}
}
