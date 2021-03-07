package provider

import (
	"context"
	"errors"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/paultyng/go-unifi/unifi"
)

func resourceSite() *schema.Resource {
	return &schema.Resource{
		Description: "`unifi_site` manages Unifi sites",

		Create: resourceSiteCreate,
		Read:   resourceSiteRead,
		Update: resourceSiteUpdate,
		Delete: resourceSiteDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceSiteImport,
		},

		Schema: map[string]*schema.Schema{
			"id": {
				Description: "The ID of the site.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"description": {
				Description: "The description of the site.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"name": {
				Description: "The name of the site.",
				Type:        schema.TypeString,
				Computed:    true,
			},
		},
	}
}

func resourceSiteImport(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	c := meta.(*client)

	id := d.Id()
	_, err := c.c.GetSite(ctx, id)
	if err != nil {
		var nf *unifi.NotFoundError
		if !errors.As(err, &nf) {
			return nil, err
		}
	} else {
		// id is a valid site
		return []*schema.ResourceData{d}, nil
	}

	// lookup site by name
	sites, err := c.c.ListSites(ctx)
	if err != nil {
		return nil, err
	}

	for _, s := range sites {
		if s.Name == id {
			d.SetId(s.ID)
			return []*schema.ResourceData{d}, nil
		}
	}

	return nil, fmt.Errorf("unable to find site %q on controller", id)
}

func resourceSiteCreate(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*client)

	description := d.Get("description").(string)

	resp, err := c.c.CreateSite(context.TODO(), description)
	if err != nil {
		return err
	}

	site := resp[0]
	d.SetId(site.ID)

	return resourceSiteSetResourceData(&site, d)
}

func resourceSiteSetResourceData(resp *unifi.Site, d *schema.ResourceData) error {
	d.Set("name", resp.Name)
	d.Set("description", resp.Description)
	return nil
}

func resourceSiteRead(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*client)

	id := d.Id()

	site, err := c.c.GetSite(context.TODO(), id)
	if _, ok := err.(*unifi.NotFoundError); ok {
		d.SetId("")
		return nil
	}
	if err != nil {
		return err
	}

	return resourceSiteSetResourceData(site, d)
}

func resourceSiteUpdate(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*client)

	site := &unifi.Site{
		ID:          d.Id(),
		Name:        d.Get("name").(string),
		Description: d.Get("description").(string),
	}

	resp, err := c.c.UpdateSite(context.TODO(), site.Name, site.Description)
	if err != nil {
		return err
	}

	return resourceSiteSetResourceData(&resp[0], d)
}

func resourceSiteDelete(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*client)
	id := d.Id()
	_, err := c.c.DeleteSite(context.TODO(), id)
	return err
}
