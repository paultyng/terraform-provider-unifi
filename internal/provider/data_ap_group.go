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
			"site": {
				Description: "The name of the site the AP group is associated with.",
				Type:        schema.TypeString,
				Computed:    true,
				Optional:    true,
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
	site := d.Get("site").(string)
	if site == "" {
		site = c.site
	}

	groups, err := c.c.ListAPGroup(context.TODO(), site)
	if err != nil {
		return err
	}
	for _, g := range groups {
		if (name == "" && g.HiddenID == "default") || g.Name == name {
			d.SetId(g.ID)
			d.Set("site", site)
			return nil
		}
	}

	return fmt.Errorf("AP group not found with name %s", name)
}
