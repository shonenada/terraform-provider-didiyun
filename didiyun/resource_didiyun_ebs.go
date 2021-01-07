package didiyun

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	ddy "github.com/shonenada/didiyun-go"
	ebs "github.com/shonenada/didiyun-go/ebs"
)

func resourceDidiyunEBS() *schema.Resource {
	return &schema.Resource{
		ReadContext:   resourceDidiyunEbsRead,
		CreateContext: resourceDidiyunEbsCreate,
		UpdateContext: resourceDidiyunEbsUpdate,
		DeleteContext: resourceDidiyunEbsDelete,
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
				Type: schema.TypeString,
			},
			"auto_continue": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"pay_period": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"count": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"coupon_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"size": {
				Type: schema.TypeInt,
			},
			"disk_type": {
				Type: schema.TypeString,
			},
			"dc2_uuid": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"snap_uuid": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"tags": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: validation.NoZeroValues,
				},
			},
		},
	}
}

func resourceDidiyunEbsRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	client := meta.(*ddy.Client)

	uuid := d.Id()
	regionId := d.Get("region_id").(string)

	req := ebs.GetRequest{
		RegionId: regionId,
		EbsUuid:  uuid,
	}

	data, err := client.Ebs().Get(&req)

	if err != nil {
		return diag.Errorf("Failed to read EBS: %v", err)
	}

	d.Set("name", data.Name)
	d.Set("size", data.Size)
	d.Set("disk_type", data.Spec.DiskType)
	d.Set("region_id", data.Region.Id)
	d.Set("zone_id", data.Region.Zone[0].Id)
	d.Set("dc2_uuid", data.Dc2.Uuid)

	if err := d.Set("tags", FlattenDidiyunTags(data.EbsTags)); err != nil {
		return diag.Errorf("Failed to set `tags`: %v", err)
	}

	return diags
}

func resourceDidiyunEbsCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*ddy.Client)

	var tags []string
	if v, ok := d.GetOk("tags"); ok {
		for _, t := range v.(*schema.Set).List() {
			tags = append(tags, t.(string))
		}
	}

	req := ebs.CreateRequest{
		RegionId:     d.Get("region_id").(string),
		ZoneId:       d.Get("zone_id").(string),
		Name:         d.Get("name").(string),
		Size:         d.Get("size").(int64),
		DiskType:     d.Get("disk_type").(string),
		AutoContinue: d.Get("auto_continue").(bool),
		PayPeriod:    d.Get("pay_period").(int),
		Count:        d.Get("dc2_count").(int),
		CouponId:     d.Get("coupon_Id").(string),
		Dc2Uuid:      d.Get("dc2_uuid").(string),
		SnapUuid:     d.Get("snap_uuid").(string),
		EbsTags:      tags,
	}

	job, err := client.Ebs().Create(&req)

	if err != nil {
		return diag.Errorf("Failed to create EBS: %v", err)
	}

	if err := WaitForJob(client, d.Get("region_id").(string), job.Uuid); err != nil {
		return diag.Errorf("Failed to create DC2: %v", err)
	}

	d.SetId(job.ResourceUuid)

	return resourceDidiyunEbsRead(ctx, d, meta)
}

func resourceDidiyunEbsUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {

	client := meta.(*ddy.Client)

	id := d.Id()
	regionId := d.Get("region_id").(string)

	if d.HasChange("name") {
		name := d.Get("name").(string)
		req := ebs.ChangeNameRequest{
			RegionId: regionId,
			Ebs: []ebs.ChangeNameInput{
				{
					EbsUuid: id,
					Name:    name,
				},
			},
		}

		job, err := client.Ebs().ChangeName(&req)
		if err != nil {
			return diag.Errorf("Failed update name of EBS: %v", id)
		}
		if err := WaitForJob(client, regionId, job.Uuid); err != nil {
			return diag.Errorf("Failed update name of EBS: %v", id)
		}
	}

	if d.HasChange("size") {
		size := d.Get("size").(int64)
		req := ebs.ChangeSizeRequest{
			RegionId: regionId,
			Ebs: []ebs.ChangeSizeInput{
				{
					EbsUuid: id,
					Size:    size,
				},
			},
		}

		job, err := client.Ebs().ChangeSize(&req)
		if err != nil {
			return diag.Errorf("Failed update size of EBS: %v", id)
		}
		if err := WaitForJob(client, regionId, job.Uuid); err != nil {
			return diag.Errorf("Failed update size of EBS: %v", id)
		}
	}

	if d.HasChange("dc2_uuid") {
		dc2Uuid := d.Get("dc2_uuid").(string)
		req := ebs.AttachRequest{
			RegionId: regionId,
			Ebs: []ebs.AttachInput{
				{
					EbsUuid: d.Id(),
					Dc2Uuid: dc2Uuid,
				},
			},
		}

		job, err := client.Ebs().Attach(&req)
		if err != nil {
			return diag.Errorf("Failed attach EBS: %v", id)
		}
		if err := WaitForJob(client, regionId, job.Uuid); err != nil {
			return diag.Errorf("Failed attach EBS: %v", id)
		}
	}

	return resourceDidiyunEbsRead(ctx, d, meta)
}

func resourceDidiyunEbsDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	client := meta.(*ddy.Client)

	req := ebs.DeleteRequest{
		RegionId: d.Get("region_id").(string),
		Ebs: []ebs.DeleteInput{
			{
				EbsUuid: d.Id(),
			},
		},
	}

	job, err := client.Ebs().Delete(&req)

	if err != nil {
		return diag.Errorf("Failed to delete EBS: %v", err)
	}

	if err := WaitForJob(client, d.Get("region_id").(string), job.Uuid); err != nil {
		return diag.Errorf("Failed to delete EBS: %v", err)
	}

	d.SetId("")

	return diags
}
