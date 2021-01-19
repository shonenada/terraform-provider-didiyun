package didiyun

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	ddy "github.com/shonenada/didiyun-go"
)

func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"access_token": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("DIDIYUN_ACCESS_TOKEN", nil),
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"didiyun_dc2":    resourceDidiyunDC2(),
			"didiyun_ebs":    resourceDidiyunEBS(),
			"didiyun_sshkey": resourceDidiyunSSHKey(),
			"didiyun_vpc":    resourceDidiyunVPC(),
		},
		DataSourcesMap: map[string]*schema.Resource{
			"didiyun_dc2_regions":  dataSourceDidiyunDc2Regions(),
			"didiyun_ebs_regions":  dataSourceDidiyunEbsRegions(),
			"didiyun_eip_regions":  dataSourceDidiyunEipRegions(),
			"didiyun_sg_regions":   dataSourceDidiyunSgRegions(),
			"didiyun_snap_regions": dataSourceDidiyunSnapRegions(),
			"didiyun_vpc_regions":  dataSourceDidiyunVpcRegions(),

			"didiyun_dc2_image":  dataSourceDidiyunDc2Image(),
			"didiyun_dc2_images": dataSourceDidiyunDc2Images(),

			"didiyun_eips": dataSourceDidiyunEips(),
			"didiyun_eip":  dataSourceDidiyunEip(),

			"didiyun_sshkeys": dataSourceDidiyunSSHKeys(),
		},
		ConfigureContextFunc: providerConfigure,
	}
}

func providerConfigure(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
	accessToken := d.Get("access_token").(string)

	var diags diag.Diagnostics

	if accessToken != "" {
		c := ddy.Client{
			AccessToken: accessToken,
		}
		return &c, diags
	}

	return nil, diag.Errorf("Failed to create didiyun client. Access token cannot be empty.")
}
