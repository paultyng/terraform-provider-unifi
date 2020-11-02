package provider

import (
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func ImportHandleSite(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	if strings.Contains(d.Id(), ":") {
		importParts := strings.Split(d.Id(), ":")
		d.SetId(importParts[1])
		d.Set("site", importParts[0])
	}
	return []*schema.ResourceData{d}, nil
}
