package provider

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataUserGroup() *schema.Resource {
	return &schema.Resource{
		Read: dataUserGroupRead,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "Default",
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

	groups, err := c.c.ListUserGroup(c.site)
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
