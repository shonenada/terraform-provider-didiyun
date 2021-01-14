package didiyun

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	ddy "github.com/shonenada/didiyun-go"
	eip "github.com/shonenada/didiyun-go/eip"
)

func dataSourceDidiyunEip() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceDidiyunEipRead,
		Schema: map[string]*schema.Schema{
			"region_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"uuid": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"ip": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"id": {
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
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func getByUuid(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	client := meta.(*ddy.Client)
	uuid := d.Get("uuid").(string)
	regionId := d.Get("region_id").(string)
	data, err := client.Eip().Get(&eip.GetRequest{
		RegionId: regionId,
		Uuid:     uuid,
	})
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(data.Uuid)
	d.Set("id", data.Id)
	d.Set("uuid", data.Uuid)
	d.Set("ip", data.Ip)
	d.Set("state", data.State)
	d.Set("status", data.Status)

	region := map[string]interface{}{}
	region["id"] = data.Region.Id
	region["name"] = data.Region.Name

	regions := make([]map[string]interface{}, 0, 1)
	regions = append(regions, region)
	d.Set("region", regions)

	tags := []string{}
	for _, tag := range data.Tags {
		tags = append(tags, tag)
	}
	d.Set("tags", tags)
	return diags
}

func filterByIp(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	client := meta.(*ddy.Client)
	regionId := d.Get("region_id").(string)
	ip := d.Get("ip").(string)
	data, err := client.Eip().List(&eip.ListRequest{
		RegionId:   regionId,
		Start:      0,
		Limit:      100,
		IsSimplify: false,
	})

	if err != nil {
		return diag.FromErr(err)
	}

	for _, eip := range *data {
		if eip.Ip == ip {
			d.SetId(eip.Uuid)
			d.Set("id", eip.Id)
			d.Set("uuid", eip.Uuid)
			d.Set("ip", eip.Ip)
			d.Set("state", eip.State)
			d.Set("status", eip.Status)

			region := map[string]interface{}{}
			region["id"] = eip.Region.Id
			region["name"] = eip.Region.Name

			regions := make([]map[string]interface{}, 0, 1)
			regions = append(regions, region)
			d.Set("region", regions)

			tags := []string{}
			for _, tag := range eip.Tags {
				tags = append(tags, tag)
			}
			d.Set("tags", tags)
			return diags
		}
	}

	return diags

}

func dataSourceDidiyunEipRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	uuid := d.Get("uuid").(string)
	ip := d.Get("ip").(string)
	if len(uuid) > 0 {
		return getByUuid(ctx, d, meta)
	}
	if len(ip) > 0 {
		return filterByIp(ctx, d, meta)
	}
	return diag.Errorf("uuid or ip is required.")
}
