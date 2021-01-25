output "eip_regions" {
  description = "Avaliable regions of EIP."
  value = data.didiyun_eip_regions.eip_regions.regions
}
