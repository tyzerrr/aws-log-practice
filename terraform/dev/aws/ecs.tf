resource "aws_ecs_cluster" "cluster" {
  name = "${local.project}-${local.env}-ecs-cluster"

  setting {
    name  = "containerInsights"
    value = "enabled"
  }
}

resource "aws_ecs_cluster_capacity_providers" "cluster_capacity_provider" {
  cluster_name = aws_ecs_cluster.cluster.name
  capacity_providers = [
    "FARGATE"
  ]
}
