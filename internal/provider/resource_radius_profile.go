package provider

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/paultyng/go-unifi/unifi"
	"strings"
)

func resourceRadiusProfile() *schema.Resource {
	return &schema.Resource{
		Description: "`unifi_setting_radius` manages settings for the built-in RADIUS server.",

		CreateContext: resourceRadiusProfileCreate,
		ReadContext:   resourceRadiusProfileRead,
		UpdateContext: resourceRadiusProfileUpdate,
		DeleteContext: resourceRadiusProfileDelete,
		Importer: &schema.ResourceImporter{
			StateContext: importRadiusProfile,
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
			"name": {
				Description: "The name of the profile.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"accounting_enabled": {
				Description: "Specifies whether to use radius accounting.",
				Type:        schema.TypeBool,
				Default:     false,
				Optional:    true,
			},
			"interim_update_enabled": {
				Description: "Specifies whether to use interim_update.",
				Type:        schema.TypeBool,
				Default:     false,
				Optional:    true,
			},
			"interim_update_interval": {
				Description: "Specifies interim_update interval.",
				Type:        schema.TypeInt,
				Default:     3600,
				Optional:    true,
			},
			"use_usg_acct_server": {
				Description: "Specifies whether to use usg as a radius accounting server.",
				Type:        schema.TypeBool,
				Default:     false,
				Optional:    true,
			},
			"use_usg_auth_server": {
				Description: "Specifies whether to use usg as a radius authentication server.",
				Type:        schema.TypeBool,
				Default:     false,
				Optional:    true,
			},
			"vlan_enabled": {
				Description: "Specifies whether to use vlan on wired connections.",
				Type:        schema.TypeBool,
				Default:     false,
				Optional:    true,
			},
			"vlan_wlan_mode": {
				Description:  "Specifies whether to use vlan on wireless connections. Must be one of `disabled`, `optional`, or `required`.",
				Type:         schema.TypeString,
				Default:      "",
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"disabled", "optional", "required"}, false),
			},
			"auth_server": {
				Description: "RADIUS authentication servers.",
				Type:        schema.TypeList,
				Optional:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"ip": {
							Description:  "IP address of authentication service server.",
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.IsIPAddress,
						},
						"port": {
							Description:  "Port of authentication service.",
							Type:         schema.TypeInt,
							Optional:     true,
							Default:      1812,
							ValidateFunc: validation.IsPortNumber,
						},
						"xsecret": {
							Description: "RADIUS secret.",
							Type:        schema.TypeString,
							Required:    true,
							Sensitive:   true,
						},
					},
				},
			},
			"acct_server": {
				Description: "RADIUS accounting servers.",
				Type:        schema.TypeList,
				Optional:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"ip": {
							Description:  "IP address of accounting service server.",
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.IsIPAddress,
						},
						"port": {
							Description:  "Port of accounting service.",
							Type:         schema.TypeInt,
							Optional:     true,
							Default:      1813,
							ValidateFunc: validation.IsPortNumber,
						},
						"xsecret": {
							Description: "RADIUS secret.",
							Type:        schema.TypeString,
							Required:    true,
							Sensitive:   true,
						},
					},
				},
			},
		},
	}
}

func setToAuthServers(set []interface{}) ([]unifi.RADIUSProfileAuthServers, error) {
	var authServers []unifi.RADIUSProfileAuthServers
	for _, item := range set {
		data, ok := item.(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("unexpected data in block")
		}
		authServer, err := toAuthServer(data)
		if err != nil {
			return nil, fmt.Errorf("unable to create port override: %w", err)
		}
		authServers = append(authServers, authServer)
	}
	return authServers, nil
}

