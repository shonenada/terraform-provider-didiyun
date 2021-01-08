provider "didiyun" {}

data "didiyun_dc2_images" "dc2_images" {
  region_id = var.region_id

  filter {
    img_type = "standard"
    os_family = "Ubuntu"
    os_version = "18.04"
    platform = "Linux"
    scene = "base"
  }
}

