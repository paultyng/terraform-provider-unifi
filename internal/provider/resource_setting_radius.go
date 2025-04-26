package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/ubiquiti-community/go-unifi/unifi"
)

func resourceSettingRadius() *schema.Resource {
	return &schema.Resource{
		Description: "`unifi_setting_radius` manages settings for the built-in RADIUS server.",

		CreateContext: resourceSettingRadiusCreate,
		ReadContext:   resourceSettingRadiusRead,
		UpdateContext: resourceSettingRadiusUpdate,
		DeleteContext: schema.NoopContext,
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
			"accounting_enabled": {
				Description: "Enable RADIUS accounting",
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
			},
			"accounting_port": {
				Description:  "The port for accounting communications.",
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      1813,
				ValidateFunc: validation.IsPortNumber,
			},
			"auth_port": {
				Description:  "The port for authentication communications.",
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      1812,
				ValidateFunc: validation.IsPortNumber,
			},
			"interim_update_interval": {
				Description: "Statistics will be collected from connected clients at this interval.",
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     3600,
			},
			"tunneled_reply": {
				Description: "Encrypt communication between the server and the client.",
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
			},
			"secret": {
				Description: "RAIDUS secret passphrase.",
				Type:        schema.TypeString,
				Sensitive:   true,
				Optional:    true,
				Default:     "",
			},
			"enabled": {
				Description: "RAIDUS server enabled.",
				Type:        schema.TypeBool,
				Default:     true,
				Optional:    true,
			},
		},
	}
}

func resourceSettingRadiusGetResourceData(d *schema.ResourceData, meta any) (*unifi.SettingRadius, error) {
	return &unifi.SettingRadius{
		AccountingEnabled:     d.Get("accounting_enabled").(bool),
		Enabled:               d.Get("enabled").(bool),
		AcctPort:              d.Get("accounting_port").(int),
		AuthPort:              d.Get("auth_port").(int),
		ConfigureWholeNetwork: true,
		TunneledReply:         d.Get("tunneled_reply").(bool),
		XSecret:               d.Get("secret").(string),
		InterimUpdateInterval: d.Get("interim_update_interval").(int),
	}, nil
}

func resourceSettingRadiusCreate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	c := meta.(*client)

	req, err := resourceSettingRadiusGetResourceData(d, meta)
	if err != nil {
		return diag.FromErr(err)
	}

	site := d.Get("site").(string)
	if site == "" {
		site = c.site
	}

	resp, err := c.c.UpdateSettingRadius(ctx, site, req)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(resp.ID)

	return resourceSettingRadiusSetResourceData(resp, d, meta, site)
}

func resourceSettingRadiusSetResourceData(resp *unifi.SettingRadius, d *schema.ResourceData, meta any, site string) diag.Diagnostics {
	d.Set("site", site)
	d.Set("enabled", resp.Enabled)
	d.Set("accounting_enabled", resp.AccountingEnabled)
	d.Set("accounting_port", resp.AcctPort)
	d.Set("auth_port", resp.AuthPort)
	d.Set("tunneled_reply", resp.TunneledReply)
	d.Set("secret", resp.XSecret)
	d.Set("interim_update_interval", resp.InterimUpdateInterval)
	return nil
}

func resourceSettingRadiusRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	c := meta.(*client)

	site := d.Get("site").(string)
	if site == "" {
		site = c.site
	}

	resp, err := c.c.GetSettingRadius(ctx, site)
	if _, ok := err.(*unifi.NotFoundError); ok {
		d.SetId("")
		return nil
	}
	if err != nil {
		return diag.FromErr(err)
	}

	return resourceSettingRadiusSetResourceData(resp, d, meta, site)
}

func resourceSettingRadiusUpdate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	c := meta.(*client)

	req, err := resourceSettingRadiusGetResourceData(d, meta)
	if err != nil {
		return diag.FromErr(err)
	}

	req.ID = d.Id()
	site := d.Get("site").(string)
	if site == "" {
		site = c.site
	}

	resp, err := c.c.UpdateSettingRadius(ctx, site, req)
	if err != nil {
		return diag.FromErr(err)
	}

	return resourceSettingRadiusSetResourceData(resp, d, meta, site)
}
