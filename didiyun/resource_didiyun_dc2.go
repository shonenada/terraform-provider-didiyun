package didiyun

import (
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceDidiyunDC2() *schema.Resource {
	return &schema.Resource{
		Read:   resourceDidiyunDC2Read,
		Create: resourceDidiyunDC2Create,
		Update: resourceDidiyunDC2Update,
		Delete: resourceDidiyunDC2Delete,

		Schema: map[string]*schema.Schema{
			"region_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"zone_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"dc2_uuid": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func resourceDidiyunDC2Read(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func resourceDidiyunDC2Create(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func resourceDidiyunDC2Update(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func resourceDidiyunDC2Delete(d *schema.ResourceData, meta interface{}) error {
	return nil
}
