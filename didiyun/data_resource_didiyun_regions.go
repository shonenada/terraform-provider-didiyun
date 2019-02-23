package didiyun

import (
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform/helper/schema"
)

func dataSourceDidiyunRegions() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceDidiyunRegionsRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}

}

func dataSourceDidiyunRegionsRead(d *schema.ResourceData, meta interface{}) error {
	accessToken := d.Get("access_token")
	client := &http.Client{}
	req, err := http.NewRequest("POST", LIST_REGIONS_URL, nil)
	if err != nil {
		return err
	}

	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", accessToken))
	resp, err := client.Do(req)
	fmt.Errorf("resp: %s", resp)
	return nil
}
