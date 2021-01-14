package didiyun

import (
	"context"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	ddy "github.com/shonenada/didiyun-go"
	dc2 "github.com/shonenada/didiyun-go/dc2"
	ds "github.com/shonenada/didiyun-go/schema"
)

func flattenDidiyunEip(eip ds.Dc2Eip) []map[string]interface{} {
	result := []map[string]interface{}{
		{
			"ip_address": eip.Ip,
			"uuid":       eip.Uuid,
		},
	}
	return result
}

func flattenDidiyunEbs(ebs []ds.Dc2Ebs) []map[string]interface{} {
	result := make([]map[string]interface{}, 0, len(ebs))
	for _, v := range ebs {
		r := make(map[string]interface{})
		r["attr"] = v.Attr
		r["name"] = v.Name
		result = append(result, r)
	}
	return result
}

func expandDc2Eip(m *schema.Set) dc2.EipParams {
	l := m.List()
	d := l[0].(map[string]interface{})

	eip := dc2.EipParams{}

	if v, ok := d["band_width"]; ok {
		eip.BandWidth, _ = v.(int)
	}

	if v, ok := d["charge_with_flow"]; ok {
		eip.IsChargeWithFlow = v.(bool)
	}

	if v, ok := d["tags"]; ok {
		rawTags := v.(*schema.Set).List()
		tags := make([]string, 0)
		for _, e := range rawTags {
			tags = append(tags, e.(string))
		}
		eip.Tags = tags
	}

	return eip
}

func expandDc2Ebs(m *schema.Set) []dc2.EbsParams {
	l := m.List()
	rv := make([]dc2.EbsParams, 0)
	for _, raw := range l {
		d := raw.(map[string]interface{})
		e := dc2.EbsParams{}
		if v, ok := d["count"]; ok {
			e.Count, _ = strconv.Atoi(v.(string))
		}
		if v, ok := d["name"]; ok {
			e.Name = v.(string)
		}
		if v, ok := d["size"]; ok {
			e.Size = v.(int64)
		}
		if v, ok := d["disktype"]; ok {
			e.DiskType = v.(string)
		}
		if v, ok := d["snap_uuid"]; ok {
			e.SnapUuid = v.(string)
		}
		if v, ok := d["tags"]; ok {
			e.Tags = v.([]string)
		}
		rv = append(rv, e)
	}
	return rv
}

func resourceDidiyunDC2() *schema.Resource {
	return &schema.Resource{
		ReadContext:   resourceDidiyunDC2Read,
		CreateContext: resourceDidiyunDC2Create,
		UpdateContext: resourceDidiyunDC2Update,
		DeleteContext: resourceDidiyunDC2Delete,
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
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"name": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.NoZeroValues,
			},
			"auto_continue": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"pay_period": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  0,
			},
			"dc2_count": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  1,
			},
			"subnet_uuid": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"model": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Model of DC2, all avaliable models: https://docs.didiyun.com/static/docs-content/products/DC2/%E5%88%9B%E5%BB%BADC2%EF%BC%88CreateDC2%EF%BC%89.html#Dc2Models",
			},
			"image_uuid": {
				Type:     schema.TypeString,
				ForceNew: true,
				Optional: true,
			},
			"snap_uuid": {
				Type:     schema.TypeString,
				ForceNew: true,
				Optional: true,
			},
			"sshkeys": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: validation.NoZeroValues,
				},
				Description: "List of uuids of SSHKey",
			},
			"password": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"rootdisk_type": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"rootdisk_size": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: validation.IntBetween(40, 500),
			},
			"tags": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: validation.NoZeroValues,
				},
			},
			"sg_uuids": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: validation.NoZeroValues,
				},
			},
			"ip": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Intranet IP",
			},
			"os_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"eip": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"band_width": {
							Type:     schema.TypeInt,
							Optional: true,
						},
						"charge_with_flow": {
							Type:     schema.TypeBool,
							Optional: true,
							Default:  false,
						},
						"tags": {
							Type:     schema.TypeSet,
							Optional: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"uuid": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"ip_address": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"ebs": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"ebs_count": {
							Type:     schema.TypeInt,
							Optional: true,
							Default:  1,
						},
						"attr": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"name": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"size": {
							Type:         schema.TypeInt,
							Optional:     true,
							ValidateFunc: validation.IntBetween(20, 16384),
						},
						"snap_uuid": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"tags": {
							Type:     schema.TypeSet,
							Optional: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
					},
				},
			},
		},
	}
}

func resourceDidiyunDC2Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	client := meta.(*ddy.Client)

	uuid := d.Id()
	regionId := d.Get("region_id").(string)

	req := dc2.GetRequest{
		RegionId: regionId,
		Uuid:     uuid,
	}

	data, err := client.Dc2().Get(&req)
	if err != nil {
		return diag.Errorf("Failed to read Dc2: %v", err)
	}

	d.Set("name", data.Name)
	d.Set("ip", data.Ip)
	d.Set("os_type", data.OSType)
	d.Set("region_id", data.Region.Id)
	d.Set("zone_id", data.Region.Zone.Id)
	d.Set("eip", flattenDidiyunEip(data.Eip))

	if err := d.Set("ebs", flattenDidiyunEbs(data.Ebs)); err != nil {
		return diag.Errorf("[DEBUG] Error setting Dc2 Ebs - error: %#v", err)
	}

	if err := d.Set("tags", FlattenDidiyunTags(data.Tags)); err != nil {
		return diag.Errorf("Failed to set `tags`: %v", err)
	}

	return diags
}

