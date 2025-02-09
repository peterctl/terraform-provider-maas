package maas

import (
	"context"
	"fmt"

	"github.com/canonical/gomaasclient/client"
	"github.com/canonical/gomaasclient/entity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceMaasZone() *schema.Resource {
	return &schema.Resource{
		Description:   "Provides a resource to manage MAAS zones.",
		CreateContext: resourceZoneCreate,
		ReadContext:   resourceZoneRead,
		UpdateContext: resourceZoneUpdate,
		DeleteContext: resourceZoneDelete,
		Importer: &schema.ResourceImporter{
			StateContext: func(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
				client := meta.(*client.Client)
				zone, err := getZone(client, d.Id())
				if err != nil {
					return nil, err
				}
				d.SetId(zone.Name)
				return []*schema.ResourceData{d}, nil
			},
		},

		Schema: map[string]*schema.Schema{
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "A brief description of the new zone.",
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of the new zone.",
			},
		},
	}
}

func resourceZoneRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*client.Client)

	zone, err := client.Zone.Get(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(zone.Name)

	tfstate := map[string]any{
		"name":        zone.Name,
		"description": zone.Description,
	}

	if err := setTerraformState(d, tfstate); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceZoneCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*client.Client)

	params := getZoneParams(d)
	zone, err := client.Zones.Create(params)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(zone.Name)

	return resourceZoneRead(ctx, d, meta)
}

func resourceZoneUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*client.Client)

	params := getZoneParams(d)
	zone, err := client.Zone.Update(d.Id(), params)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(zone.Name)

	return resourceZoneRead(ctx, d, meta)
}

func resourceZoneDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*client.Client)

	if err := client.Zone.Delete(d.Id()); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func getZoneParams(d *schema.ResourceData) *entity.ZoneParams {
	return &entity.ZoneParams{
		Name:        d.Get("name").(string),
		Description: d.Get("description").(string),
	}
}

func findZone(client *client.Client, identifier string) (*entity.Zone, error) {
	zones, err := client.Zones.Get()
	if err != nil {
		return nil, err
	}
	for _, z := range zones {
		if fmt.Sprintf("%v", z.ID) == identifier || z.Name == identifier {
			return &z, nil
		}
	}
	return nil, nil
}

func getZone(client *client.Client, identifier string) (*entity.Zone, error) {
	zone, err := findZone(client, identifier)
	if err != nil {
		return nil, err
	}
	if zone == nil {
		return nil, fmt.Errorf("zone (%s) was not found", identifier)
	}
	return zone, nil
}
