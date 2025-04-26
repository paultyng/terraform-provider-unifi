package provider

import (
	"context"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func dataUser() *schema.Resource {
	return &schema.Resource{
		Description: "`unifi_user` retrieves properties of a user (or \"client\" in the UI) of the network by MAC address.",

		ReadContext: dataUserRead,

		Schema: map[string]*schema.Schema{
			"site": {
				Description: "The name of the site the user is associated with.",
				Type:        schema.TypeString,
				Computed:    true,
				Optional:    true,
			},
			"mac": {
				Description:      "The MAC address of the user.",
				Type:             schema.TypeString,
				Required:         true,
				DiffSuppressFunc: macDiffSuppressFunc,
				ValidateFunc:     validation.StringMatch(macAddressRegexp, "Mac address is invalid"),
			},

			// read-only / computed
			"id": {
				Description: "The ID of the user.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"name": {
				Description: "The name of the user.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"user_group_id": {
				Description: "The user group ID for the user.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"note": {
				Description: "A note with additional information for the user.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"fixed_ip": {
				Description: "fixed IPv4 address set for this user.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"network_id": {
				Description: "The network ID for this user.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"blocked": {
				Description: "Specifies whether this user should be blocked from the network.",
				Type:        schema.TypeBool,
				Computed:    true,
			},
			"dev_id_override": {
				Description: "Override the device fingerprint.",
				Type:        schema.TypeInt,
				Computed:    true,
			},
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
			"local_dns_record": {
				Description: "The local DNS record for this user.",
				Type:        schema.TypeString,
				Computed:    true,
			},
		},
	}
}

func dataUserRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	c := meta.(*client)

	site := d.Get("site").(string)
	if site == "" {
		site = c.site
	}
	mac := d.Get("mac").(string)

	macResp, err := c.c.GetUserByMAC(ctx, site, strings.ToLower(mac))
	if err != nil {
		return diag.FromErr(err)
	}

	resp, err := c.c.GetUser(ctx, site, macResp.ID)
	if err != nil {
		return diag.FromErr(err)
	}

	// for some reason the IP address is only on this endpoint, so issue another request

	resp.IP = macResp.IP
	fixedIP := ""
	if resp.UseFixedIP {
		fixedIP = resp.FixedIP
	}
	localDnsRecord := ""
	if resp.LocalDNSRecordEnabled {
		localDnsRecord = resp.LocalDNSRecord
	}
	d.SetId(resp.ID)
	d.Set("site", site)
	d.Set("mac", resp.MAC)
	d.Set("name", resp.Name)
	d.Set("user_group_id", resp.UserGroupID)
	d.Set("note", resp.Note)
	d.Set("fixed_ip", fixedIP)
	d.Set("network_id", resp.NetworkID)
	d.Set("blocked", resp.Blocked)
	d.Set("dev_id_override", resp.DevIdOverride)
	d.Set("hostname", resp.Hostname)
	d.Set("ip", resp.IP)
	d.Set("ip", localDnsRecord)

	return nil
}
