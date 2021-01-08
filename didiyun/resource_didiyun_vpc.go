package didiyun

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	ddy "github.com/shonenada/didiyun-go"
	vpc "github.com/shonenada/didiyun-go/vpc"
)

func resourceDidiyunVPC() *schema.Resource {
	return &schema.Resource{
		ReadContext:   resourceDidiyunVPCRead,
		CreateContext: resourceDidiyunVPCCreate,
		UpdateContext: resourceDidiyunVPCUpdate,
		DeleteContext: resourceDidiyunVPCDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"region_id": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.NoZeroValues,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"is_default": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"desc": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"cidr": {
				Type:     schema.TypeString,
				Required: true,
			},
			"subnet": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Required: true,
						},
						"cidr": {
							Type:     schema.TypeString,
							Required: true,
						},
						"zone_id": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},
		},
	}
}

func resourceDidiyunVPCRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	client := meta.(*ddy.Client)

	uuid := d.Id()
	regionId := d.Get("region_id").(string)

	req := vpc.GetRequest{
		RegionId: regionId,
		VpcUuid:  uuid,
	}

	data, err := client.Vpc().Get(&req)
	if err != nil {
		return diag.Errorf("Failed to read VPC: %v", err)
	}

	d.Set("name", data.Name)
	d.Set("cidr", data.CIDR)
	d.Set("desc", data.Desc)
	d.Set("is_default", data.IsDefault)

	return diags
}

func resourceDidiyunVPCCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*ddy.Client)

	var subnets []vpc.SubnetInput
	if data, ok := d.GetOk("subnet"); ok {
		d := data.([]interface{})
		for _, each := range d {
			e := each.(map[string]interface{})
			t := vpc.SubnetInput{}

			if v, ok := e["name"]; ok {
				t.Name = v.(string)
			}
			if v, ok := e["cidr"]; ok {
				t.CIDR = v.(string)
			}
			if v, ok := e["zone_id"]; ok {
				t.ZoneId = v.(string)
			}

			subnets = append(subnets, t)
		}
	}

	req := vpc.CreateRequest{
		RegionId: d.Get("region_id").(string),
		Name:     d.Get("name").(string),
		CIDR:     d.Get("cidr").(string),
		Subnet:   subnets,
	}

	job, err := client.Vpc().Create(&req)

	if err != nil {
		return diag.Errorf("Failed to create VPC: %v", err)
	}

	if err := WaitForJob(client, d.Get("region_id").(string), job.Uuid); err != nil {
		return diag.Errorf("Failed to create VPC: %v", err)
	}

	d.SetId(job.ResourceUuid)

	return resourceDidiyunVPCRead(ctx, d, meta)
}

func resourceDidiyunVPCUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*ddy.Client)

	id := d.Id()
	region_id := d.Get("region_id").(string)

	if d.HasChange("name") {
		name := d.Get("name").(string)
		req := vpc.ChangeNameRequest{
			RegionId: region_id,
			Vpc: []vpc.ChangeNameInput{{
				VpcUuid: id,
				Name:    name,
			}},
		}
		job, err := client.Vpc().ChangeName(&req)
		if err != nil {
			return diag.Errorf("Failed update name of VPC: %v", err)
		}

		if err := WaitForJob(client, region_id, job.Uuid); err != nil {
			return diag.Errorf("Failed update name of VPC: %v", err)
		}
	}
	return resourceDidiyunVPCRead(ctx, d, meta)
}

func resourceDidiyunVPCDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	client := meta.(*ddy.Client)

	req := vpc.DeleteRequest{
		RegionId: d.Get("region_id").(string),
		Vpc: []vpc.DeleteInput{
			{
				VpcUuid: d.Id(),
			},
		},
	}

	job, err := client.Vpc().Delete(&req)

	if err != nil {
		return diag.Errorf("Failed to delete VPC: %v", err)
	}

	if err := WaitForJob(client, d.Get("region_id").(string), job.Uuid); err != nil {
		return diag.Errorf("Failed to delete VPC: %v", err)
	}

	d.SetId("")

	return diags
}
