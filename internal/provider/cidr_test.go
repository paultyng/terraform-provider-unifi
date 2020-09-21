package provider

import (
	"testing"
)

func TestCIDRValidate(t *testing.T) {
	for _, c := range []struct {
		expectedError string
		cidr          string
	}{
		{"invalid CIDR address: ", ""},
		{"invalid CIDR address: abc", "abc"},
		{"invalid CIDR address: 192.1.2.3", "192.1.2.3"},
		{"invalid CIDR address: 500.1.2.3/20", "500.1.2.3/20"},
		{"invalid CIDR address: 192.1.2.3/500", "192.1.2.3/500"},

		{"", "192.1.2.1/20"},
	} {
		t.Run(c.cidr, func(t *testing.T) {
			_, actualErrs := cidrValidate(c.cidr, "key")
			switch len(actualErrs) {
			case 0:
				if c.expectedError != "" {
					t.Fatalf("expected no error, got %d: %#v", len(actualErrs), actualErrs)
				}
			case 1:
				actualErr := actualErrs[0].Error()
				if actualErr != c.expectedError {
					t.Fatalf("expected %q, got %q", c.expectedError, actualErr)
				}
			default:
				t.Fatalf("expected 0 or 1 errors, got %d: %#v", len(actualErrs), actualErrs)
			}
		})
	}
}