func resourceDidiyunDC2Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*ddy.Client)

	var sshkeys []string

	if v, ok := d.GetOk("sshkeys"); ok {
		for _, k := range v.(*schema.Set).List() {
			sshkeys = append(sshkeys, k.(string))
		}
	}

	var dc2Tags []string
	if v, ok := d.GetOk("tags"); ok {
		for _, t := range v.(*schema.Set).List() {
			dc2Tags = append(dc2Tags, t.(string))
		}
	}

	var sgUuids []string
	if v, ok := d.GetOk("sg_uuids"); ok {
		for _, id := range v.(*schema.Set).List() {
			sgUuids = append(sgUuids, id.(string))
		}
	}

	eip := expandDc2Eip(d.Get("eip").(*schema.Set))
	ebs := expandDc2Ebs(d.Get("ebs").(*schema.Set))

	req := dc2.CreateRequest{
		RegionId:       d.Get("region_id").(string),
		ZoneId:         d.Get("zone_id").(string),
		Name:           d.Get("name").(string),
		IsAutoContinue: d.Get("auto_continue").(bool),
		PayPeriod:      d.Get("pay_period").(int),
		Count:          d.Get("dc2_count").(int),
		SubnetUuid:     d.Get("subnet_uuid").(string),
		Model:          d.Get("model").(string),
		ImgUuid:        d.Get("image_uuid").(string),
		PubKeyUuids:    sshkeys,
		Password:       d.Get("password").(string),
		RootDiskType:   d.Get("rootdisk_type").(string),
		RootDiskSize:   d.Get("rootdisk_size").(int),
		Tags:           dc2Tags,
		SgUuids:        sgUuids,
		Eip:            eip,
		Ebs:            ebs,
	}

	job, err := client.Dc2().Create(&req)

	if err != nil {
		return diag.Errorf("Failed to create DC2: %v", err)
	}

	if err := WaitForJob(client, d.Get("region_id").(string), job.Uuid); err != nil {
		return diag.Errorf("Failed to create DC2: %v", err)
	}

	d.SetId(job.ResourceUuid)

	return resourceDidiyunDC2Read(ctx, d, meta)
}

func resourceDidiyunDC2Update(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {

	client := meta.(*ddy.Client)

	id := d.Id()
	regionId := d.Get("region_id").(string)

	if d.HasChange("name") {
		name := d.Get("name").(string)
		req := dc2.ChangeNameRequest{
			RegionId: regionId,
			Dc2: []dc2.ChangeNameParams{
				{
					Uuid: id,
					Name: name,
				},
			},
		}

		job, err := client.Dc2().ChangeName(&req)

		if err != nil {
			return diag.Errorf("Failed update name of Dc2: %v", err)
		}

		if err := WaitForJob(client, regionId, job.Uuid); err != nil {
			return diag.Errorf("Failed update name of Dc2: %v", id)
		}
	}

	if d.HasChange("password") {
		password := d.Get("password").(string)
		req := dc2.ChangePasswordRequest{
			RegionId: regionId,
			Dc2: []dc2.ChangePasswordParams{
				{
					Uuid:     id,
					Password: password,
				},
			},
		}

		job, err := client.Dc2().ChangePassword(&req)

		if err != nil {
			return diag.Errorf("Failed to change password of dc2: %v", id)
		}

		if err := WaitForJob(client, regionId, job.Uuid); err != nil {
			return diag.Errorf("Failed to change password of dc2: %v", id)
		}
	}

	if d.HasChange("model") {
		model := d.Get("model").(string)
		req := dc2.ChangeSpecRequest{
			RegionId: regionId,
			Dc2: []dc2.ChangeSpecParams{
				{
					Uuid:  id,
					Model: model,
				},
			},
		}

		job, err := client.Dc2().ChangeSpec(&req)

		if err != nil {
			return diag.Errorf("Failed to change model of dc2: %v", id)
		}

		if err := WaitForJob(client, regionId, job.Uuid); err != nil {
			return diag.Errorf("Failed to change model of dc2: %v", id)
		}
	}

	return resourceDidiyunDC2Read(ctx, d, meta)
}

func resourceDidiyunDC2Delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	client := meta.(*ddy.Client)

	req := dc2.DeleteRequest{
		RegionId:    d.Get("region_id").(string),
		IsDeleteEip: true,
		IsDeleteEbs: true,
		IsIgnoreSLB: true,
		Dc2: []dc2.DeleteParams{
			dc2.DeleteParams{
				Uuid: d.Id(),
			},
		},
	}

	job, err := client.Dc2().Delete(&req)

	if err != nil {
		return diag.Errorf("Failed to delete DC2: %v", err)
	}

	if err := WaitForJob(client, d.Get("region_id").(string), job.Uuid); err != nil {
		return diag.Errorf("Failed to delete DC2: %v", err)
	}

	d.SetId("")

	return diags
}
