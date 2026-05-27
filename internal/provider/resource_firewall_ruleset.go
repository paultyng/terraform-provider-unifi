package provider

import (
	"context"
	"fmt"
	"sort"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/paultyng/go-unifi/unifi"
)

func resourceFirewallRuleset() *schema.Resource {
	return &schema.Resource{
		Description: "`unifi_firewall_ruleset` manages the order of individual firewall rules in a ruleset. You must provide all rule IDs present in the set for this to succeed. There can only be one ruleset resource per site and ruleset combination. Since this resource will be managed on-the-fly, importing it is optional.",

		CreateContext: reorderFirewallRules,
		ReadContext:   resourceFirewallRulesetRead,
		UpdateContext: reorderFirewallRules,
		DeleteContext: resourceFirewallRulesetDelete,
		Importer: &schema.ResourceImporter{
			StateContext: importRuleset,
		},

		Schema: map[string]*schema.Schema{
			"id": {
				Description: "The ID of the firewall ruleset. It is a concatenation of `<name>:<ruleset>`.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"site": {
				Description: "The name of the site this ruleset is associated with.",
				Type:        schema.TypeString,
				Computed:    true,
				Optional:    true,
				ForceNew:    true,
			},
			"ruleset": {
				Description: "The ruleset to manage. This is from the perspective of the security gateway. " +
					"Must be one of `WAN_IN`, `WAN_OUT`, `WAN_LOCAL`, `LAN_IN`, `LAN_OUT`, `LAN_LOCAL`, `GUEST_IN`, " +
					"`GUEST_OUT`, `GUEST_LOCAL`, `WANv6_IN`, `WANv6_OUT`, `WANv6_LOCAL`, `LANv6_IN`, `LANv6_OUT`, " +
					"`LANv6_LOCAL`, `GUESTv6_IN`, `GUESTv6_OUT`, or `GUESTv6_LOCAL`.",
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"WAN_IN", "WAN_OUT", "WAN_LOCAL", "LAN_IN", "LAN_OUT", "LAN_LOCAL", "GUEST_IN", "GUEST_OUT", "GUEST_LOCAL", "WANv6_IN", "WANv6_OUT", "WANv6_LOCAL", "LANv6_IN", "LANv6_OUT", "LANv6_LOCAL", "GUESTv6_IN", "GUESTv6_OUT", "GUESTv6_LOCAL"}, false),
			},
			"before_predefined": {
				Description: "List of unique rule IDs present in this ruleset in order of their designated index that should be applied before predefined rules.",
				Type:        schema.TypeList,
				Optional:    true,
				// ValidateFunc: validation.ListOfUniqueStrings,
				Elem: &schema.Schema{Type: schema.TypeString},
			},
			"after_predefined": {
				Description: "List of unique rule IDs present in this ruleset in order of their designated index that should be applied after predefined rules.",
				Type:        schema.TypeList,
				Optional:    true,
				// ValidateFunc: validation.ListOfUniqueStrings,
				Elem: &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func resourceFirewallRulesetRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*client)
	var site, ruleset string
	var err error

	id := d.Id()
	if id != "" {
		site, ruleset, err = parseRulesetId(id)
		if err != nil {
			return diag.FromErr(err)
		}
	} else {
		site = d.Get("site").(string)
		if site == "" {
			site = c.site
		}

		ruleset = d.Get("ruleset").(string)
	}

	currentRuleset, err := getRuleset(ctx, c, site, ruleset)
	if err != nil {
		return diag.FromErr(err)
	}

	// Only one ruleset resource can be created per site and ruleset.
	id = site + ":" + strings.ToUpper(ruleset)
	d.SetId(id)
	d.Set("site", site)
	d.Set("ruleset", ruleset)
	d.Set("before_predefined", currentRuleset.getIds(true, false))
	d.Set("after_predefined", currentRuleset.getIds(false, true))
	return nil
}

func resourceFirewallRulesetDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// Since this is an internal representation, there is nothing to delete
	return nil
}

func reorderFirewallRules(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*client)

	site := d.Get("site").(string)
	if site == "" {
		site = c.site
	}

	ruleset := d.Get("ruleset").(string)
	// Ensure all rules are represented.
	// When a new rule is created but unlisted, this check can accidentally
	// pass because it might run before the creation...
	currentRuleset, err := getRuleset(ctx, c, site, ruleset)
	if err != nil {
		return diag.FromErr(err)
	}
	newBeforeRuleIds, err := listToStringSlice(d.Get("before_predefined").([]interface{}))
	if err != nil {
		return diag.FromErr(err)
	}
	newAfterRuleIds, err := listToStringSlice(d.Get("after_predefined").([]interface{}))
	if err != nil {
		return diag.FromErr(err)
	}
	newRuleIds := append(newBeforeRuleIds, newAfterRuleIds...)
	if !compareStringSlicesWithoutOrder(currentRuleset.getIds(true, true), newRuleIds) {
		return diag.Errorf("The set of existing rule IDs in ruleset %s of site %s does not match the set of managed rule IDs", ruleset, site)
	}
	req, err := currentRuleset.buildUpdateRequest(newBeforeRuleIds, newAfterRuleIds)
	if err != nil {
		return diag.FromErr(err)
	}
	err = c.c.ReorderFirewallRules(ctx, site, ruleset, req)
	if err != nil {
		return diag.FromErr(err)
	}

	return resourceFirewallRulesetRead(ctx, d, meta)
}

