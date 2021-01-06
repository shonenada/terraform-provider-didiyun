package didiyun

import (
	"github.com/hashicorp/terraform/helper/schema"
)

func FlattenDidiyunTags(tags []string) *schema.Set {
	flattentags := schema.NewSet(schema.HashString, []interface{}{})
	for _, v := range tags {
		flattentags.Add(v)
	}
	return flattentags
}
