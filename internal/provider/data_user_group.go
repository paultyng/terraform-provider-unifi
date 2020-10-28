package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataUserGroup() *schema.Resource {
	return &schema.Resource{
		Description: "`unifi_user_group` data source can be used to retrieve the ID for a user group by name.",

		Read: dataUserGroupRead,

		Schema: map[string]*schema.Schema{
			"id": {
				Description: "The ID of this AP group.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"site": {
				Description: "The name of the site the user group is associated with.",
				Type:        schema.TypeString,
				Computed:    true,
				Optional:    true,
			},
			"name": {
				Description: "The name of the user group to look up.",
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "Default",
			},

			"qos_rate_max_down": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"qos_rate_max_up": {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
}

func dataUserGroupRead(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*client)

	name := d.Get("name").(string)
	site := d.Get("site").(string)
	if site == "" {
		site = c.site
	}

	groups, err := c.c.ListUserGroup(context.TODO(), site)
	if err != nil {
		return err
	}
	for _, g := range groups {
		if g.Name == name {
			d.SetId(g.ID)

			d.Set("site", site)
			d.Set("qos_rate_max_down", g.QOSRateMaxDown)
			d.Set("qos_rate_max_up", g.QOSRateMaxUp)

			return nil
		}
	}

	return fmt.Errorf("user group not found with name %s", name)
}
