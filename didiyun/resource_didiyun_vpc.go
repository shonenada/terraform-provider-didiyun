package didiyun

import (
	"fmt"
	"time"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
	didi_job "github.com/shonenada/didiyun-go/job"
	ds "github.com/shonenada/didiyun-go/schema"
	vpc "github.com/shonenada/didiyun-go/vpc"
)

func resourceDidiyunVPC() *schema.Resource {
	return &schema.Resource{
		Read:   resourceDidiyunVPCRead,
		Create: resourceDidiyunVPCCreate,
		Update: resourceDidiyunVPCUpdate,
		Delete: resourceDidiyunVPCDelete,
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
				Optional: false,
			},
			"cidr": {
				Type:     schema.TypeString,
				Required: true,
			},
			"subnet": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					"name": {
						Type:     schema.TypString,
						Required: true,
					},
					"cidr": {
						Type:     schema.TypString,
						Required: true,
					},
					"zone_id": {
						Type:     schema.TypString,
						Required: true,
					},
				},
			},
		},
	}
}

func resourceDidiyunVPCRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).Client()

	uuid := d.Id()
	regionId := d.Get("region_id").(string)

	req := vpc.GetRequest{
		RegionId: regionId,
		VpcUuid:  uuid,
	}

	data, err := client.VPC().Get(&req)
	if err != nil {
		return fmt.Error("Failed to read VPC: %v", err)
	}

	d.Set("name", data.Name)
	d.Set("cidr", data.CIDR)
	d.Set("desc", data.Desc)
	d.Set("is_default", data.IsDefault)

	return nil
}

func resourceDidiyunVPCCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).Client()

	uuid := d.Get()
	regionId := d.Get("region_id").(string)

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
		RegionId: d.Get("region_id"),
		Name:     d.Get("name"),
		CIDR:     d.Get("cidr"),
		Subnet:   subnets,
	}

	job, err := client.VPC().Create(&req)

	if err != nil {
		return fmt.Errorf("Failed to create VPC: %v", err)
	}

	if err := WaitForJob(d.Get("region_id").(string), job.Uuid); err != nil {
		return fmt.Errorf("Failed to create VPC: %v", err)
	}

	d.SetId(job.ResourceUuid)

	return resourceDidiyunVPCRead(d, meta)
}

func resourceDidiyunVPCUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).Client()

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
		job, err := client.VPC().ChangeName(&req)
		if err != nil {
			return fmt.Errorf("Failed update name of VPC: %v", err)
		}

		if err := WaitForJob(region_id, job.Uuid); err != nil {
			return fmt.Errorf("Failed update name of VPC: %v", err)
		}
	}
	return resourceDidiyunVPCRead(d, meta)
}

func resourceDidiyunVPCDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).Client()
	req := vpc.DeleteRequest{
		RegionId: d.Get("region_id").(string),
		Vpc: []vpc.DeleteInput{
			{
				vpcUuid: d.Id(),
			},
		},
	}

	job, err := client.VPC().Delete(&req)

	if err != nil {
		return fmt.Errorf("Failed to delete VPC: %v", err)
	}

	if err := WaitForJob(d.Get("region_id").(string), job.Uuid); err != nil {
		return fmt.Errorf("Failed to delete VPC: %v", err)
	}

	d.SetId("")

	return nil
}
