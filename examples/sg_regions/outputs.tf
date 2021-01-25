output "sg_regions" {
  description = "Avaliable regions of EBS."
  value = data.didiyun_sg_regions.sg_regions.regions
}
