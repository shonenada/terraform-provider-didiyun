package didiyun

import (
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
)

func Provider() terraform.ResourceProvider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"token": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"didiyun_dc2": resourceDidiyunDC2(),
		},
		DataSourcesMap: map[string]*schema.Resource{
			"didiyun_regions": dataSourceDidiyunRegions(),
		},
	}
}
