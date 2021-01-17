package didiyun

import (
	"context"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	ddy "github.com/shonenada/didiyun-go"
	ds "github.com/shonenada/didiyun-go/schema"
	sg "github.com/shonenada/didiyun-go/sg"
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

		Schema: map[string]*schema.Schema{},
	}
}

func resourceDidiyunSgRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	client := meta.(*ddy.Client)
	return diags
}

func resourceDidiyunSgCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*ddy.Client)
	return resourceDidiyunDC2Read(ctx, d, meta)
}

func resourceDidiyunSgUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*ddy.Client)
	return resourceDidiyunDC2Read(ctx, d, meta)
}

func resourceDidiyunSgDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	client := meta.(*ddy.Client)
	return diags
}
