package didiyun

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	ddy "github.com/shonenada/didiyun-go"
	dc2 "github.com/shonenada/didiyun-go/dc2"
	ddyds "github.com/shonenada/didiyun-go/schema"
)

func dataSourceDidiyunDc2Image() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceDidiyunDc2ImageRead,
		Schema: map[string]*schema.Schema{
			"region_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"filter": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"img_type": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"os_arch": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"os_family": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"os_version": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"platform": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"scene": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},

			// Computed values
			"uuid": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"os_arch": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"os_family": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"os_version": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"platform": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"img_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceDidiyunDc2ImageRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	client := meta.(*ddy.Client)

	var imgType string = ""
	var osArch string = ""
	var osFamily string = ""
	var osVersion string = ""
	var platform string = ""
	var scene string = ""

	if v, ok := d.GetOk("filter"); ok {
		fs := v.(*schema.Set).List()
		for _, rf := range fs {
			f := rf.(map[string]interface{})
			if f["img_type"] != nil && len(f["img_type"].(string)) > 0 {
				imgType = f["img_type"].(string)
			}
			if f["os_arch"] != nil && len(f["os_arch"].(string)) > 0 {
				osArch = f["os_arch"].(string)
			}
			if f["os_family"] != nil && len(f["os_family"].(string)) > 0 {
				osFamily = f["os_family"].(string)
			}
			if f["os_version"] != nil && len(f["os_version"].(string)) > 0 {
				osVersion = f["os_version"].(string)
			}
			if f["platform"] != nil && len(f["platform"].(string)) > 0 {
				platform = f["platform"].(string)
			}
			if f["scene"] != nil && len(f["scene"].(string)) > 0 {
				scene = f["scene"].(string)
			}
		}
	}

	regionId := d.Get("region_id").(string)

	data, err := client.Dc2().ListImage(&dc2.ListImageRequest{
		RegionId: regionId,
	})

	if err != nil {
		return diag.FromErr(err)
	}

	rv := make([]interface{}, 0)

	for _, e := range *data {
		if len(imgType) > 0 && e.ImgType != imgType {
			continue
		}
		if len(osArch) > 0 && e.OsArch != osArch {
			continue
		}
		if len(osFamily) > 0 && e.OsFamily != osFamily {
			continue
		}
		if len(osVersion) > 0 && e.OsVersion != osVersion {
			continue
		}
		if len(platform) > 0 && e.Platform != platform {
			continue
		}
		if len(scene) > 0 {
			found := false
			for _, s := range e.Scenes {
				if s == scene {
					found = true
					break
				}
			}
			if !found {
				continue
			}
		}
		rv = append(rv, e)
	}

	if len(rv) > 1 {
		return diag.Errorf("Filterd images more than 1")
	}

	img := rv[0].(ddyds.ImageInfo)

	d.Set("uuid", img.Uuid)
	d.Set("name", img.Name)
	d.Set("description", img.Description)
	d.Set("os_arch", img.OsArch)
	d.Set("os_family", img.OsFamily)
	d.Set("os_version", img.OsVersion)
	d.Set("platform", img.Platform)
	d.Set("img_type", img.ImgType)

	return diags

}
