package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/ubiquiti-community/go-unifi/unifi"
)

// TODO: probably need to update this to be more like setting_usg,
// using locking, and upsert, more computed, etc.

func resourceSettingMgmt() *schema.Resource {
	return &schema.Resource{
		Description: "`unifi_setting_mgmt` manages settings for a unifi site.",

		CreateContext: resourceSettingMgmtCreate,
		ReadContext:   resourceSettingMgmtRead,
		UpdateContext: resourceSettingMgmtUpdate,
		DeleteContext: resourceSettingMgmtDelete,
		Importer: &schema.ResourceImporter{
			StateContext: importSiteAndID,
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
			"ssh_enabled": {
				Description: "Enable SSH authentication.",
				Type:        schema.TypeBool,
				Optional:    true,
			},
			"ssh_key": {
				Description: "SSH key.",
				Type:        schema.TypeSet,
				Optional:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Description: "Name of SSH key.",
							Type:        schema.TypeString,
							Required:    true,
						},
						"type": {
							Description: "Type of SSH key, e.g. ssh-rsa.",
							Type:        schema.TypeString,
							Required:    true,
						},
						"key": {
							Description: "Public SSH key.",
							Type:        schema.TypeString,
							Optional:    true,
						},
						"comment": {
							Description: "Comment.",
							Type:        schema.TypeString,
							Optional:    true,
						},
					},
				},
			},
		},
	}
}

func setToSshKeys(set *schema.Set) ([]unifi.SettingMgmtXSshKeys, error) {
	var sshKeys []unifi.SettingMgmtXSshKeys
	for _, item := range set.List() {
		data, ok := item.(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("unexpected data in block")
		}
		sshKey, err := toSshKey(data)
		if err != nil {
			return nil, fmt.Errorf("unable to create port override: %w", err)
		}
		sshKeys = append(sshKeys, sshKey)
	}
	return sshKeys, nil
}

func toSshKey(data map[string]interface{}) (unifi.SettingMgmtXSshKeys, error) {
	return unifi.SettingMgmtXSshKeys{
		Name:    data["name"].(string),
		KeyType: data["type"].(string),
		Key:     data["key"].(string),
		Comment: data["comment"].(string),
	}, nil
}

func setFromSshKeys(sshKeys []unifi.SettingMgmtXSshKeys) ([]map[string]interface{}, error) {
	list := make([]map[string]interface{}, 0, len(sshKeys))
	for _, sshKey := range sshKeys {
		v, err := fromSshKey(sshKey)
		if err != nil {
			return nil, fmt.Errorf("unable to parse ssh key: %w", err)
		}
		list = append(list, v)
	}
	return list, nil
}

func fromSshKey(sshKey unifi.SettingMgmtXSshKeys) (map[string]interface{}, error) {
	return map[string]interface{}{
		"name":    sshKey.Name,
		"type":    sshKey.KeyType,
		"key":     sshKey.Key,
		"comment": sshKey.Comment,
	}, nil
}

func resourceSettingMgmtGetResourceData(d *schema.ResourceData, meta interface{}) (*unifi.SettingMgmt, error) {
	sshKeys, err := setToSshKeys(d.Get("ssh_key").(*schema.Set))
	if err != nil {
		return nil, fmt.Errorf("unable to process ssh_key block: %w", err)
	}

	return &unifi.SettingMgmt{
		AutoUpgrade: d.Get("auto_upgrade").(bool),
		XSshEnabled: d.Get("ssh_enabled").(bool),
		XSshKeys:    sshKeys,
	}, nil
}

func resourceSettingMgmtCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*client)

	req, err := resourceSettingMgmtGetResourceData(d, meta)
	if err != nil {
		return diag.FromErr(err)
	}

	site := d.Get("site").(string)
	if site == "" {
		site = c.site
	}

	resp, err := c.c.UpdateSettingMgmt(ctx, site, req)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(resp.ID)

	return resourceSettingMgmtSetResourceData(resp, d, meta, site)
}

func resourceSettingMgmtSetResourceData(resp *unifi.SettingMgmt, d *schema.ResourceData, meta interface{}, site string) diag.Diagnostics {
	sshKeys, err := setFromSshKeys(resp.XSshKeys)
	if err != nil {
		return diag.FromErr(err)
	}

	d.Set("site", site)
	d.Set("auto_upgrade", resp.AutoUpgrade)
	d.Set("ssh_enabled", resp.XSshEnabled)
	d.Set("ssh_key", sshKeys)
	return nil
}

func resourceSettingMgmtRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*client)

	site := d.Get("site").(string)
	if site == "" {
		site = c.site
	}

	resp, err := c.c.GetSettingMgmt(ctx, site)
	if _, ok := err.(*unifi.NotFoundError); ok {
		d.SetId("")
		return nil
	}
	if err != nil {
		return diag.FromErr(err)
	}

	return resourceSettingMgmtSetResourceData(resp, d, meta, site)
}

func resourceSettingMgmtUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*client)

	req, err := resourceSettingMgmtGetResourceData(d, meta)
	if err != nil {
		return diag.FromErr(err)
	}

	req.ID = d.Id()
	site := d.Get("site").(string)
	if site == "" {
		site = c.site
	}

	resp, err := c.c.UpdateSettingMgmt(ctx, site, req)
	if err != nil {
		return diag.FromErr(err)
	}

	return resourceSettingMgmtSetResourceData(resp, d, meta, site)
}

func resourceSettingMgmtDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return nil
}
