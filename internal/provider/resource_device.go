package provider

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/paultyng/go-unifi/unifi"
)

func resourceDevice() *schema.Resource {
	return &schema.Resource{
		Description: "`unifi_device` manages a device of the network.\n\n" +
			"Devices are adopted by the controller, so it is not possible for this resource to be created through " +
			"Terraform, the create operation instead will simply start managing the device specified by MAC address. " +
			"It's safer to start this process with an explicit import of the device.",

		CreateContext: resourceDeviceCreate,
		ReadContext:   resourceDeviceRead,
		UpdateContext: resourceDeviceUpdate,
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
				Description:      "The MAC address of the device. This can be specified so that the provider can take control of a device (since devices are created through adoption).",
				Type:             schema.TypeString,
				Optional:         true,
				Computed:         true,
				ForceNew:         true,
				DiffSuppressFunc: macDiffSuppressFunc,
				ValidateFunc:     validation.StringMatch(macAddressRegexp, "Mac address is invalid"),
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
				// see https://github.com/hashicorp/terraform-plugin-sdk/issues/62
				Type:     schema.TypeSet,
				Optional: true,
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

			"allow_adoption": {
				Description: "Specifies whether this resource should tell the controller to adopt the device on create.",
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
			},
			"forget_on_destroy": {
				Description: "Specifies whether this resource should tell the controller to forget the device on destroy.",
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
			},
		},
	}
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
		mac := cleanMAC(id)
		device, err := c.c.GetDeviceByMAC(ctx, site, mac)

		if err != nil {
			return nil, err
		}

		id = device.ID
	}

	if id != "" {
		d.SetId(id)
	}
	if site != "" {
		d.Set("site", site)
	}

	return []*schema.ResourceData{d}, nil
}

func resourceDeviceCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*client)

	site := d.Get("site").(string)
	if site == "" {
		site = c.site
	}

	mac := d.Get("mac").(string)
	if mac == "" {
		return diag.Errorf("no MAC address specified, please import the device using terraform import")
	}

	mac = cleanMAC(mac)
	device, err := c.c.GetDeviceByMAC(ctx, site, mac)

	if device == nil {
		return diag.Errorf("device not found using mac %q", mac)
	}
	if err != nil {
		return diag.FromErr(err)
	}

	if !device.Adopted {
		if !d.Get("allow_adoption").(bool) {
			return diag.Errorf("Device must be adopted before it can be managed")
		}

		err := c.c.AdoptDevice(ctx, site, mac)
		if err != nil {
			return diag.FromErr(err)
		}

		device, err = waitForDeviceState(ctx, d, meta, unifi.DeviceStateConnected, []unifi.DeviceState{unifi.DeviceStateAdopting, unifi.DeviceStatePending, unifi.DeviceStateProvisioning, unifi.DeviceStateUpgrading}, 2*time.Minute)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	d.SetId(device.ID)
	return resourceDeviceUpdate(ctx, d, meta)
}

func resourceDeviceUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*client)

	site := d.Get("site").(string)
	if site == "" {
		site = c.site
	}

	req, err := resourceDeviceGetResourceData(d)
	if err != nil {
		return diag.FromErr(err)
	}

	req.ID = d.Id()
	req.SiteID = site

	resp, err := c.c.UpdateDevice(ctx, site, req)
	if err != nil {
		return diag.FromErr(err)
	}

	_, err = waitForDeviceState(ctx, d, meta, unifi.DeviceStateConnected, []unifi.DeviceState{unifi.DeviceStateAdopting, unifi.DeviceStateProvisioning}, 30*time.Second)
	if err != nil {
		return diag.FromErr(err)
	}

	return resourceDeviceSetResourceData(resp, d, site)
}

func resourceDeviceDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*client)

	if !d.Get("forget_on_destroy").(bool) {
		return nil
	}

	site := d.Get("site").(string)
	mac := d.Get("mac").(string)

	if site == "" {
		site = c.site
	}

	err := c.c.ForgetDevice(ctx, site, mac)
	if err != nil {
		return diag.FromErr(err)
	}

	_, err = waitForDeviceState(ctx, d, meta, unifi.DeviceStatePending, []unifi.DeviceState{unifi.DeviceStateConnected, unifi.DeviceStateDeleting}, 30*time.Second)
	if _, ok := err.(*unifi.NotFoundError); !ok {
		return diag.FromErr(err)
	}

	return nil
}

func resourceDeviceRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*client)

	id := d.Id()

	site := d.Get("site").(string)
	if site == "" {
		site = c.site
	}

	resp, err := c.c.GetDevice(ctx, site, id)
	if _, ok := err.(*unifi.NotFoundError); ok {
		d.SetId("")
		return nil
	}
	if err != nil {
		return diag.FromErr(err)
	}

	return resourceDeviceSetResourceData(resp, d, site)
}

func resourceDeviceSetResourceData(resp *unifi.Device, d *schema.ResourceData, site string) diag.Diagnostics {
	portOverrides, err := setFromPortOverrides(resp.PortOverrides)
	if err != nil {
		return diag.FromErr(err)
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

func setFromPortOverrides(pos []unifi.DevicePortOverrides) ([]map[string]interface{}, error) {
	list := make([]map[string]interface{}, 0, len(pos))
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

func waitForDeviceState(ctx context.Context, d *schema.ResourceData, meta interface{}, targetState unifi.DeviceState, pendingStates []unifi.DeviceState, timeout time.Duration) (*unifi.Device, error) {
	c := meta.(*client)

	site := d.Get("site").(string)
	mac := d.Get("mac").(string)

	if site == "" {
		site = c.site
	}

	// Always consider unknown to be a pending state.
	pendingStates = append(pendingStates, unifi.DeviceStateUnknown)

	var pending []string
	for _, state := range pendingStates {
		pending = append(pending, state.String())
	}

	wait := resource.StateChangeConf{
		Pending: pending,
		Target:  []string{targetState.String()},
		Refresh: func() (interface{}, string, error) {
			device, err := c.c.GetDeviceByMAC(ctx, site, mac)

			if _, ok := err.(*unifi.NotFoundError); ok {
				err = nil
			}

			// When a device is forgotten, it will disappear from the UI for a few seconds before reappearing.
			// During this time, `device.GetDeviceByMAC` will return a 400.
			//
			// TODO: Improve handling of this situation in `go-unifi`.
			if err != nil && strings.Contains(err.Error(), "api.err.UnknownDevice") {
				err = nil
			}

			var state string
			if device != nil {
				state = device.State.String()
			}

			// TODO: Why is this needed???
			if device == nil {
				return nil, state, err
			}

			return device, state, err
		},
		Timeout:        timeout,
		NotFoundChecks: 30,
	}

	outputRaw, err := wait.WaitForStateContext(ctx)

	if output, ok := outputRaw.(*unifi.Device); ok {
		return output, err
	}

	return nil, err
}
