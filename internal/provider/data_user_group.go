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

	groups, err := c.c.ListUserGroup(context.TODO(), c.site)
	if err != nil {
		return err
	}
	for _, g := range groups {
		if g.Name == name {
			d.SetId(g.ID)

			d.Set("qos_rate_max_down", g.QOSRateMaxDown)
			d.Set("qos_rate_max_up", g.QOSRateMaxUp)

			return nil
		}
	}

	return fmt.Errorf("user group not found with name %s", name)
}
