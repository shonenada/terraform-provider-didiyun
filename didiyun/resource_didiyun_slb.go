package didiyun

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	ddy "github.com/shonenada/didiyun-go"
	"github.com/shonenada/didiyun-go/slb"
)

func resourceDidiyunSlb() *schema.Resource {
	return &schema.Resource{
		ReadContext:   resourceDidiyunSlbRead,
		CreateContext: resourceDidiyunSlbCreate,
		UpdateContext: resourceDidiyunSlbUpdate,
		DeleteContext: resourceDidiyunSlbDelete,
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
			"zone_id": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.NoZeroValues,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"is_auto_continue": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"pay_period": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  false,
			},
			"vpc_uuid": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"address_type": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func resourceDidiyunSlbRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	client := meta.(*ddy.Client)

	uuid := d.Id()
	regionId := d.Get("region_id").(string)

	req := slb.GetRequest{
		RegionId: regionId,
		Uuid:     uuid,
	}

	data, err := client.Slb().Get(&req)
	if err != nil {
		return diag.Errorf("Failed to read SLB: %v", err)
	}

	d.Set("name", data.Name)

	return diags
}

func resourceDidiyunSlbCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*ddy.Client)
	req := slb.CreateRequest{
		RegionId:       d.Get("region_id").(string),
		ZoneId:         d.Get("zone_id").(string),
		Name:           d.Get("name").(string),
		IsAutoContinue: d.Get("is_auto_continue").(bool),
		PayPeriod:      d.Get("pay_period").(int),
		VpcUuid:        d.Get("vpc_uuid").(string),
		AddressType:    d.Get("address_type").(string),
	}
	job, err := client.Slb().Create(&req)

	if err != nil {
		return diag.Errorf("Failed to create Slb: %v", err)
	}

	if err := WaitForJob(client, d.Get("region_id").(string), job.Uuid); err != nil {
		return diag.Errorf("Failed to create Ebs: %v", err)
	}

	d.SetId(job.ResourceUuid)

	return resourceDidiyunSlbRead(ctx, d, meta)
}

func resourceDidiyunSlbUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*ddy.Client)

	id := d.Id()
	regionId := d.Get("region_id").(string)

	if d.HasChange("name") {
		name := d.Get("name").(string)
		req := slb.ChangeNameRequest{
			Slb: []slb.ChangeNameParams{
				{
					Uuid: id,
					Name: name,
				},
			},
		}

		job, err := client.slb().ChangeName(&req)

		if err != nil {
			return diag.Errorf("Failed update name of SLB: %v", err)
		}

		if err := WaitForJob(client, regionId, job.Uuid); err != nil {
			return diag.Errorf("Failed update name of SLB: %v", id)
		}
	}

	return resourceDidiyunSlbRead(ctx, d, meta)
}

func resourceDidiyunSlbDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	client := meta.(*ddy.Client)
	req := eip.DeleteRequest{
		Slb: []slb.DeleteParam{
			{
				Uuid: d.Id(),
			},
		},
	}

	job, err := client.Slb().Delete(&req)

	if err != nil {
		return diag.Errorf("Failed to delete EBS: %v", err)
	}

	if err := WaitForJob(client, d.Get("region_id").(string), job.Uuid); err != nil {
		return diag.Errorf("Failed to delete Slb: %v", err)
	}

	d.SetId("")

	return resourceDidiyunSlbRead(ctx, d, meta)
}
