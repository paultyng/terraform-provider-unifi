package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/paultyng/go-unifi/unifi"
)

func resourceUser() *schema.Resource {
	return &schema.Resource{
		Description: "`unifi_user` manages a user (or \"client\" in the UI) of the network, these are identified " +
			"by unique MAC addresses.\n\n" +
			"Users are created in the controller when observed on the network, so the resource defaults to allowing " +
			"itself to just take over management of a MAC address, but this can be turned off.",

		Create: resourceUserCreate,
		Read:   resourceUserRead,
		Update: resourceUserUpdate,
		Delete: resourceUserDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"id": {
				Description: "The ID of the user.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"site": {
				Description: "The name of the site to associate the user with.",
				Type:        schema.TypeString,
				Computed:    true,
				Optional:    true,
				ForceNew:    true,
			},
			"mac": {
				Description:      "The MAC address of the user.",
				Type:             schema.TypeString,
				Required:         true,
				ForceNew:         true,
				DiffSuppressFunc: macDiffSuppressFunc,
				ValidateFunc:     validation.StringMatch(macAddressRegexp, "Mac address is invalid"),
			},
			"name": {
				Description: "The name of the user.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"user_group_id": {
				Description: "The user group ID for the user.",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"note": {
				Description: "A note with additional information for the user.",
				Type:        schema.TypeString,
				Optional:    true,
			},
			// TODO: combine this with output IP for a single attribute ip_address?
			"fixed_ip": {
				Description:  "A fixed IPv4 address for this user.",
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.IsIPv4Address,
			},
			"network_id": {
				Description: "The network ID for this user.",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"blocked": {
				Description: "Specifies whether this user should be blocked from the network.",
				Type:        schema.TypeBool,
				Optional:    true,
			},

			// these are "meta" attributes that control TF UX
			"allow_existing": {
				Description: "Specifies whether this resource should just take over control of an existing user.",
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
			},
			"skip_forget_on_destroy": {
				Description: "Specifies whether this resource should tell the controller to \"forget\" the user on destroy.",
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
			},

			// computed only attributes
			"hostname": {
				Description: "The hostname of the user.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"ip": {
				Description: "The IP address of the user.",
				Type:        schema.TypeString,
				Computed:    true,
			},
		},
	}
}

func resourceUserCreate(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*client)

	req, err := resourceUserGetResourceData(d)
	if err != nil {
		return err
	}

	allowExisting := d.Get("allow_existing").(bool)

	site := d.Get("site").(string)
	if site == "" {
		site = c.site
	}

	resp, err := c.c.CreateUser(context.TODO(), site, req)
	if err != nil {
		apiErr, ok := err.(*unifi.APIError)
		if !ok || (apiErr.Message != "api.err.MacUsed" || !allowExisting) {
			return err
		}

		// mac in use, just absorb it
		mac := d.Get("mac").(string)
		existing, err := c.c.GetUserByMAC(context.TODO(), site, mac)
		if err != nil {
			return err
		}

		req.ID = existing.ID
		req.SiteID = existing.SiteID

		resp, err = c.c.UpdateUser(context.TODO(), site, req)
		if err != nil {
			return err
		}
	}

	d.SetId(resp.ID)

	if d.Get("blocked").(bool) {
		err := c.c.BlockUserByMAC(context.TODO(), site, d.Get("mac").(string))
		if err != nil {
			return err
		}
	}

	return resourceUserSetResourceData(resp, d, site)
}

func resourceUserGetResourceData(d *schema.ResourceData) (*unifi.User, error) {
	fixedIP := d.Get("fixed_ip").(string)

	return &unifi.User{
		MAC:         d.Get("mac").(string),
		Name:        d.Get("name").(string),
		UserGroupID: d.Get("user_group_id").(string),
		Note:        d.Get("note").(string),
		FixedIP:     fixedIP,
		UseFixedIP:  fixedIP != "",
		NetworkID:   d.Get("network_id").(string),
		// not sure if this matters/works
		Blocked: d.Get("blocked").(bool),
	}, nil
}

func resourceUserSetResourceData(resp *unifi.User, d *schema.ResourceData, site string) error {
	fixedIP := ""
	if resp.UseFixedIP {
		fixedIP = resp.FixedIP
	}

	d.Set("site", site)
	d.Set("mac", resp.MAC)
	d.Set("name", resp.Name)
	d.Set("user_group_id", resp.UserGroupID)
	d.Set("note", resp.Note)
	d.Set("fixed_ip", fixedIP)
	d.Set("network_id", resp.NetworkID)
	d.Set("blocked", resp.Blocked)

	d.Set("hostname", resp.Hostname)
	d.Set("ip", resp.IP)

	return nil
}

func resourceUserRead(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*client)

	id := d.Id()

	site := d.Get("site").(string)
	if site == "" {
		site = c.site
	}

	resp, err := c.c.GetUser(context.TODO(), site, id)
	if _, ok := err.(*unifi.NotFoundError); ok {
		d.SetId("")
		return nil
	}
	if err != nil {
		return err
	}

	// for some reason the IP address is only on this endpoint, so issue another request
	macResp, err := c.c.GetUserByMAC(context.TODO(), site, resp.MAC)
	if _, ok := err.(*unifi.NotFoundError); ok {
		d.SetId("")
		return nil
	}
	if err != nil {
		return err
	}

	resp.IP = macResp.IP

	return resourceUserSetResourceData(resp, d, site)
}

func resourceUserUpdate(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*client)

	site := d.Get("site").(string)
	if site == "" {
		site = c.site
	}

	if d.HasChange("blocked") {
		mac := d.Get("mac").(string)
		if d.Get("blocked").(bool) {
			err := c.c.BlockUserByMAC(context.TODO(), site, mac)
			if err != nil {
				return err
			}
		} else {
			err := c.c.UnblockUserByMAC(context.TODO(), site, mac)
			if err != nil {
				return err
			}
		}
	}

	req, err := resourceUserGetResourceData(d)
	if err != nil {
		return err
	}

	req.ID = d.Id()
	req.SiteID = site

	resp, err := c.c.UpdateUser(context.TODO(), site, req)
	if err != nil {
		return err
	}

	return resourceUserSetResourceData(resp, d, site)
}

func resourceUserDelete(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*client)

	id := d.Id()

	if d.Get("skip_forget_on_destroy").(bool) {
		return nil
	}

	site := d.Get("site").(string)
	if site == "" {
		site = c.site
	}

	// lookup MAC instead of trusting state
	u, err := c.c.GetUser(context.TODO(), site, id)
	if _, ok := err.(*unifi.NotFoundError); ok {
		return nil
	}
	if err != nil {
		return err
	}

	err = c.c.DeleteUserByMAC(context.TODO(), site, u.MAC)
	return err
}
