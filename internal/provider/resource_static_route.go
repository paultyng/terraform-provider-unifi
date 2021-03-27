package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/paultyng/go-unifi/unifi"
)

func resourceStaticRoute() *schema.Resource {
	return &schema.Resource{
		Description: "`unifi_static_route` manages a static route.",

		Create: resourceStaticRouteCreate,
		Read:   resourceStaticRouteRead,
		Update: resourceStaticRouteUpdate,
		Delete: resourceStaticRouteDelete,
		Importer: &schema.ResourceImporter{
			State: importSiteAndID,
		},

		Schema: map[string]*schema.Schema{
			"id": {
				Description: "The ID of the static route.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"site": {
				Description: "The name of the site to associate the static route with.",
				Type:        schema.TypeString,
				Computed:    true,
				Optional:    true,
				ForceNew:    true,
			},
			"name": {
				Description: "The name of the static route.",
				Type:        schema.TypeString,
				Required:    true,
			},

			"network": {
				Description:      "The network subnet address.",
				Type:             schema.TypeString,
				Required:         true,
				ValidateFunc:     cidrValidate,
				DiffSuppressFunc: cidrDiffSuppress,
			},
			"type": {
				Description:  "The type of static route. Can be `interface-route`, `nexthop-route`, or `blackhole`.",
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringInSlice([]string{"interface-route", "nexthop-route", "blackhole"}, false),
			},
			"distance": {
				Description: "The distance of the static route.",
				Type:        schema.TypeInt,
				Required:    true,
			},

			"next_hop": {
				Description:  "The next hop of the static route (only valid for `nexthop-route` type).",
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.IsIPv4Address,
			},
			"interface": {
				Description: "The interface of the static route (only valid for `interface-route` type). This can be `WAN1`, `WAN2`, or a network ID.",
				Type:        schema.TypeString,
				Optional:    true,
			},
		},
	}
}

func resourceStaticRouteCreate(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*client)

	req, err := resourceStaticRouteGetResourceData(d)
	if err != nil {
		return err
	}

	site := d.Get("site").(string)
	if site == "" {
		site = c.site
	}

	resp, err := c.c.CreateRouting(context.TODO(), site, req)
	if err != nil {
		return err
	}

	d.SetId(resp.ID)

	return resourceStaticRouteSetResourceData(resp, d, site)
}

func resourceStaticRouteGetResourceData(d *schema.ResourceData) (*unifi.Routing, error) {
	t := d.Get("type").(string)

	r := &unifi.Routing{
		Enabled: true,
		Type:    "static-route",

		Name:                d.Get("name").(string),
		StaticRouteNetwork:  cidrZeroBased(d.Get("network").(string)),
		StaticRouteDistance: d.Get("distance").(int),
		StaticRouteType:     t,
	}

	switch t {
	case "interface-route":
		r.StaticRouteInterface = d.Get("interface").(string)
	case "nexthop-route":
		r.StaticRouteNexthop = d.Get("next_hop").(string)
	case "blackhole":
	default:
		return nil, fmt.Errorf("unexpected route type: %q", t)
	}

	return r, nil
}

func resourceStaticRouteSetResourceData(resp *unifi.Routing, d *schema.ResourceData, site string) error {
	d.Set("site", site)
	d.Set("name", resp.Name)
	d.Set("network", cidrZeroBased(resp.StaticRouteNetwork))
	d.Set("distance", resp.StaticRouteDistance)

	t := resp.StaticRouteType
	d.Set("type", t)

	d.Set("next_hop", "")
	d.Set("interface", "")

	switch t {
	case "interface-route":
		d.Set("interface", resp.StaticRouteInterface)
	case "nexthop-route":
		d.Set("next_hop", resp.StaticRouteNexthop)
	case "blackhole":
		// no additional attributes
	default:
		return fmt.Errorf("unexpected static route type: %q", t)
	}

	return nil
}

func resourceStaticRouteRead(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*client)

	id := d.Id()

	site := d.Get("site").(string)
	if site == "" {
		site = c.site
	}

	resp, err := c.c.GetRouting(context.TODO(), site, id)
	if _, ok := err.(*unifi.NotFoundError); ok {
		d.SetId("")
		return nil
	}
	if err != nil {
		return err
	}

	return resourceStaticRouteSetResourceData(resp, d, site)
}

func resourceStaticRouteUpdate(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*client)

	req, err := resourceStaticRouteGetResourceData(d)
	if err != nil {
		return err
	}

	req.ID = d.Id()

	site := d.Get("site").(string)
	if site == "" {
		site = c.site
	}
	req.SiteID = site

	resp, err := c.c.UpdateRouting(context.TODO(), site, req)
	if err != nil {
		return err
	}

	return resourceStaticRouteSetResourceData(resp, d, site)
}

func resourceStaticRouteDelete(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*client)

	id := d.Id()

	site := d.Get("site").(string)
	if site == "" {
		site = c.site
	}
	err := c.c.DeleteRouting(context.TODO(), site, id)
	if _, ok := err.(*unifi.NotFoundError); ok {
		return nil
	}
	return err
}
