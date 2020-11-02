package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/paultyng/go-unifi/unifi"
)

func resourceFirewallGroup() *schema.Resource {
	return &schema.Resource{
		Description: "`unifi_firewall_group` manages groups of addresses or ports for use in firewall rules (`unifi_firewall_rule`).",

		Create: resourceFirewallGroupCreate,
		Read:   resourceFirewallGroupRead,
		Update: resourceFirewallGroupUpdate,
		Delete: resourceFirewallGroupDelete,
		Importer: &schema.ResourceImporter{
			State: ImportHandleSite,
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

func resourceFirewallGroupCreate(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*client)

	req, err := resourceFirewallGroupGetResourceData(d)
	if err != nil {
		return err
	}

	site := d.Get("site").(string)
	if site == "" {
		site = c.site
	}

	resp, err := c.c.CreateFirewallGroup(context.TODO(), site, req)
	if err != nil {
		return err
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

func resourceFirewallGroupSetResourceData(resp *unifi.FirewallGroup, d *schema.ResourceData, site string) error {
	d.Set("site", site)
	d.Set("name", resp.Name)
	d.Set("type", resp.GroupType)
	d.Set("members", stringSliceToSet(resp.GroupMembers))

	return nil
}

func resourceFirewallGroupRead(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*client)

	id := d.Id()

	site := d.Get("site").(string)
	if site == "" {
		site = c.site
	}

	resp, err := c.c.GetFirewallGroup(context.TODO(), site, id)
	if _, ok := err.(*unifi.NotFoundError); ok {
		d.SetId("")
		return nil
	}
	if err != nil {
		return err
	}

	return resourceFirewallGroupSetResourceData(resp, d, site)
}

func resourceFirewallGroupUpdate(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*client)

	req, err := resourceFirewallGroupGetResourceData(d)
	if err != nil {
		return err
	}

	req.ID = d.Id()

	site := d.Get("site").(string)
	if site == "" {
		site = c.site
	}
	req.SiteID = site

	resp, err := c.c.UpdateFirewallGroup(context.TODO(), site, req)
	if err != nil {
		return err
	}

	return resourceFirewallGroupSetResourceData(resp, d, site)
}

func resourceFirewallGroupDelete(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*client)

	id := d.Id()

	site := d.Get("site").(string)
	if site == "" {
		site = c.site
	}

	err := c.c.DeleteFirewallGroup(context.TODO(), site, id)
	if _, ok := err.(*unifi.NotFoundError); ok {
		return nil
	}
	return err
}
