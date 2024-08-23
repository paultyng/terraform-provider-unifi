package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/sayedh/go-unifi/unifi"
)

func resourceUserGroup() *schema.Resource {
	return &schema.Resource{
		Description: "`unifi_user_group` manages a user group (called \"client group\" in the UI), which can be used " +
			"to limit bandwidth for groups of users.",

		CreateContext: resourceUserGroupCreate,
		ReadContext:   resourceUserGroupRead,
		UpdateContext: resourceUserGroupUpdate,
		DeleteContext: resourceUserGroupDelete,
		Importer: &schema.ResourceImporter{
			StateContext: importSiteAndID,
		},

		Schema: map[string]*schema.Schema{
			"id": {
				Description: "The ID of the user group.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"site": {
				Description: "The name of the site to associate the user group with.",
				Type:        schema.TypeString,
				Computed:    true,
				Optional:    true,
				ForceNew:    true,
			},
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

func resourceUserGroupCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*client)

	req, err := resourceUserGroupGetResourceData(d)
	if err != nil {
		return diag.FromErr(err)
	}

	site := d.Get("site").(string)
	if site == "" {
		site = c.site
	}

	resp, err := c.c.CreateUserGroup(context.TODO(), site, req)
	if err != nil {
		return diag.FromErr(err)
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

func resourceUserGroupSetResourceData(resp *unifi.UserGroup, d *schema.ResourceData) diag.Diagnostics {
	d.Set("name", resp.Name)
	d.Set("qos_rate_max_down", resp.QOSRateMaxDown)
	d.Set("qos_rate_max_up", resp.QOSRateMaxUp)

	return nil
}

func resourceUserGroupRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*client)

	id := d.Id()

	site := d.Get("site").(string)
	if site == "" {
		site = c.site
	}

	resp, err := c.c.GetUserGroup(context.TODO(), site, id)
	if _, ok := err.(*unifi.NotFoundError); ok {
		d.SetId("")
		return nil
	}
	if err != nil {
		return diag.FromErr(err)
	}

	return resourceUserGroupSetResourceData(resp, d)
}

func resourceUserGroupUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*client)

	req, err := resourceUserGroupGetResourceData(d)
	if err != nil {
		return diag.FromErr(err)
	}

	req.ID = d.Id()

	site := d.Get("site").(string)
	if site == "" {
		site = c.site
	}
	req.SiteID = site

	resp, err := c.c.UpdateUserGroup(context.TODO(), site, req)
	if err != nil {
		return diag.FromErr(err)
	}

	return resourceUserGroupSetResourceData(resp, d)
}

func resourceUserGroupDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*client)

	id := d.Id()

	site := d.Get("site").(string)
	if site == "" {
		site = c.site
	}
	err := c.c.DeleteUserGroup(context.TODO(), site, id)
	if _, ok := err.(*unifi.NotFoundError); ok {
		return nil
	}
	return diag.FromErr(err)
}
