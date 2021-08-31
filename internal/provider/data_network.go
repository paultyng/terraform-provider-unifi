package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataNetwork() *schema.Resource {
	return &schema.Resource{
		Description: "`unifi_network` data source can be used to retrieve the ID for a network by name.",

		ReadContext: dataNetworkRead,

		Schema: map[string]*schema.Schema{
			"id": {
				Description: "The ID of this network.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"site": {
				Description: "The name of the site the network is associated with.",
				Type:        schema.TypeString,
				Computed:    true,
				Optional:    true,
			},
			"name": {
				Description: "The name of the network to look up.",
				Type:        schema.TypeString,
				Required:    true,
			},
		},
	}
}

func dataNetworkRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*client)

	name := d.Get("name").(string)
	site := d.Get("site").(string)
	if site == "" {
		site = c.site
	}

	networks, err := c.c.ListNetwork(ctx, site)
	if err != nil {
		return diag.FromErr(err)
	}
	for _, n := range networks {
		if n.Name == name {
			d.SetId(n.ID)

			d.Set("site", site)

			return nil
		}
	}

	return diag.Errorf("network not found with name %s", name)
}
