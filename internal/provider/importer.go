package provider

import (
	"context"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func importSiteAndID(ctx context.Context, d *schema.ResourceData, meta any) ([]*schema.ResourceData, error) {
	if id := d.Id(); strings.Contains(id, ":") {
		importParts := strings.SplitN(id, ":", 2)
		d.SetId(importParts[1])
		d.Set("site", importParts[0])
	}
	return []*schema.ResourceData{d}, nil
}
