package didiyun

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	ddy "github.com/shonenada/didiyun-go"
	ebs "github.com/shonenada/didiyun-go/ebs"
	ds "github.com/shonenada/didiyun-go/schema"
)

func dataSourceDidiyunEbss() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceDidiyunEbssRead,
		Schema: map[string]*schema.Schema{
			"region_id": {
				Type:     schema.TypeString,
				Required: true,
			},

			"ebss": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"attr": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"size": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"uuid": {
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
							Type:     schema.TypeList,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
					},
				},
			},
		},
	}
}

func flattenDidiyunEbss(ebss []ds.Ebs) []map[string]interface{} {
	result := make([]map[string]interface{}, 0, len(ebss))

	for _, ebs := range ebss {
		r := make(map[string]interface{})
		r["attr"] = ebs.Attr
		r["name"] = ebs.Name
		r["size"] = ebs.Size
		r["uuid"] = ebs.Uuid

		region := map[string]interface{}{}
		region["id"] = ebs.Region.Id
		region["name"] = ebs.Region.Name
		regions := make([]map[string]interface{}, 0, 1)
		regions = append(regions, region)
		r["region"] = regions

		tags := make([]string, 0, len(ebs.Tags))
		for _, e := range ebs.Tags {
			tags = append(tags, e)
		}
		r["tags"] = tags
	}

	return result
}

func dataSourceDidiyunEbssRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	client := meta.(*ddy.Client)
	regionId := d.Get("region_id").(string)
	data, err := client.Ebs().List(&ebs.ListRequest{
		RegionId: regionId,
		Start:    0,
		Limit:    100,
	})
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("ebss")
	d.Set("ebss", flattenDidiyunEbss(*data))

	return diags
}
