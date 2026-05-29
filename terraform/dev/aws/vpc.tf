data "aws_availability_zones" "availability_zones" {}

locals {
  availability_zones = slice(sort(data.aws_availability_zones.availability_zones.names), 0, var.az_count)
  public_subnets = {
    for i, az in local.availability_zones :
    az => {
      cidr              = cidrsubnet(var.vpc_cidr_block, var.subnet_cidr_allocated_bit, i)
      availability_zone = az
    }
  }

  private_subnets = {
    for i, az in local.availability_zones :
    az => {
      cidr              = cidrsubnet(var.vpc_cidr_block, var.subnet_cidr_allocated_bit, i + var.az_count)
      availability_zone = az
    }
  }
}

resource "aws_vpc" "vpc" {
  cidr_block = var.vpc_cidr_block
  tags = {
    Name = "${local.project}-${local.env}-vpc"
  }
}

resource "aws_subnet" "private_subnets" {
  for_each          = local.private_subnets
  vpc_id            = aws_vpc.vpc.id
  cidr_block        = each.value.cidr
  availability_zone = each.value.availability_zone

  tags = {
    Name             = "${local.project}-${local.env}-${each.value.availability_zone}"
    AvailabilityZone = each.value.availability_zone
    Scope            = "private"
  }
}

resource "aws_subnet" "public_subnets" {
  for_each          = local.public_subnets
  vpc_id            = aws_vpc.vpc.id
  cidr_block        = each.value.cidr
  availability_zone = each.value.availability_zone

  tags = {
    Name             = "${local.project}-${local.env}-${each.value.availability_zone}"
    AvailabilityZone = each.value.availability_zone
    Scope            = "public"
  }
}
