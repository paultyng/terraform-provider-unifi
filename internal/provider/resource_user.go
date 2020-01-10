package provider

import (
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/paultyng/go-unifi/unifi"
)

func resourceUser() *schema.Resource {
	return &schema.Resource{
		Create: resourceUserCreate,
		Read:   resourceUserRead,
		Update: resourceUserUpdate,
		Delete: resourceUserDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"mac": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					old = strings.TrimSpace(strings.ReplaceAll(strings.ToLower(old), "-", ":"))
					new = strings.TrimSpace(strings.ReplaceAll(strings.ToLower(new), "-", ":"))
					return old == new
				},
				// Validation:
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"user_group_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"note": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"fixed_ip": {
				Type:     schema.TypeString,
				Optional: true,
				// TODO: Validate
			},
			"network_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"blocked": {
				Type:     schema.TypeBool,
				Optional: true,
			},

			// these are "meta" attributes that control TF UX
			"allow_existing": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			"skip_forget_on_destroy": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
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

	resp, err := c.c.CreateUser(c.site, req)
	if err != nil {
		apiErr, ok := err.(*unifi.APIError)
		if !ok || (apiErr.Message != "api.err.MacUsed" || !allowExisting) {
			return err
		}

		// mac in use, just absorb it
		mac := d.Get("mac").(string)
		existing, err := c.c.GetUserByMAC(c.site, mac)
		if err != nil {
			return err
		}

		req.ID = existing.ID
		req.SiteID = existing.SiteID

		resp, err = c.c.UpdateUser(c.site, req)
		if err != nil {
			return err
		}
	}

	d.SetId(resp.ID)

	if d.Get("blocked").(bool) {
		err := c.c.BlockUserByMAC(c.site, d.Get("mac").(string))
		if err != nil {
			return err
		}
	}

	return resourceUserSetResourceData(resp, d)
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

func resourceUserSetResourceData(resp *unifi.User, d *schema.ResourceData) error {
	fixedIP := ""
	if resp.UseFixedIP {
		fixedIP = resp.FixedIP
	}

	d.Set("mac", resp.MAC)
	d.Set("name", resp.Name)
	d.Set("user_group_id", resp.UserGroupID)
	d.Set("note", resp.Note)
	d.Set("fixed_ip", fixedIP)
	d.Set("network_id", resp.NetworkID)
	d.Set("blocked", resp.Blocked)

	return nil
}

func resourceUserRead(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*client)

	id := d.Id()

	resp, err := c.c.GetUser(c.site, id)
	if _, ok := err.(*unifi.NotFoundError); ok {
		d.SetId("")
		return nil
	}
	if err != nil {
		return err
	}

	return resourceUserSetResourceData(resp, d)
}

func resourceUserUpdate(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*client)

	if d.HasChange("blocked") {
		mac := d.Get("mac").(string)
		if d.Get("blocked").(bool) {
			err := c.c.BlockUserByMAC(c.site, mac)
			if err != nil {
				return err
			}
		} else {
			err := c.c.UnblockUserByMAC(c.site, mac)
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
	req.SiteID = c.site

	resp, err := c.c.UpdateUser(c.site, req)
	if err != nil {
		return err
	}

	return resourceUserSetResourceData(resp, d)
}

func resourceUserDelete(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*client)

	id := d.Id()

	if d.Get("skip_forget_on_destroy").(bool) {
		return nil
	}

	// lookup MAC instead of trusting state
	u, err := c.c.GetUser(c.site, id)
	if _, ok := err.(*unifi.NotFoundError); ok {
		return nil
	}
	if err != nil {
		return err
	}

	err = c.c.DeleteUserByMAC(c.site, u.MAC)
	return err
}
