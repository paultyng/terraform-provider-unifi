package provider

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"

	"github.com/paultyng/terraform-provider-unifi/unifi"
)

func resourceWLAN() *schema.Resource {
	return &schema.Resource{
		Create: resourceWLANCreate,
		Read:   resourceWLANRead,
		Update: resourceWLANUpdate,
		Delete: resourceWLANDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"vlan_id": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  1,
			},
			"passphrase": {
				Type:      schema.TypeString,
				Required:  true,
				Sensitive: true,
			},
			"wlan_group_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"user_group_id": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func resourceWLANCreate(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*client)

	req := &unifi.WLAN{
		Name:        d.Get("name").(string),
		VLAN:        fmt.Sprintf("%d", d.Get("vlan_id").(int)),
		XPassphrase: d.Get("passphrase").(string),

		WLANGroupID: d.Get("wlan_group_id").(string),
		UserGroupID: d.Get("user_group_id").(string),

		Enabled:                  true,
		VLANEnabled:              true,
		WPAEnc:                   "ccmp",
		Security:                 "wpapsk",
		WPAMode:                  "wpa2",
		NameCombineEnabled:       true,
		GroupRekey:               3600,
		DTIMMode:                 "default",
		No2GhzOui:                true,
		MinrateNaBeaconRateKbps:  6000,
		MinrateNaDataRateKbps:    6000,
		MinrateNaMgmtRateKbps:    6000,
		MinrateNgBeaconRateKbps:  1000,
		MinrateNgCckRatesEnabled: true,
		MinrateNgDataRateKbps:    1000,
		MinrateNgMgmtRateKbps:    1000,
	}

	resp, err := c.c.CreateWLAN(c.site, req)
	if err != nil {
		return err
	}

	d.SetId(resp.ID)

	return nil
}

func resourceWLANRead(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*client)

	id := d.Id()

	_, err := c.c.GetWLAN(c.site, id)
	if _, ok := err.(*unifi.NotFoundError); ok {
		d.SetId("")
		return nil
	}
	if err != nil {
		return err
	}

	return nil
}

func resourceWLANUpdate(d *schema.ResourceData, meta interface{}) error {
	panic("not implemented")
}

func resourceWLANDelete(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*client)

	id := d.Id()

	err := c.c.DeleteWLAN(c.site, id)
	if _, ok := err.(*unifi.NotFoundError); ok {
		return nil
	}
	return err
}
