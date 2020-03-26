package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/paultyng/go-unifi/unifi"
)

func resourceFirewallGroup() *schema.Resource {
	return &schema.Resource{
		Create: resourceFirewallGroupCreate,
		Read:   resourceFirewallGroupRead,
		Update: resourceFirewallGroupUpdate,
		Delete: resourceFirewallGroupDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"type": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringInSlice([]string{"address-group", "port-group", "ipv6-address-group"}, false),
			},
			"members": {
				Type:     schema.TypeSet,
				Required: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
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

	resp, err := c.c.CreateFirewallGroup(context.TODO(), c.site, req)
	if err != nil {
		return err
	}

	d.SetId(resp.ID)

	return resourceFirewallGroupSetResourceData(resp, d)
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

func resourceFirewallGroupSetResourceData(resp *unifi.FirewallGroup, d *schema.ResourceData) error {
	d.Set("name", resp.Name)
	d.Set("type", resp.GroupType)
	d.Set("members", stringSliceToSet(resp.GroupMembers))

	return nil
}

func resourceFirewallGroupRead(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*client)

	id := d.Id()

	resp, err := c.c.GetFirewallGroup(context.TODO(), c.site, id)
	if _, ok := err.(*unifi.NotFoundError); ok {
		d.SetId("")
		return nil
	}
	if err != nil {
		return err
	}

	return resourceFirewallGroupSetResourceData(resp, d)
}

func resourceFirewallGroupUpdate(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*client)

	req, err := resourceFirewallGroupGetResourceData(d)
	if err != nil {
		return err
	}

	req.ID = d.Id()
	req.SiteID = c.site

	resp, err := c.c.UpdateFirewallGroup(context.TODO(), c.site, req)
	if err != nil {
		return err
	}

	return resourceFirewallGroupSetResourceData(resp, d)
}

func resourceFirewallGroupDelete(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*client)

	id := d.Id()

	err := c.c.DeleteFirewallGroup(context.TODO(), c.site, id)
	if _, ok := err.(*unifi.NotFoundError); ok {
		return nil
	}
	return err
}