func setToAcctServers(set []interface{}) ([]unifi.RADIUSProfileAcctServers, error) {
	var acctServers []unifi.RADIUSProfileAcctServers
	for _, item := range set {
		data, ok := item.(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("unexpected data in block")
		}
		accServer, err := toAcctServer(data)
		if err != nil {
			return nil, fmt.Errorf("unable to create port override: %w", err)
		}
		acctServers = append(acctServers, accServer)
	}
	return acctServers, nil
}

func toAuthServer(data map[string]interface{}) (unifi.RADIUSProfileAuthServers, error) {
	return unifi.RADIUSProfileAuthServers{
		IP:      data["ip"].(string),
		Port:    data["port"].(int),
		XSecret: data["xsecret"].(string),
	}, nil
}

func toAcctServer(data map[string]interface{}) (unifi.RADIUSProfileAcctServers, error) {
	return unifi.RADIUSProfileAcctServers{
		IP:      data["ip"].(string),
		Port:    data["port"].(int),
		XSecret: data["xsecret"].(string),
	}, nil
}

func setFromAuthServers(authServers []unifi.RADIUSProfileAuthServers) ([]map[string]interface{}, error) {
	list := make([]map[string]interface{}, 0, len(authServers))
	for _, authServer := range authServers {
		v, err := fromAuthServer(authServer)
		if err != nil {
			return nil, fmt.Errorf("unable to parse ssh key: %w", err)
		}
		list = append(list, v)
	}
	return list, nil
}

func setFromAcctServers(acctServers []unifi.RADIUSProfileAcctServers) ([]map[string]interface{}, error) {
	list := make([]map[string]interface{}, 0, len(acctServers))
	for _, acctServer := range acctServers {
		v, err := fromAcctServer(acctServer)
		if err != nil {
			return nil, fmt.Errorf("unable to parse ssh key: %w", err)
		}
		list = append(list, v)
	}
	return list, nil
}

func fromAuthServer(sshKey unifi.RADIUSProfileAuthServers) (map[string]interface{}, error) {
	return map[string]interface{}{
		"ip":      sshKey.IP,
		"port":    sshKey.Port,
		"xsecret": sshKey.XSecret,
	}, nil
}

func fromAcctServer(sshKey unifi.RADIUSProfileAcctServers) (map[string]interface{}, error) {
	return map[string]interface{}{
		"ip":      sshKey.IP,
		"port":    sshKey.Port,
		"xsecret": sshKey.XSecret,
	}, nil
}

func resourceRadiusProfileCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*client)
	req, err := resourceRadiusProfileGetResourceData(d)
	if err != nil {
		return diag.FromErr(err)
	}

	site := d.Get("site").(string)
	if site == "" {
		site = c.site
	}
	resp, err := c.c.CreateRADIUSProfile(ctx, site, req)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(resp.ID)

	return resourceRadiusProfileSetResourceData(resp, d, site)
}

func resourceRadiusProfileGetResourceData(d *schema.ResourceData) (*unifi.RADIUSProfile, error) {
	authServers, err := setToAuthServers(d.Get("auth_server").([]interface{}))
	if err != nil {
		return nil, fmt.Errorf("unable to auth_server ssh_key block: %w", err)
	}
	acctServers, err := setToAcctServers(d.Get("acct_server").([]interface{}))
	if err != nil {
		return nil, fmt.Errorf("unable to acct_server ssh_key block: %w", err)
	}
	return &unifi.RADIUSProfile{
		Name:                  d.Get("name").(string),
		InterimUpdateEnabled:  d.Get("interim_update_enabled").(bool),
		InterimUpdateInterval: d.Get("interim_update_interval").(int),
		AccountingEnabled:     d.Get("accounting_enabled").(bool),
		UseUsgAcctServer:      d.Get("use_usg_acct_server").(bool),
		UseUsgAuthServer:      d.Get("use_usg_auth_server").(bool),
		VLANEnabled:           d.Get("vlan_enabled").(bool),
		VLANWLANMode:          d.Get("vlan_wlan_mode").(string),
		AuthServers:           authServers,
		AcctServers:           acctServers,
	}, nil
}

