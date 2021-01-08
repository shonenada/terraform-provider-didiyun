provider "didiyun" {}

data "didiyun_dc2_images" "dc2_images" {
  region_id = var.region_id
}

