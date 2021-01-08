output "ebs_regions" {
  description = "Avaliable regions of EBS."
  value = data.didiyun_ebs_regions.ebs_regions.regions
}
