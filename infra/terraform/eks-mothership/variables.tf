variable "region" {
  description = "AWS region"
  type        = string
  default     = "us-east-1"
}

variable "cluster_name" {
  description = "Cluster name"
  type        = string
  default     = null
}

variable "lepton_cloud_route53_zone_id" {
  description = "cloud.lepton.ai Route53 zone ID"
  type        = string
  default     = "Z007822916VK7B4DFVMP7"
}

variable "aurora_master_username" {
  description = "Username for the Aurora master DB user. Required unless `snapshot_identifier` or `replication_source_identifier` is provided or unless a `global_cluster_identifier` is provided when the cluster is the secondary cluster of a global database"
  type        = string
  default     = "root"
}
