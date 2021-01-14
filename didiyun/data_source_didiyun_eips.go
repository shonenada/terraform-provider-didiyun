package didiyun

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	ddy "github.com/shonenada/didiyun-go"
	eip "github.com/shonenada/didiyun-go/eip"
	ds "github.com/shonenada/didiyun-go/schema"
)

func dataSourceDidiyunEips() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceDidiyunEipsRead,
		Schema: map[string]*schema.Schema{
			"region_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"eips": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"uuid": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"ip": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"state": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"region": {
							Type:     schema.TypeSet,
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
						"tags": {
							Type:     schema.TypeSet,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
					},
				},
			},
		},
	}
}

func flattenDidiyunEips(eips []ds.Eip) []map[string]interface{} {
	result := make([]map[string]interface{}, 0, len(eips))

	for _, eip := range eips {
		r := make(map[string]interface{})
		r["uuid"] = eip.Uuid
		r["id"] = eip.Id
		r["ip"] = eip.Ip
		r["state"] = eip.State
		r["status"] = eip.Status

		region := map[string]interface{}{}
		region["id"] = eip.Region.Id
		region["name"] = eip.Region.Name

		regions := make([]map[string]interface{}, 0, 1)
		regions = append(regions, region)
		r["region"] = regions

		tags := make([]string, 0, len(eip.Tags))
		for _, e := range eip.Tags {
			tags = append(tags, e)
		}
		r["tags"] = tags

		result = append(result, r)
	}

	return result
}

func dataSourceDidiyunEipsRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	client := meta.(*ddy.Client)

	regionId := d.Get("region_id").(string)

	data, err := client.Eip().List(&eip.ListRequest{
		RegionId:   regionId,
		Start:      0,
		Limit:      100,
		IsSimplify: false,
	})
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("eips")
	d.Set("eips", flattenDidiyunEips(*data))

	return diags
}
