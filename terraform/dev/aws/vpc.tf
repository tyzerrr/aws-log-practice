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

# VPC
resource "aws_vpc" "vpc" {
  cidr_block = var.vpc_cidr_block
  tags = {
    Name = "${local.project}-${local.env}-vpc"
  }
}

# Subnets
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

# Internet Gateway
resource "aws_internet_gateway" "igw" {
  vpc_id = aws_vpc.vpc.id
  tags = {
    Name = "${local.project}-${local.env}-igw"
  }
}

# NAT Gataway
resource "aws_eip" "eips" {
  for_each = local.public_subnets
  domain   = "vpc"

  tags = {
    AvailabilityZone = each.value.availability_zone
  }
}

resource "aws_nat_gateway" "nat_gateways" {
  for_each      = local.public_subnets
  allocation_id = aws_eip.eips[each.key].allocation_id
  subnet_id     = aws_subnet.public_subnets[each.key].id

  tags = {
    AvailabilityZone = each.value.availability_zone
  }
}

# Route tables
resource "aws_route_table" "public_route_tables" {
  for_each = local.public_subnets
  vpc_id   = aws_vpc.vpc.id

  tags = {
    AvailabilityZone = each.key
    Scope            = "public"
  }
}

resource "aws_route" "public_route" {
  for_each               = local.public_subnets
  route_table_id         = aws_route_table.public_route_tables[each.key].id
  destination_cidr_block = "0.0.0.0/0"
  gateway_id             = aws_internet_gateway.igw.id
}

resource "aws_route_table_association" "public_route_table_associations" {
  for_each       = local.public_subnets
  route_table_id = aws_route_table.public_route_tables[each.key].id
  subnet_id      = aws_subnet.public_subnets[each.key].id
}

resource "aws_route_table" "private_route_tables" {
  for_each = local.private_subnets
  vpc_id   = aws_vpc.vpc.id

  tags = {
    AvailabilityZone = each.key
    Scope            = "private"
  }
}

resource "aws_route" "private_route" {
  for_each               = local.private_subnets
  route_table_id         = aws_route_table.private_route_tables[each.key].id
  destination_cidr_block = "0.0.0.0/0"
  nat_gateway_id         = aws_nat_gateway.nat_gateways[each.key].id
}

resource "aws_route_table_association" "private_route_table_associations" {
  for_each       = local.private_subnets
  route_table_id = aws_route_table.private_route_tables[each.key].id
  subnet_id      = aws_subnet.private_subnets[each.key].id
}