package didiyun

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	ddy "github.com/shonenada/didiyun-go"
	ds "github.com/shonenada/didiyun-go/schema"
	sg "github.com/shonenada/didiyun-go/sg"
)

func dataSourceDidiyunSgs() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceDidiyunSgsRead,
		Schema: map[string]*schema.Schema{
			"region_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"sgs": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"uuid": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"is_default": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"dc2_count": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"sg_rule_count": {
							Type:     schema.TypeInt,
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
						"vpc": {
							Type:     schema.TypeSet,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"name": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"description": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"is_default": {
										Type:     schema.TypeBool,
										Computed: true,
									},
									"cidr": {
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

func flattenDidiyunSgs(sgs []ds.Sg) []map[string]interface{} {
	result := make([]map[string]interface{}, 0, len(sgs))

	for _, sg := range sgs {
		r := make(map[string]interface{})
		r["uuid"] = sg.Uuid
		r["name"] = sg.Name
		r["is_default"] = sg.IsDefault
		r["dc2_count"] = sg.Dc2Count
		r["sg_rule_count"] = sg.SgRuleCount

		region := map[string]interface{}{}
		region["id"] = sg.Region.Id
		region["name"] = sg.Region.Name

		regions := make([]map[string]interface{}, 0, 1)
		regions = append(regions, region)
		r["region"] = regions

		vpc := map[string]interface{}{}
		vpc["name"] = sg.Vpc.Name
		vpc["description"] = sg.Vpc.Description
		vpc["is_default"] = sg.Vpc.IsDefault
		vpc["cidr"] = sg.Vpc.CIDR

		vpcs := make([]map[string]interface{}, 0, 1)
		vpcs = append(vpcs, vpc)
		r["vpc"] = vpcs

		result = append(result, r)
	}

	return result
}

func dataSourceDidiyunSgsRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	client := meta.(*ddy.Client)

	regionId := d.Get("region_id").(string)

	data, err := client.Sg().List(&sg.ListRequest{
		RegionId: regionId,
		Start:    0,
		Limit:    100,
	})
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("sgs")
	d.Set("sgs", flattenDidiyunSgs(*data))

	return diags
}
