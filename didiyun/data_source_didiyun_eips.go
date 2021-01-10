package didiyun

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	ddy "github.com/shonenada/didiyun-go"
)

func dataSourceDidiyunEips() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceDidiyunEipsRead,
		Schema: map[string]*schema.Schema{
			"eips": {
				Type: schema.TypeSet,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"charge": {
							Type: schem.String,
							Computed: true,
						},
						"id": {
							Type: schem.String,
							Computed: true,
						},
						"uuid": {
							Type: schem.String,
							Computed: true,
						},
						"ip": {
							Type: schem.String,
							Computed: true,
						},
						"state": {
							Type: schem.String,
							Computed: true,
						},
						"status": {
							Type: schem.String,
							Computed: true,
						},
						"spec": {
							Type: schema.TypSet,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"name": {
										Type: schema.TypeString,
										Computed: true,
									},
									"uuid": {
										Type: schema.TypeString,
										Computed: true,
									},
									"type": {
										Type: schema.TypeString,
										Computed: true,
									},
									"state": {
										Type: schema.TypeString,
										Computed: true,
									},
									"bandwidth": {
										Type: schema.TypeInt,
										Computed: true,
									},
									"charge_type": {
										Type: schema.TypeString,
										Computed: true,
									},
									"description": {
										Type: schema.TypeString,
										Computed: true,
									},
									"inbound_bandwidth": {
										Type: schema.TypeInt,
										Computed: true,
									},
									"outbound_bandwidth": {
										Type: schema.TypeInt,
										Computed: true,
									},
									"offering_uuid": {
										Type: schema.TypeString,
										Computed: true,
									},
									"peer_offering_uuid": {
										Type: schema.TypeString,
										Computed: true,
									},
								}
							},
						},
						"region": {
							Type: schema.TypSet,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"id": {
										Type: schema.TypeString,
										Computed: true,
									},
									"area_name": {
										Type: schema.TypeString,
										Computed: true,
									},
									"name": {
										Type: schema.TypeString,
										Computed: true,
									},
								}
							},
						},
						"tags": {
							Type: schema.TypSet,
							Computed: true,
							Elem: &schema.Schema{Type: schema.TypeString},
						},
					},
				},
			},
		},
	}
}

func dataSourceDidiyunEipsRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	client := meta.(*ddy.Client)

	data, err := client.Eip().List()
	if err != nil {
		return diag.FromErr(err)
	}
	return diags
}
