package didiyun

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	ddy "github.com/shonenada/didiyun-go"
	eip "github.com/shonenada/didiyun-go/eip"
)

func resourceDidiyunEip() *schema.Resource {
	return &schema.Resource{
		ReadContext:   resourceDidiyunEipRead,
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
			"pay_period": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"bandwidth": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"charge_with_flow": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"auto_continue": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
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

func resourceDidiyunEipRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	client := meta.(*ddy.Client)

	uuid := d.Id()
	regionId := d.Get("region").(string)

	data, err := client.Eip().Get(&eip.GetRequest{
		RegionId: regionId,
		Uuid:     uuid,
	})

	if err != nil {
		return diag.Errorf("Failed to read Eip: %v", err)
	}

	d.SetId(data.Uuid)
	d.Set("ip", data.Ip)
	d.Set("state", data.State)
	d.Set("status", data.Status)

	if err := d.Set("tags", FlattenDidiyunTags(data.Tags)); err != nil {
		return diag.Errorf("Failed to set `tags`: %v", err)
	}

	return diags
}

func resourceDidiyunEipCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*ddy.Client)

	var tags []string
	if v, ok := d.GetOk("tags"); ok {
		for _, t := range v.(*schema.Set).List() {
			tags = append(tags, t.(string))
		}
	}

	req := eip.CreateRequest{
		RegionId:         d.Get("region_id").(string),
		IsAutoContinue:   d.Get("auto_continue").(bool),
		IsChargeWithFlow: d.Get("charge_with_flow").(bool),
		PayPeriod:        d.Get("pay_period").(int),
		Count:            d.Get("count").(int),
		BandWidth:        d.Get("bandwidth").(int),
		Tags:             tags,
	}

	job, err := client.Eip().Create(&req)

	if err != nil {
		return diag.Errorf("Failed to create EBS: %v", err)
	}

	if err := WaitForJob(client, d.Get("region_id").(string), job.Uuid); err != nil {
		return diag.Errorf("Failed to create Eip: %v", err)
	}

	d.SetId(job.ResourceUuid)

	return resourceDidiyunEipRead(ctx, d, meta)
}

func resourceDidiyunEipUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	// client := meta.(*ddy.Client)
	return diags
}

func resourceDidiyunEipDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	// client := meta.(*ddy.Client)
	return diags
}
