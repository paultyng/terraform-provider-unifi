package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/paultyng/go-unifi/unifi"
)

func resourceDevice() *schema.Resource {
	return &schema.Resource{
		Description: "`unifi_device` manages a device of the network.\n\n" +
			"Devices are adopted by the controller, so it may not be possible " +
			"for this resource to be created through terraform.",

		//Create: resourceDeviceCreate,
		Read:   resourceDeviceRead,
		Update: resourceDeviceUpdate,
		Delete: resourceDeviceDelete,
		Importer: &schema.ResourceImporter{
			State: importSiteAndID,
		},

		Schema: map[string]*schema.Schema{
			"id": {
				Description: "The ID of the device.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"site": {
				Description: "The name of the site to associate the device with.",
				Type:        schema.TypeString,
				Computed:    true,
				Optional:    true,
				ForceNew:    true,
			},
			"mac": {
				Description:      "The MAC address of the device.",
				Type:             schema.TypeString,
				Required:         true,
				ForceNew:         true,
				DiffSuppressFunc: macDiffSuppressFunc,
				ValidateFunc:     validation.StringMatch(macAddressRegexp, "Mac address is invalid"),
			},
			"name": {
				Description: "The name of the device.",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"disabled": {
				Description: "Specifies whether this device should be disabled.",
				Type:        schema.TypeBool,
				Computed:    true,
			},
			"port_overrides": {
				Description: "Settings overrides for specific switch ports.",
				Type:        schema.TypeSet,
				Optional:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"port_idx": {
							Description: "Switch port number.",
							Type:        schema.TypeInt,
							Required:    true,
						},
						"name": {
							Description: "Human-readable name of the port.",
							Type:        schema.TypeString,
							Optional:    true,
						},
						"port_profile_id": {
							Description: "ID of the Port Profile used on this port.",
							Type:        schema.TypeString,
							Optional:    true,
						},
					},
				},
			},
		},
	}
}

func resourceDeviceCreate(d *schema.ResourceData, meta interface{}) error {
	return &unifi.NotFoundError{}
}

func resourceDeviceUpdate(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*client)

	site := d.Get("site").(string)
	if site == "" {
		site = c.site
	}

	req, err := resourceDeviceGetResourceData(d)
	if err != nil {
		return err
	}

	req.ID = d.Id()
	req.SiteID = site

	resp, err := c.c.UpdateDevice(context.TODO(), site, req)
	if err != nil {
		return err
	}

	return resourceDeviceSetResourceData(resp, d, site)
}

func resourceDeviceDelete(d *schema.ResourceData, meta interface{}) error {
	return &unifi.NotFoundError{}
}

func resourceDeviceRead(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*client)

	id := d.Id()

	site := d.Get("site").(string)
	if site == "" {
		site = c.site
	}

	resp, err := c.c.GetDevice(context.TODO(), site, id)
	if _, ok := err.(*unifi.NotFoundError); ok {
		d.SetId("")
		return nil
	}
	if err != nil {
		return err
	}

	return resourceDeviceSetResourceData(resp, d, site)
}

func resourceDeviceSetResourceData(resp *unifi.Device, d *schema.ResourceData, site string) error {
	port_overrides, err := listFromPortOverrides(resp.PortOverrides)
	if err != nil {
		return err
	}

	d.Set("site", site)
	d.Set("mac", resp.MAC)
	d.Set("name", resp.Name)
	d.Set("disabled", resp.Disabled)
	d.Set("port_overrides", port_overrides)

	return nil
}

func resourceDeviceGetResourceData(d *schema.ResourceData) (*unifi.Device, error) {
	pos, err := listToPortOverrides(d.Get("port_overrides").([]interface{}))
	if err != nil {
		return nil, fmt.Errorf("unable to process port_overrides block: %w", err)
	}

	//TODO: pass Disabled once we figure out how to enable the device afterwards

	return &unifi.Device{
		MAC:           d.Get("mac").(string),
		Name:          d.Get("name").(string),
		PortOverrides: pos,
	}, nil
}

func listToPortOverrides(list []interface{}) ([]unifi.DevicePortOverrides, error) {
	pos := make([]unifi.DevicePortOverrides, 0, len(list))
	for _, item := range list {
		data, ok := item.(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("unexpected data in block")
		}
		po, err := toPortOverride(data)
		if err != nil {
			return nil, fmt.Errorf("unable to create port override: %w", err)
		}
		pos = append(pos, po)
	}
	return pos, nil
}

func listFromPortOverrides(pos []unifi.DevicePortOverrides) ([]interface{}, error) {
	list := make([]interface{}, 0, len(pos))
	for _, po := range pos {
		v, err := fromPortOverride(po)
		if err != nil {
			return nil, fmt.Errorf("unable to parse port override: %w", err)
		}
		list = append(list, v)
	}
	return list, nil
}

func toPortOverride(data map[string]interface{}) (unifi.DevicePortOverrides, error) {
	// TODO: error check these?
	idx := data["port_idx"].(int)
	name := data["name"].(string)
	profile_id := data["port_profile_id"].(string)
	return unifi.DevicePortOverrides{
		PortIDX:       idx,
		Name:          name,
		PortProfileID: profile_id,
	}, nil
}

func fromPortOverride(po unifi.DevicePortOverrides) (map[string]interface{}, error) {
	return map[string]interface{}{
		"port_idx":        po.PortIDX,
		"name":            po.Name,
		"port_profile_id": po.PortProfileID,
	}, nil
}