func getRuleset(ctx context.Context, c *client, site, ruleset string) (FirewallRuleset, error) {
	rules, err := c.c.ListFirewallRule(ctx, site)
	if err != nil {
		return FirewallRuleset{}, err
	}
	rules = filterRulesByRuleset(rules, ruleset)
	return newRuleset(rules), nil
}

func newRuleset(rules []unifi.FirewallRule) FirewallRuleset {
	preRules := make([]unifi.FirewallRule, 0, len(rules))
	postRules := make([]unifi.FirewallRule, 0, len(rules))
	// Filter rules into separate slices
	for _, rule := range rules {
		if rule.RuleIndex >= 2000 && rule.RuleIndex < 3000 {
			preRules = append(preRules, rule)
		} else if rule.RuleIndex >= 4000 && rule.RuleIndex < 5000 {
			postRules = append(postRules, rule)
		}
	}
	// Sort rules in ascending index order.
	sort.Slice(preRules, func(i, j int) bool {
		return preRules[i].RuleIndex < preRules[j].RuleIndex
	})
	sort.Slice(postRules, func(i, j int) bool {
		return postRules[i].RuleIndex < postRules[j].RuleIndex
	})
	preRuleAt := make([]FirewallRuleAtIndex, len(preRules))
	ruleLookup := make(map[string]*FirewallRuleAtIndex, len(rules))
	postRuleAt := make([]FirewallRuleAtIndex, len(postRules))
	for i, rule := range preRules {
		preRuleAt[i] = FirewallRuleAtIndex{ID: rule.ID, RuleIndex: rule.RuleIndex}
		ruleLookup[rule.ID] = &preRuleAt[i]
	}
	for i, rule := range postRules {
		postRuleAt[i] = FirewallRuleAtIndex{ID: rule.ID, RuleIndex: rule.RuleIndex}
		ruleLookup[rule.ID] = &postRuleAt[i]
	}
	return FirewallRuleset{preRules: preRuleAt, postRules: postRuleAt, ruleLookup: ruleLookup}
}

type FirewallRuleAtIndex struct {
	ID        string
	RuleIndex int
}

type FirewallRuleset struct {
	preRules   []FirewallRuleAtIndex
	postRules  []FirewallRuleAtIndex
	ruleLookup map[string]*FirewallRuleAtIndex
}

func (r *FirewallRuleset) buildUpdateRequest(preRules []string, postRules []string) ([]unifi.FirewallRuleIndexUpdate, error) {
	preRuleUpdates, err := r.buildRuleIndexUpdates(preRules, 2000)
	if err != nil {
		return nil, err
	}
	postRuleUpdates, err := r.buildRuleIndexUpdates(postRules, 4000)
	if err != nil {
		return nil, err
	}
	return append(preRuleUpdates, postRuleUpdates...), nil
}

func (r *FirewallRuleset) buildRuleIndexUpdates(rules []string, baseIndex int) ([]unifi.FirewallRuleIndexUpdate, error) {
	ruleUpdates := make([]unifi.FirewallRuleIndexUpdate, 0, len(rules))
	for i, ruleId := range rules {
		rule, ok := r.ruleLookup[ruleId]
		if !ok {
			return ruleUpdates, fmt.Errorf("Lookup in ruleset failed! Rule ID: %s", ruleId)
		}
		if rule.RuleIndex != baseIndex+i {
			ruleUpdates = append(ruleUpdates, unifi.FirewallRuleIndexUpdate{ID: ruleId, RuleIndex: baseIndex + i})
		}
	}
	return ruleUpdates, nil
}

func (r *FirewallRuleset) getIds(pre, post bool) []string {
	totalLength := 0
	if pre {
		totalLength += len(r.preRules)
	}
	if post {
		totalLength += len(r.postRules)
	}
	ids := make([]string, totalLength)
	cur := 0
	if pre {
		for _, rule := range r.preRules {
			ids[cur] = rule.ID
			cur += 1
		}
	}
	if post {
		for _, rule := range r.postRules {
			ids[cur] = rule.ID
			cur += 1
		}
	}
	return ids
}

func importRuleset(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	site, ruleset, err := parseRulesetId(d.Id())
	if err != nil {
		return nil, err
	}
	d.Set("site", site)
	d.Set("ruleset", ruleset)
	d.SetId(site + ":" + strings.ToUpper(ruleset))
	return []*schema.ResourceData{d}, nil
}

func parseRulesetId(id string) (string, string, error) {
	parts := strings.SplitN(id, ":", 2)
	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		return "", "", fmt.Errorf("unexpected format of ID (%s), expected attribute1:attribute2", id)
	}
	// rulesets in IDs have an uppercased v6
	return parts[0], strings.Replace(strings.ToUpper(parts[1]), "V", "v", 1), nil
}
