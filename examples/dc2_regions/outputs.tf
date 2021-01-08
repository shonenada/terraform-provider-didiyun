output "dc2_regions" {
  description = "Avaliable regions of DC2."
  value = data.didiyun_dc2_regions.dc2_regions.regions
}
