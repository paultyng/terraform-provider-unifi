package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataRADIUSProfile() *schema.Resource {
	return &schema.Resource{
		Description: "`unifi_radius_profile` data source can be used to retrieve the ID for a RADIUS profile by name.",

		Read: dataRADIUSProfileRead,

		Schema: map[string]*schema.Schema{
			"name": {
				Description: "The name of the RADIUS profile to look up.",
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "Default",
			},
		},
	}
}

func dataRADIUSProfileRead(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*client)

	name := d.Get("name").(string)

	profiles, err := c.c.ListRADIUSProfile(context.TODO(), c.site)
	if err != nil {
		return err
	}
	for _, g := range profiles {
		if g.Name == name {
			d.SetId(g.ID)
			return nil
		}
	}

	return fmt.Errorf("RADIUS profile not found with name %s", name)
}
