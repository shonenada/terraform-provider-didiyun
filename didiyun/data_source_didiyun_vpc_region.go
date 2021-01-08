package didiyun

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	ddy "github.com/shonenada/didiyun-go"
)

func dataSourceDidiyunVpcRegion() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceDidiyunVpcRegionRead,
		Schema: map[string]*schema.Schema{
			"regions": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"area_name": {
							Type:         schema.TypeString,
							Computed:     true,
							ValidateFunc: validation.NoZeroValues,
						},
						"id": {
							Type:         schema.TypeString,
							Computed:     true,
							ValidateFunc: validation.NoZeroValues,
						},
						"name": {
							Type:         schema.TypeString,
							Computed:     true,
							ValidateFunc: validation.NoZeroValues,
						},
						"zone": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"id": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"name": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func dataSourceDidiyunVpcRegionRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	client := meta.(*ddy.Client)

	data, err := client.Region().ListVpcRegions()
	if err != nil {
		return diag.FromErr(err)
	}

	var regions []map[string]interface{}
	for _, r := range *data {
		e := make(map[string]interface{})
		e["area_name"] = r.AreaName
		e["id"] = r.Id
		e["name"] = r.Name
		var zones []map[string]interface{}
		for _, ez := range r.Zone {
			z := make(map[string]interface{})
			z["id"] = ez.Id
			z["name"] = ez.Name
			zones = append(zones, z)
		}
		e["zome"] = zones
		regions = append(regions, e)
	}

	d.Set("regions", regions)

	return diags
}
