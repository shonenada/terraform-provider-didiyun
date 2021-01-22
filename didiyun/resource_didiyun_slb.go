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
}

func resourceDidiyunSlbCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
}

func resourceDidiyunSlbUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
}

func resourceDidiyunSlbDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
}
