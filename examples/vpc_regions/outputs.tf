output "vpc_regions" {
  description = "Avaliable regions of EBS."
  value = data.didiyun_vpc_regions.vpc_regions.regions
}
