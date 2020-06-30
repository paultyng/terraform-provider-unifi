package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/paultyng/go-unifi/unifi"
)

func resourceUserGroup() *schema.Resource {
	return &schema.Resource{
		Description: "`unifi_user_group` manages a user group (called \"client group\" in the UI), which can be used " +
			"to limit bandwidth for groups of users.",

		Create: resourceUserGroupCreate,
		Read:   resourceUserGroupRead,
		Update: resourceUserGroupUpdate,
		Delete: resourceUserGroupDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Description: "The name of the user group.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"qos_rate_max_down": {
				Description: "The QOS maximum download rate.",
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     -1,
				// TODO: validate does not equal 0,1
			},
			"qos_rate_max_up": {
				Description: "The QOS maximum upload rate.",
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     -1,
				// TODO: validate does not equal 0,1
			},
		},
	}
}

func resourceUserGroupCreate(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*client)

	req, err := resourceUserGroupGetResourceData(d)
	if err != nil {
		return err
	}

	resp, err := c.c.CreateUserGroup(context.TODO(), c.site, req)
	if err != nil {
		return err
	}

	d.SetId(resp.ID)

	return resourceUserGroupSetResourceData(resp, d)
}

func resourceUserGroupGetResourceData(d *schema.ResourceData) (*unifi.UserGroup, error) {
	return &unifi.UserGroup{
		Name:           d.Get("name").(string),
		QOSRateMaxDown: d.Get("qos_rate_max_down").(int),
		QOSRateMaxUp:   d.Get("qos_rate_max_up").(int),
	}, nil
}

func resourceUserGroupSetResourceData(resp *unifi.UserGroup, d *schema.ResourceData) error {
	d.Set("name", resp.Name)
	d.Set("qos_rate_max_down", resp.QOSRateMaxDown)
	d.Set("qos_rate_max_up", resp.QOSRateMaxUp)

	return nil
}

func resourceUserGroupRead(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*client)

	id := d.Id()

	resp, err := c.c.GetUserGroup(context.TODO(), c.site, id)
	if _, ok := err.(*unifi.NotFoundError); ok {
		d.SetId("")
		return nil
	}
	if err != nil {
		return err
	}

	return resourceUserGroupSetResourceData(resp, d)
}

func resourceUserGroupUpdate(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*client)

	req, err := resourceUserGroupGetResourceData(d)
	if err != nil {
		return err
	}

	req.ID = d.Id()
	req.SiteID = c.site

	resp, err := c.c.UpdateUserGroup(context.TODO(), c.site, req)
	if err != nil {
		return err
	}

	return resourceUserGroupSetResourceData(resp, d)
}

func resourceUserGroupDelete(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*client)

	id := d.Id()

	err := c.c.DeleteUserGroup(context.TODO(), c.site, id)
	if _, ok := err.(*unifi.NotFoundError); ok {
		return nil
	}
	return err
}
