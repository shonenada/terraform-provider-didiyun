package didiyun

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	ddy "github.com/shonenada/didiyun-go"
)

func dataSourceDidiyunSSHKeys() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceDidiyunSSHKeysRead,
		Schema: map[string]*schema.Schema{
			"sshkeys": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"uuid": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"fingerprint": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceDidiyunSSHKeysRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	client := meta.(*ddy.Client)

	data, err := client.SSHKey().List()
	if err != nil {
		return diag.FromErr(err)
	}

	sshkeys := make([]map[string]interface{}, 0, len(*data))
	for _, r := range *data {
		e := make(map[string]interface{})
		e["uuid"] = r.Uuid
		e["name"] = r.Name
		e["fingerprint"] = r.Fingerprint
		sshkeys = append(sshkeys, e)
	}

	d.SetId("sshkeys")
	d.Set("sshkeys", sshkeys)

	return diags
}
