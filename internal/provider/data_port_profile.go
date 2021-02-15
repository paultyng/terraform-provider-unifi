package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataPortProfile() *schema.Resource {
	return &schema.Resource{
		Description: "`unifi_port_profile` data source can be used to retrieve the ID for a port profile by name.",

		Read: dataPortProfileRead,

		Schema: map[string]*schema.Schema{
			"id": {
				Description: "The ID of this port profile.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"site": {
				Description: "The name of the site the port profile is associated with.",
				Type:        schema.TypeString,
				Computed:    true,
				Optional:    true,
			},
			"name": {
				Description: "The name of the port profile to look up.",
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "All",
			},
		},
	}
}

func dataPortProfileRead(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*client)

	name := d.Get("name").(string)
	site := d.Get("site").(string)
	if site == "" {
		site = c.site
	}

	groups, err := c.c.ListPortProfile(context.TODO(), site)
	if err != nil {
		return err
	}
	for _, g := range groups {
		if g.Name == name {
			d.SetId(g.ID)

			d.Set("site", site)

			return nil
		}
	}

	return fmt.Errorf("port profile not found with name %s", name)
}
