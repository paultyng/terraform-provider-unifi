package provider

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataWLANGroup() *schema.Resource {
	return &schema.Resource{
		Read: dataWLANGroupRead,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "Default",
			},
		},
	}
}

func dataWLANGroupRead(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*client)

	name := d.Get("name").(string)

	groups, err := c.c.ListWLANGroup(c.site)
	if err != nil {
		return err
	}
	for _, g := range groups {
		if g.Name == name {
			d.SetId(g.ID)
			return nil
		}
	}

	return fmt.Errorf("WLAN group not found with name %s", name)
}
