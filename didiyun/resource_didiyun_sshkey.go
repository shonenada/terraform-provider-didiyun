package didiyun

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	ddy "github.com/shonenada/didiyun-go"
	sshkey "github.com/shonenada/didiyun-go/sshkey"
)

func resourceDidiyunSSHKey() *schema.Resource {
	return &schema.Resource{
		ReadContext:   resourceDidiyunSSHKeyRead,
		CreateContext: resourceDidiyunSSHKeyCreate,
		UpdateContext: resourceDidiyunSSHKeyUpdate,
		DeleteContext: resourceDidiyunSSHKeyDelete,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				Description:  "name of this SSH Key",
				ValidateFunc: validation.NoZeroValues,
			},
			"key": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.NoZeroValues,
			},
			"fingerprint": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

// Didiyun does not supported get sshkey by id,
// so we need list all keys, then find the key by id.
func resourceDidiyunSSHKeyRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	client := meta.(*ddy.Client)
	keys, err := client.SSHKey().List()
	if err != nil {
		d.SetId("")
		return diag.Errorf("Failed to request SSH Keys: %v", err)
	}

	for _, ele := range *keys {
		if ele.PubKeyUuid == d.Id() {
			d.Set("name", ele.Name)
			d.Set("key", ele.Key)
			d.Set("fingerprint", ele.Fingerprint)
			return diags
		}
	}

	d.SetId("")
	return diag.Errorf("Failed to find SSH Keys: %v", err)
}

func resourceDidiyunSSHKeyCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*ddy.Client)

	req := sshkey.CreateRequest{
		Name: d.Get("name").(string),
		Key:  d.Get("key").(string),
	}

	data, err := client.SSHKey().Create(&req)
	if err != nil {
		return diag.Errorf("Failed to create SSH Key: %v", err)
	}

	d.SetId(data.PubKeyUuid)

	return resourceDidiyunSSHKeyRead(ctx, d, meta)
}

func resourceDidiyunSSHKeyUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	return diags
}

func resourceDidiyunSSHKeyDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	client := meta.(*ddy.Client)

	req := sshkey.DeleteRequest{
		PubKeyUuid: d.Id(),
	}

	_, err := client.SSHKey().Delete(&req)

	if err != nil {
		return diag.Errorf("Failed to delete SSH Key: %v", err)
	}

	d.SetId("")
	return diags
}
