package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/paultyng/go-unifi/unifi"
)

func resourceAccount() *schema.Resource {
	return &schema.Resource{
		Description: "`unifi_account` manages a radius user account\n\n" +
			"To authenticate devices based on MAC address, use the MAC address as the username and password under client creation. \n" +
			"Convert lowercase letters to uppercase, and also remove colons or periods from the MAC address. \n\n" +
			"ATTENTION: If the user profile does not include a VLAN, the client will fall back to the untagged VLAN. \n\n" +
			"NOTE: MAC-based authentication accounts can only be used for wireless and wired clients. L2TP remote access does not apply.",

		CreateContext: resourceAccountCreate,
		ReadContext:   resourceAccountRead,
		UpdateContext: resourceAccountUpdate,
		DeleteContext: resourceAccountDelete,
		Importer: &schema.ResourceImporter{
			StateContext: importSiteAndID,
		},

		Schema: map[string]*schema.Schema{
			"id": {
				Description: "The ID of the account.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"site": {
				Description: "The name of the site to associate the account with.",
				Type:        schema.TypeString,
				Computed:    true,
				Optional:    true,
				ForceNew:    true,
			},
			"name": {
				Description: "The name of the account.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"password": {
				Description: "The password of the account.",
				Type:        schema.TypeString,
				Required:    true,
				Sensitive:   true,
			},
			"tunnel_type": {
				Description:  "See [RFC 2868](https://www.rfc-editor.org/rfc/rfc2868) section 3.1", // @TODO: better documentation https://help.ui.com/hc/en-us/articles/360015268353-UniFi-USG-UDM-Configuring-RADIUS-Server#6
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      13,
				ValidateFunc: validation.IntBetween(1, 13),
			},
			"tunnel_medium_type": {
				Description:  "See [RFC 2868](https://www.rfc-editor.org/rfc/rfc2868) section 3.2", // @TODO: better documentation https://help.ui.com/hc/en-us/articles/360015268353-UniFi-USG-UDM-Configuring-RADIUS-Server#6
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      6,
				ValidateFunc: validation.IntBetween(1, 15),
			},
			"network_id": {
				Description: "ID of the network for this account",
				Type:        schema.TypeString,
				Optional:    true,
			},
		},
	}
}

func resourceAccountCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*client)

	req, err := resourceAccountGetResourceData(d)
	if err != nil {
		return diag.FromErr(err)
	}

	site := d.Get("site").(string)
	if site == "" {
		site = c.site
	}

	resp, err := c.c.CreateAccount(ctx, site, req)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(resp.ID)

	return resourceAccountSetResourceData(resp, d, site)
}

func resourceAccountUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*client)

	site := d.Get("site").(string)
	if site == "" {
		site = c.site
	}

	req, err := resourceAccountGetResourceData(d)
	if err != nil {
		return diag.FromErr(err)
	}

	req.ID = d.Id()
	req.SiteID = site

	resp, err := c.c.UpdateAccount(ctx, site, req)
	if err != nil {
		return diag.FromErr(err)
	}

	return resourceAccountSetResourceData(resp, d, site)
}

func resourceAccountDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*client)

	//name := d.Get("name").(string)
	site := d.Get("site").(string)
	if site == "" {
		site = c.site
	}

	id := d.Id()
	err := c.c.DeleteAccount(ctx, site, id)
	if _, ok := err.(*unifi.NotFoundError); ok {
		return nil
	}
	return diag.FromErr(err)
}

func resourceAccountRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*client)

	id := d.Id()

	site := d.Get("site").(string)
	if site == "" {
		site = c.site
	}

	resp, err := c.c.GetAccount(ctx, site, id)
	if _, ok := err.(*unifi.NotFoundError); ok {
		d.SetId("")
		return nil
	}
	if err != nil {
		return diag.FromErr(err)
	}

	return resourceAccountSetResourceData(resp, d, site)
}

func resourceAccountSetResourceData(resp *unifi.Account, d *schema.ResourceData, site string) diag.Diagnostics {
	d.Set("site", site)
	d.Set("name", resp.Name)
	d.Set("password", resp.XPassword)
	d.Set("tunnel_type", resp.TunnelType)
	d.Set("tunnel_medium_type", resp.TunnelMediumType)
	d.Set("network_id", resp.NetworkID)
	return nil
}

func resourceAccountGetResourceData(d *schema.ResourceData) (*unifi.Account, error) {
	return &unifi.Account{
		Name:             d.Get("name").(string),
		XPassword:        d.Get("password").(string),
		TunnelType:       d.Get("tunnel_type").(int),
		TunnelMediumType: d.Get("tunnel_medium_type").(int),
		NetworkID:        d.Get("network_id").(string),
	}, nil
}