func resourceRadiusProfileSetResourceData(resp *unifi.RADIUSProfile, d *schema.ResourceData, site string) diag.Diagnostics {
	authServers, err := setFromAuthServers(resp.AuthServers)
	if err != nil {
		return diag.FromErr(err)
	}
	acctServers, err := setFromAcctServers(resp.AcctServers)
	if err != nil {
		return diag.FromErr(err)
	}

	d.Set("site", site)
	d.Set("name", resp.Name)

	d.Set("interim_update_enabled", resp.InterimUpdateEnabled)
	d.Set("interim_update_interval", resp.InterimUpdateInterval)
	d.Set("accounting_enabled", resp.AccountingEnabled)
	d.Set("use_usg_acct_server", resp.UseUsgAcctServer)
	d.Set("use_usg_auth_server", resp.UseUsgAuthServer)
	d.Set("vlan_enabled", resp.VLANEnabled)
	d.Set("vlan_wlan_mode", resp.VLANWLANMode)
	d.Set("auth_server", authServers)
	d.Set("acct_server", acctServers)
	return nil
}

func resourceRadiusProfileRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*client)

	id := d.Id()

	site := d.Get("site").(string)
	if site == "" {
		site = c.site
	}
	resp, err := c.c.GetRADIUSProfile(ctx, site, id)
	if _, ok := err.(*unifi.NotFoundError); ok {
		d.SetId("")
		return nil
	}
	if err != nil {
		return diag.FromErr(err)
	}

	return resourceRadiusProfileSetResourceData(resp, d, site)
}

func resourceRadiusProfileUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*client)

	req, err := resourceRadiusProfileGetResourceData(d)
	if err != nil {
		return diag.FromErr(err)
	}

	req.ID = d.Id()

	site := d.Get("site").(string)
	if site == "" {
		site = c.site
	}
	req.SiteID = site

	resp, err := c.c.UpdateRADIUSProfile(ctx, site, req)
	if err != nil {
		return diag.FromErr(err)
	}

	return resourceRadiusProfileSetResourceData(resp, d, site)
}

func resourceRadiusProfileDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*client)

	id := d.Id()

	site := d.Get("site").(string)
	if site == "" {
		site = c.site
	}

	err := c.c.DeleteRADIUSProfile(ctx, site, id)
	return diag.FromErr(err)
}

func importRadiusProfile(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	c := meta.(*client)
	id := d.Id()
	site := d.Get("site").(string)
	if site == "" {
		site = c.site
	}

	if strings.Contains(id, ":") {
		importParts := strings.SplitN(id, ":", 2)
		site = importParts[0]
		id = importParts[1]
	}

	if strings.HasPrefix(id, "name=") {
		targetName := strings.TrimPrefix(id, "name=")
		var err error
		if id, err = getRadiusProfileIDByName(ctx, c.c, targetName, site); err != nil {
			return nil, err
		}
	}

	if id != "" {
		d.SetId(id)
	}
	if site != "" {
		d.Set("site", site)
	}

	return []*schema.ResourceData{d}, nil
}

func getRadiusProfileIDByName(ctx context.Context, client unifiClient, profileName, site string) (string, error) {
	radiusProfiles, err := client.ListRADIUSProfile(ctx, site)
	if err != nil {
		return "", err
	}

	idMatchingName := ""
	allNames := []string{}
	for _, profile := range radiusProfiles {
		allNames = append(allNames, profile.Name)
		if profile.Name != profileName {
			continue
		}
		if idMatchingName != "" {
			return "", fmt.Errorf("Found multiple radius profiles with name '%s'", profileName)
		}
		idMatchingName = profile.ID
	}
	if idMatchingName == "" {
		return "", fmt.Errorf("Found no radius profile with name '%s', found: %s", profileName, strings.Join(allNames, ", "))
	}
	return idMatchingName, nil
}
