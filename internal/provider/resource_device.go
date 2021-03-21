package provider

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/paultyng/go-unifi/unifi"
)

func resourceDevice() *schema.Resource {
	return &schema.Resource{
		Description: "`unifi_device` manages a device of the network.\n\n" +
			"Devices are adopted by the controller, so it is not possible " +
			"for this resource to be created through Terraform.",

		Create:        resourceDeviceCreate,
		Read:          resourceDeviceRead,
		Update:        resourceDeviceUpdate,
		DeleteContext: resourceDeviceDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceDeviceImport,
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
				Description: "The MAC address of the device.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"name": {
				Description: "The name of the device.",
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
			},
			"disabled": {
				Description: "Specifies whether this device should be disabled.",
				Type:        schema.TypeBool,
				Computed:    true,
			},
			"port_override": {
				Description: "Settings overrides for specific switch ports.",
				// TODO: this should really be a map or something when possible in the SDK
				Type:     schema.TypeSet,
				Optional: true,
				Set:      resourceDevicePortOverrideSet,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"number": {
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

func resourceDevicePortOverrideSet(v interface{}) int {
	m := v.(map[string]interface{})
	return m["number"].(int)
}

func resourceDeviceImport(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	c := meta.(*client)
	id := d.Id()
	site := d.Get("site").(string)
	if site == "" {
		site = c.site
	}

	if colons := strings.Count(id, ":"); colons == 1 || colons == 6 {
		importParts := strings.SplitN(id, ":", 2)
		site = importParts[0]
		id = importParts[1]
	}

	if macAddressRegexp.MatchString(id) {
		// look up id by mac
		find := cleanMAC(id)

		devices, err := c.c.ListDevice(ctx, site)
		if err != nil {
			return nil, err
		}
		for _, d := range devices {
			if cleanMAC(d.MAC) == find {
				id = d.ID
				break
			}
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

func resourceDeviceCreate(d *schema.ResourceData, meta interface{}) error {
	return errors.New("unifi_device can only be imported, not created")
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

func resourceDeviceDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  "Deleting a device via Terraform is not supported, the device will just be removed from state.",
		},
	}
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
	portOverrides, err := setFromPortOverrides(resp.PortOverrides)
	if err != nil {
		return err
	}

	d.Set("site", site)
	d.Set("mac", resp.MAC)
	d.Set("name", resp.Name)
	d.Set("disabled", resp.Disabled)
	d.Set("port_override", portOverrides)

	return nil
}

func resourceDeviceGetResourceData(d *schema.ResourceData) (*unifi.Device, error) {
	pos, err := setToPortOverrides(d.Get("port_override").(*schema.Set))
	if err != nil {
		return nil, fmt.Errorf("unable to process port_override block: %w", err)
	}

	//TODO: pass Disabled once we figure out how to enable the device afterwards

	return &unifi.Device{
		MAC:           d.Get("mac").(string),
		Name:          d.Get("name").(string),
		PortOverrides: pos,
	}, nil
}

func setToPortOverrides(set *schema.Set) ([]unifi.DevicePortOverrides, error) {
	// use a map here to remove any duplication
	overrideMap := map[int]unifi.DevicePortOverrides{}
	for _, item := range set.List() {
		data, ok := item.(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("unexpected data in block")
		}
		po, err := toPortOverride(data)
		if err != nil {
			return nil, fmt.Errorf("unable to create port override: %w", err)
		}
		overrideMap[po.PortIDX] = po
	}

	pos := make([]unifi.DevicePortOverrides, 0, len(overrideMap))
	for _, item := range overrideMap {
		pos = append(pos, item)
	}
	return pos, nil
}

func setFromPortOverrides(pos []unifi.DevicePortOverrides) (*schema.Set, error) {
	list := make([]interface{}, 0, len(pos))
	for _, po := range pos {
		v, err := fromPortOverride(po)
		if err != nil {
			return nil, fmt.Errorf("unable to parse port override: %w", err)
		}
		list = append(list, v)
	}
	return schema.NewSet(resourceDevicePortOverrideSet, list), nil
}

func toPortOverride(data map[string]interface{}) (unifi.DevicePortOverrides, error) {
	// TODO: error check these?
	idx := data["number"].(int)
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
		"number":          po.PortIDX,
		"name":            po.Name,
		"port_profile_id": po.PortProfileID,
	}, nil
}
