package provider

import (
	"context"
	"errors"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/paultyng/go-unifi/unifi"
)

func resourceFirewallGroup() *schema.Resource {
	return &schema.Resource{
		Description: "`unifi_firewall_group` manages groups of addresses or ports for use in firewall rules (`unifi_firewall_rule`).",

		CreateContext: resourceFirewallGroupCreate,
		ReadContext:   resourceFirewallGroupRead,
		UpdateContext: resourceFirewallGroupUpdate,
		DeleteContext: resourceFirewallGroupDelete,
		Importer: &schema.ResourceImporter{
			StateContext: importSiteAndID,
		},

		Schema: map[string]*schema.Schema{
			"id": {
				Description: "The ID of the firewall group.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"site": {
				Description: "The name of the site to associate the firewall group with.",
				Type:        schema.TypeString,
				Computed:    true,
				Optional:    true,
				ForceNew:    true,
			},
			"name": {
				Description: "The name of the firewall group.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"type": {
				Description:  "The type of the firewall group. Must be one of: `address-group`, `port-group`, or `ipv6-address-group`.",
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringInSlice([]string{"address-group", "port-group", "ipv6-address-group"}, false),
			},
			"members": {
				Description: "The members of the firewall group.",
				Type:        schema.TypeSet,
				Required:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func resourceFirewallGroupCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*client)

	req, err := resourceFirewallGroupGetResourceData(d)
	if err != nil {
		return diag.FromErr(err)
	}

	site := d.Get("site").(string)
	if site == "" {
		site = c.site
	}

	resp, err := c.c.CreateFirewallGroup(ctx, site, req)
	if err != nil {
		var apiErr *unifi.APIError
		if errors.As(err, &apiErr) && apiErr.Message == "api.err.FirewallGroupExisted" {
			return diag.Errorf("firewall groups must have unique names: %s", err)
		}

		return diag.FromErr(err)
	}

	d.SetId(resp.ID)

	return resourceFirewallGroupSetResourceData(resp, d, site)
}

func resourceFirewallGroupGetResourceData(d *schema.ResourceData) (*unifi.FirewallGroup, error) {
	members, err := setToStringSlice(d.Get("members").(*schema.Set))
	if err != nil {
		return nil, err
	}

	return &unifi.FirewallGroup{
		Name:         d.Get("name").(string),
		GroupType:    d.Get("type").(string),
		GroupMembers: members,
	}, nil
}

func resourceFirewallGroupSetResourceData(resp *unifi.FirewallGroup, d *schema.ResourceData, site string) diag.Diagnostics {
	d.Set("site", site)
	d.Set("name", resp.Name)
	d.Set("type", resp.GroupType)
	d.Set("members", stringSliceToSet(resp.GroupMembers))

	return nil
}

func resourceFirewallGroupRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*client)

	id := d.Id()

	site := d.Get("site").(string)
	if site == "" {
		site = c.site
	}

	resp, err := c.c.GetFirewallGroup(ctx, site, id)
	if _, ok := err.(*unifi.NotFoundError); ok {
		d.SetId("")
		return nil
	}
	if err != nil {
		return diag.FromErr(err)
	}

	return resourceFirewallGroupSetResourceData(resp, d, site)
}

func resourceFirewallGroupUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*client)

	req, err := resourceFirewallGroupGetResourceData(d)
	if err != nil {
		return diag.FromErr(err)
	}

	req.ID = d.Id()

	site := d.Get("site").(string)
	if site == "" {
		site = c.site
	}
	req.SiteID = site

	resp, err := c.c.UpdateFirewallGroup(ctx, site, req)
	if err != nil {
		return diag.FromErr(err)
	}

	return resourceFirewallGroupSetResourceData(resp, d, site)
}

func resourceFirewallGroupDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*client)

	id := d.Id()

	site := d.Get("site").(string)
	if site == "" {
		site = c.site
	}

	err := c.c.DeleteFirewallGroup(ctx, site, id)
	if _, ok := err.(*unifi.NotFoundError); ok {
		return nil
	}
	return diag.FromErr(err)
}
