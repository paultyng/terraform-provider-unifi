package provider

import (
	"context"
	"slices"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/ubiquiti-community/go-unifi/unifi"
)

func resourceDNSRecord() *schema.Resource {
	return &schema.Resource{
		Description: "`unifi_dns_record` manages DNS record settings for different providers.",

		CreateContext: resourceDNSRecordCreate,
		ReadContext:   resourceDNSRecordRead,
		UpdateContext: resourceDNSRecordUpdate,
		DeleteContext: resourceDNSRecordDelete,
		Importer: &schema.ResourceImporter{
			StateContext: importSiteAndID,
		},

		Schema: map[string]*schema.Schema{
			"id": {
				Description: "The ID of the DNS record.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"site": {
				Description: "The name of the site to associate the DNS record with.",
				Type:        schema.TypeString,
				Computed:    true,
				Optional:    true,
				ForceNew:    true,
			},
			"name": {
				Description: "The key of the DNS record.",
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
			},
			"enabled": {
				Description: "Whether the DNS record is enabled.",
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				ForceNew:    false,
			},
			"port": {
				Description: "The port of the DNS record.",
				Type:        schema.TypeInt,
				Required:    true,
			},
			"priority": {
				Description: "The priority of the DNS record.",
				Type:        schema.TypeInt,
				Optional:    true,
			},
			"record_type": {
				Description: "The type of the DNS record.",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"ttl": {
				Description: "The TTL of the DNS record.",
				Type:        schema.TypeInt,
				Optional:    true,
			},
			"value": {
				Description: "The value of the DNS record.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"weight": {
				Description: "The weight of the DNS record.",
				Type:        schema.TypeInt,
				Optional:    true,
			},
		},
	}
}

func resourceDNSRecordCreate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	c := meta.(*client)

	req, err := resourceDNSRecordGetResourceData(d)
	if err != nil {
		return diag.FromErr(err)
	}

	site := d.Get("site").(string)
	if site == "" {
		site = c.site
	}

	resp, err := c.c.CreateDNSRecord(ctx, site, req)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(resp.ID)

	return resourceDNSRecordSetResourceData(resp, d, site)
}

func resourceDNSRecordGetResourceData(d *schema.ResourceData) (*unifi.DNSRecord, error) {
	r := &unifi.DNSRecord{
		Enabled:    d.Get("enabled").(bool),
		Key:        d.Get("name").(string),
		Port:       d.Get("port").(int),
		Priority:   d.Get("priority").(int),
		RecordType: d.Get("record_type").(string),
		Ttl:        d.Get("ttl").(int),
		Value:      d.Get("value").(string),
		Weight:     d.Get("weight").(int),
	}

	return r, nil
}

func resourceDNSRecordSetResourceData(resp *unifi.DNSRecord, d *schema.ResourceData, site string) diag.Diagnostics {
	d.Set("enabled", resp.Enabled)
	d.Set("name", resp.Key)
	d.Set("port", resp.Port)
	d.Set("priority", resp.Priority)
	d.Set("record_type", resp.RecordType)
	d.Set("ttl", resp.Ttl)
	d.Set("value", resp.Value)
	d.Set("weight", resp.Weight)
	d.Set("site", site)

	return nil
}

func resourceDNSRecordRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	c := meta.(*client)

	id := d.Id()

	site := d.Get("site").(string)
	if site == "" {
		site = c.site
	}

	resp, err := c.c.ListDNSRecord(ctx, site)

	if err != nil {
		return diag.FromErr(err)
	}

	i := slices.IndexFunc(resp, func(r unifi.DNSRecord) bool {
		return r.ID == id
	})

	if i == -1 {
		d.SetId("")
		return nil
	}

	rec := resp[i]

	return resourceDNSRecordSetResourceData(&rec, d, site)
}

func resourceDNSRecordUpdate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	c := meta.(*client)

	req, err := resourceDNSRecordGetResourceData(d)
	if err != nil {
		return diag.FromErr(err)
	}

	req.ID = d.Id()

	site := d.Get("site").(string)
	if site == "" {
		site = c.site
	}
	req.SiteID = site

	resp, err := c.c.UpdateDNSRecord(ctx, site, req)
	if err != nil {
		return diag.FromErr(err)
	}

	return resourceDNSRecordSetResourceData(resp, d, site)
}

func resourceDNSRecordDelete(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	c := meta.(*client)

	id := d.Id()

	site := d.Get("site").(string)
	if site == "" {
		site = c.site
	}
	err := c.c.DeleteDNSRecord(ctx, site, id)
	if _, ok := err.(*unifi.NotFoundError); ok {
		return nil
	}
	return diag.FromErr(err)
}
