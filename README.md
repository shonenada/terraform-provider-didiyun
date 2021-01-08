Terraform Provider for Didiyun
==============================

Terraform Didiyun provider is a plugin for Terraform that allows to manage
resource of [Didiyun](https://didiyun.com).

Examples
--------

List Regions of DC2

```hlc
terraform {
  required_version = ">= 0.13"
  required_providers {
    didiyun = {
      version = "0.0.1"
      source = "shonenada/didiyun"
    }
  }
}

provider "didiyun" {
  access_token = "[FILL_UP_ACCESS_TOKEN]"
}

data "didiyun_dc2_regions" "dc2_regions" {}

output "dc2_regions" {
  description = "Avaliable regions of DC2."
  value = data.didiyun_dc2_regions.dc2_regions.regions
}
```
