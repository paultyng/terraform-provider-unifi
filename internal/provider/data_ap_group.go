package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataAPGroup() *schema.Resource {
	return &schema.Resource{
		Description: "`unifi_ap_group` data source can be used to retrieve the ID for an AP group by name.",

		Read: dataAPGroupRead,

		Schema: map[string]*schema.Schema{
			"id": {
				Description: "The ID of this AP group.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"name": {
				Description: "The name of the AP group to look up, leave blank to look up the default AP group.",
				Type:        schema.TypeString,
				Optional:    true,
			},
		},
	}
}

func dataAPGroupRead(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*client)

	if v := c.ControllerVersion(); !v.GreaterThanOrEqual(controllerV6) {
		return fmt.Errorf("AP groups are not supported on controller version %q", v)
	}

	name := d.Get("name").(string)

	groups, err := c.c.ListAPGroup(context.TODO(), c.site)
	if err != nil {
		return err
	}
	for _, g := range groups {
		if (name == "" && g.HiddenID == "default") || g.Name == name {
			d.SetId(g.ID)
			return nil
		}
	}

	return fmt.Errorf("AP group not found with name %s", name)
}
