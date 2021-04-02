package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/paultyng/go-unifi/unifi"
)

func resourceSettingMgmt() *schema.Resource {
	return &schema.Resource{
		Description: "`unifi_setting_mgmt` manages settings for a unifi site.",

		Create: resourceSettingMgmtCreate,
		Read:   resourceSettingMgmtRead,
		Update: resourceSettingMgmtUpdate,
		Delete: resourceSettingMgmtDelete,
		Importer: &schema.ResourceImporter{
			State: importSiteAndID,
		},

		Schema: map[string]*schema.Schema{
			"id": {
				Description: "The ID of the settings.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"site": {
				Description: "The name of the site to associate the settings with.",
				Type:        schema.TypeString,
				Computed:    true,
				Optional:    true,
				ForceNew:    true,
			},
			"auto_upgrade": {
				Description: "Automatically upgrade device firmware.",
				Type:        schema.TypeBool,
				Optional:    true,
			},
		},
	}
}

func resourceSettingMgmtGetResourceData(d *schema.ResourceData, meta interface{}) (*unifi.SettingMgmt, error) {
	return &unifi.SettingMgmt{
		AutoUpgrade: d.Get("auto_upgrade").(bool),
	}, nil
}

func resourceSettingMgmtCreate(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*client)

	req, err := resourceSettingMgmtGetResourceData(d, meta)
	if err != nil {
		return err
	}

	site := d.Get("site").(string)
	if site == "" {
		site = c.site
	}

	resp, err := c.c.UpdateSettingMgmt(context.TODO(), site, req)
	if err != nil {
		return err
	}

	d.SetId(resp.ID)

	return resourceSettingMgmtSetResourceData(resp, d, meta, site)
}

func resourceSettingMgmtSetResourceData(resp *unifi.SettingMgmt, d *schema.ResourceData, meta interface{}, site string) error {
	d.Set("site", site)
	d.Set("auto_upgrade", resp.AutoUpgrade)
	return nil
}

func resourceSettingMgmtRead(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*client)

	site := d.Get("site").(string)
	if site == "" {
		site = c.site
	}

	resp, err := c.c.GetSettingMgmt(context.TODO(), site)
	if _, ok := err.(*unifi.NotFoundError); ok {
		d.SetId("")
		return nil
	}
	if err != nil {
		return err
	}

	return resourceSettingMgmtSetResourceData(resp, d, meta, site)
}

func resourceSettingMgmtUpdate(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*client)

	req, err := resourceSettingMgmtGetResourceData(d, meta)
	if err != nil {
		return err
	}

	req.ID = d.Id()
	site := d.Get("site").(string)
	if site == "" {
		site = c.site
	}

	resp, err := c.c.UpdateSettingMgmt(context.TODO(), site, req)
	if err != nil {
		return err
	}

	return resourceSettingMgmtSetResourceData(resp, d, meta, site)
}

func resourceSettingMgmtDelete(d *schema.ResourceData, meta interface{}) error {
	return nil
}
