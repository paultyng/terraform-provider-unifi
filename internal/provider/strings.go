package provider

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func setToStringSlice(src *schema.Set) ([]string, error) {
	srcList := src.List()
	dst := make([]string, 0, len(srcList))
	for _, s := range srcList {
		d, ok := s.(string)
		if !ok {
			return nil, fmt.Errorf("unale to convert %v (%T) to string", s, s)
		}
		dst = append(dst, d)
	}
	return dst, nil
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
