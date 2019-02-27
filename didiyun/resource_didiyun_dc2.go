package didiyun

import (
	"fmt"
	"strconv"
	"time"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
	dc2 "github.com/shonenada/didiyun-go/dc2"
	didi_job "github.com/shonenada/didiyun-go/job"
	ds "github.com/shonenada/didiyun-go/schema"
)

func flattenDidiyunTags(tags []string) *schema.Set {
	flattentags := schema.NewSet(schema.HashString, []interface{}{})
	for _, v := range tags {
		flattentags.Add(v)
	}
	return flattentags
}

func flattenDidiyunEip(eip ds.Eip) map[string]string {
	result := map[string]string{
		"ip_address": eip.Ip,
	}
	return result
}

func flattenDidiyunEbs(ebs []ds.Ebs) []map[string]interface{} {
	result := make([]map[string]interface{}, 0, 1)
	for _, v := range ebs {
		r := make(map[string]interface{})
		r["attr"] = v.Attr
		r["name"] = v.Name
		r["size"] = v.Spec.Size
		r["disktype"] = v.Spec.DiskType
		r["tags"] = v.EbsTags

		result = append(result, r)
	}
	return result
}

func resourceDidiyunDC2() *schema.Resource {
	return &schema.Resource{
		Read:   resourceDidiyunDC2Read,
		Create: resourceDidiyunDC2Create,
		Update: resourceDidiyunDC2Update,
		Delete: resourceDidiyunDC2Delete,
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
			"dc2_model": {
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
				Type:     schema.TypeMap,
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
						"ip_address": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"ebs": {
				Type:     schema.TypeList,
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
							Required:     true,
							ValidateFunc: validation.IntBetween(20, 16384),
						},
						"disktype": {
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
								Type: schema.TypeString,
							},
						},
					},
				},
			},
		},
	}
}

func resourceDidiyunDC2Read(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).Client()

	uuid := d.Id()
	regionId := d.Get("region_id").(string)

	req := dc2.GetRequest{
		RegionId: regionId,
		Dc2Uuid:  uuid,
	}

	data, err := client.Dc2().Get(&req)
	if err != nil {
		return fmt.Errorf("Failed to read Dc2: %v", err)
	}

	d.Set("name", data.Name)
	d.Set("ip", data.Ip)
	d.Set("os_type", data.OSType)
	d.Set("region_id", data.Region.Id)
	d.Set("zone_id", data.Region.Zone.Id)
	d.Set("eip", flattenDidiyunEip(data.Eip))
	d.Set("ebs", flattenDidiyunEbs(data.Ebs))

	if err := d.Set("tags", flattenDidiyunTags(data.Tags)); err != nil {
		return fmt.Errorf("Failed to set `tags`: %v", err)
	}

	return nil
}

func resourceDidiyunDC2Create(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).Client()

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

	var eip dc2.EipInput
	if data, ok := d.GetOk("eip"); ok {
		d := data.(map[string]interface{})
		if v, ok := d["band_width"]; ok {
			eip.BandWidth, _ = strconv.Atoi(v.(string))
		}

		if v, ok := d["charge_with_flow"]; ok {
			eip.ChargeWithFlow = v.(bool)
		}

		if v, ok := d["tags"]; ok {
			eip.EipTags = v.([]string)
		}
	}

	var ebs []dc2.EbsInput
	if data, ok := d.GetOk("ebs"); ok {
		d := data.([]map[string]interface{})
		for _, e := range d {
			t := dc2.EbsInput{}

			if v, ok := e["count"]; ok {
				t.Count, _ = strconv.Atoi(v.(string))
			}

			if v, ok := e["name"]; ok {
				t.Name = v.(string)
			}

			if v, ok := e["size"]; ok {
				t.Size = v.(int64)
			}

			if v, ok := e["disktype"]; ok {
				t.DiskType = v.(string)
			}

			if v, ok := e["snap_uuid"]; ok {
				t.SnapUuid = v.(string)
			}

			if v, ok := e["tags"]; ok {
				t.EbsTags = v.([]string)
			}

			ebs = append(ebs, t)
		}
	}

	req := dc2.CreateRequest{
		RegionId:     d.Get("region_id").(string),
		ZoneId:       d.Get("zone_id").(string),
		Name:         d.Get("name").(string),
		AutoContinue: d.Get("auto_continue").(bool),
		PayPeriod:    d.Get("pay_period").(int),
		Count:        d.Get("dc2_count").(int),
		SubnetUuid:   d.Get("subnet_uuid").(string),
		Dc2Model:     d.Get("dc2_model").(string),
		ImgUuid:      d.Get("image_uuid").(string),
		PubKeyUuids:  sshkeys,
		Password:     d.Get("password").(string),
		RootDiskType: d.Get("rootdisk_type").(string),
		RootDiskSize: d.Get("rootdisk_size").(int),
		Dc2Tags:      dc2Tags,
		SgUuids:      sgUuids,
		Eip:          eip,
		Ebs:          ebs,
	}

	data, err := client.Dc2().Create(&req)

	if err != nil {
		return fmt.Errorf("Failed to create DC2: %v", err)
	}

	err = resource.Retry(2*time.Minute, func() *resource.RetryError {
		jobs, err := client.Job().GetResult(&didi_job.ResultRequest{
			RegionId: d.Get("region_id").(string),
			JobUuids: data.Uuid,
		})
		if err != nil {
			return resource.RetryableError(fmt.Errorf("Failed to get job: %v", err))
		}

		job := (*jobs)[0]

		if job.Progress < 100 {
			return resource.RetryableError(fmt.Errorf("Wait for job"))
		}

		if !job.Done {
			return resource.RetryableError(fmt.Errorf("Wait for job"))
		}

		if !job.Success {
			return resource.NonRetryableError(fmt.Errorf("Failed to execute job: %v", job.Result))
		}

		return nil
	})

	d.SetId(data.ResourceUuid)

	return resourceDidiyunDC2Read(d, meta)
}

func resourceDidiyunDC2Update(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func resourceDidiyunDC2Delete(d *schema.ResourceData, meta interface{}) error {
	return nil
}
