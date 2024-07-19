package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/paultyng/go-unifi/unifi"
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

func resourceDNSRecordCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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
		Key:        d.Get("key").(string),
		Port:       d.Get("port").(int),
		Priority:   d.Get("priority").(int),
		RecordType: d.Get("record_time").(string),
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
	d.Set("record_time", resp.RecordType)
	d.Set("ttl", resp.Ttl)
	d.Set("value", resp.Value)
	d.Set("weight", resp.Weight)

	return nil
}

func resourceDNSRecordRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*client)

	id := d.Id()

	site := d.Get("site").(string)
	if site == "" {
		site = c.site
	}

	resp, err := c.c.GetDNSRecord(ctx, site, id)
	if _, ok := err.(*unifi.NotFoundError); ok {
		d.SetId("")
		return nil
	}
	if err != nil {
		return diag.FromErr(err)
	}

	return resourceDNSRecordSetResourceData(resp, d, site)
}

func resourceDNSRecordUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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

func resourceDNSRecordDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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
