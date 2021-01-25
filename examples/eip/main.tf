provider "didiyun" {}

data "didiyun_eip" "main" {
  region_id = var.region_id
  uuid = var.uuid
}

