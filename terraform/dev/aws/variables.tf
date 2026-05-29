variable "vpc_cidr_block" {
  type        = string
  default     = "10.0.0.0/16"
  description = "VPCに割り当てるCIDR Block"

  validation {
    condition     = can(regex("\\d{1,3}\\.\\d{1,3}\\.\\d{1,3}\\.\\d{1,3}/\\d{1,2}", var.vpc_cidr_block))
    error_message = "Specify vpc CIDR block with the CIDR format"
  }
}

variable "az_count" {
  type        = number
  default     = 2
  description = "Availability zone count"
}

variable "subnet_cidr_allocated_bit" {
  type        = number
  default     = 8
  description = "New allocated bit from VPC"
}

variable "github_organization" {
  type = string
  default = "tyzerrr"
  description = "Github organization name"
}

variable "github_repository" {
  type = string
  default = "aws-log-practice"
  description = "Github repository name"
}