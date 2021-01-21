package didiyun

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	ddy "github.com/shonenada/didiyun-go"
	"github.com/shonenada/didiyun-go/sg"
)

func resourceDidiyunSg() *schema.Resource {
	return &schema.Resource{
		ReadContext:   resourceDidiyunSgRead,
		CreateContext: resourceDidiyunSgCreate,
		UpdateContext: resourceDidiyunSgUpdate,
		DeleteContext: resourceDidiyunSgDelete,
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
		},
	}
}

func resourceDidiyunSgRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	client := meta.(*ddy.Client)

	uuid := d.Id()
	regionId := d.Get("region_id").(string)

	data, err := client.Sg().GetByUuid(regionId, uuid)
	if err != nil {
		return diag.Errorf("Failed to read SG: %v", err)
	}

	d.Set("name", data.Name)

	return diags
}

func resourceDidiyunSgCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*ddy.Client)

	req := sg.CreateRequest{
		RegionId: d.Get("region_id").(string),
		Name:     d.Get("name").(string),
	}

	job, err := client.Sg().Create(&req)

	if err != nil {
		return diag.Errorf("Failed to create Sg: %v", err)
	}

	d.SetId(job.ResourceUuid)

	return resourceDidiyunDC2Read(ctx, d, meta)
}

func resourceDidiyunSgUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*ddy.Client)
	id := d.Id()
	regionId := d.Get("region_id").(string)

	if d.HasChange("name") {
		name := d.Get("name").(string)
		req := sg.ChangeNameRequest{
			RegionId: d.Get("region_id").(string),
			Sg: []sg.ChangeNameParams{
				{
					Uuid: id,
					Name: name,
				},
			},
		}
		job, err := client.Sg().ChangeName(&req)
		if err != nil {
			return diag.Errorf("Failed change name of Sg: %v", id)
		}
		if err := WaitForJob(client, regionId, job.Uuid); err != nil {
			return diag.Errorf("Failed attach EBS: %v", id)
		}
	}

	return resourceDidiyunDC2Read(ctx, d, meta)
}

func resourceDidiyunSgDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	client := meta.(*ddy.Client)

	req := sg.DeleteRequest{
		RegionId: d.Get("region_id").(string),
		Sg: []sg.DeleteParams{
			{
				Uuid: d.Id(),
			},
		},
	}

	job, err := client.Sg().Delete(&req)

	if err != nil {
		return diag.Errorf("Failed to delete Sg: %v", err)
	}

	if err := WaitForJob(client, d.Get("region_id").(string), job.Uuid); err != nil {
		return diag.Errorf("Failed to delete Sg: %v", err)
	}

	d.SetId("")

	return diags
}
