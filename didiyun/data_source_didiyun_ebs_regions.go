package didiyun

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	ddy "github.com/shonenada/didiyun-go"
)

func dataSourceDidiyunEbsRegions() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceDidiyunEbsRegionsRead,
		Schema: map[string]*schema.Schema{
			"regions": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"area_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
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

func dataSourceDidiyunEbsRegionsRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	client := meta.(*ddy.Client)

	data, err := client.Region().ListEbsRegions()
	if err != nil {
		return diag.FromErr(err)
	}

	regions := make([]map[string]interface{}, 0, len(*data))
	for _, r := range *data {
		e := make(map[string]interface{})
		e["area_name"] = r.AreaName
		e["id"] = r.Id
		e["name"] = r.Name
		zones := make([]map[string]interface{}, 0, len(r.Zone))
		for _, ez := range r.Zone {
			z := make(map[string]interface{})
			z["id"] = ez.Id
			z["name"] = ez.Name
			zones = append(zones, z)
		}
		e["zone"] = zones
		regions = append(regions, e)
	}

	d.SetId("ebs_regions")
	d.Set("regions", regions)

	return diags
}
