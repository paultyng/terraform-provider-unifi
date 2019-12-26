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

		// TODO: handle site + ID (or name)
		// Importer: &schema.ResourceImporter{
		// 	State: schema.ImportStatePassthrough,
		// },

		Schema: map[string]*schema.Schema{
			"site": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "default",
				ForceNew: true,
			},
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
		},
	}
}

func resourceWLANCreate(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*client)

	site := d.Get("site").(string)

	// TODO: allow passing these defaults
	wlanGroups, err := c.c.ListWLANGroup(site)
	if err != nil {
		return err
	}
	var defaultWLANGroup *unifi.WLANGroup
	for _, wg := range wlanGroups {
		if wg.HiddenID == "Default" {
			defaultWLANGroup = &wg
			break
		}
	}
	if defaultWLANGroup == nil {
		return fmt.Errorf("unable to find default WLAN group")
	}

	userGroups, err := c.c.ListUserGroup(site)
	if err != nil {
		return err
	}
	var defaultUserGroup *unifi.UserGroup
	for _, ug := range userGroups {
		if ug.HiddenID == "Default" {
			defaultUserGroup = &ug
			break
		}
	}
	if defaultUserGroup == nil {
		return fmt.Errorf("unable to find default user group")
	}

	req := &unifi.WLAN{
		Name:        d.Get("name").(string),
		VLAN:        fmt.Sprintf("%d", d.Get("vlan_id").(int)),
		XPassphrase: d.Get("passphrase").(string),

		WLANGroupID: defaultWLANGroup.ID,
		UserGroupID: defaultUserGroup.ID,

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

	resp, err := c.c.CreateWLAN(site, req)
	if err != nil {
		return err
	}

	d.SetId(resp.ID)

	return nil
}

func resourceWLANRead(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*client)

	site := d.Get("site").(string)
	id := d.Id()

	_, err := c.c.GetWLAN(site, id)
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

	site := d.Get("site").(string)
	id := d.Id()

	err := c.c.DeleteWLAN(site, id)
	if _, ok := err.(*unifi.NotFoundError); ok {
		return nil
	}
	return err
}
