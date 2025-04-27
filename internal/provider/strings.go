package provider

import (
	"fmt"
	"sort"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func listToStringSlice(src []interface{}) ([]string, error) {
	dst := make([]string, 0, len(src))
	for _, s := range src {
		d, ok := s.(string)
		if !ok {
			return nil, fmt.Errorf("unale to convert %v (%T) to string", s, s)
		}
		dst = append(dst, d)
	}
	return dst, nil
}

func setToStringSlice(src *schema.Set) ([]string, error) {
	return listToStringSlice(src.List())
}

func stringSliceToList(list []string) []interface{} {
	vs := make([]interface{}, 0, len(list))
	for _, v := range list {
		vs = append(vs, v)
	}
	return vs
}

func stringSliceToSet(src []string) *schema.Set {
	return schema.NewSet(schema.HashString, stringSliceToList(src))
}

func compareStringSlicesWithoutOrder(s, t []string) bool {
	if len(s) != len(t) {
		return false
	}

	ss := make([]string, len(s))
	tt := make([]string, len(t))

	copy(ss, s)
	copy(tt, t)

	sort.Strings(ss)
	sort.Strings(tt)

	for i, v := range ss {
		if v != tt[i] {
			return false
		}
	}
	return true
}
