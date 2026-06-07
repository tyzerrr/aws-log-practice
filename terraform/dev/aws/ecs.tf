resource "aws_ecs_cluster" "cluster" {
  name = "${local.project}-${local.env}-ecs-cluster"

  setting {
    name  = "containerInsights"
    value = "enabled"
  }

  tags = {
    Name = "${local.project}-${local.env}-ecs-cluster"
  }
}

resource "aws_ecs_cluster_capacity_providers" "cluster_capacity_provider" {
  cluster_name = aws_ecs_cluster.cluster.name
  capacity_providers = [
    "FARGATE"
  ]
}

resource "aws_ecs_task_definition" "task_definition" {
  family                   = "${local.project}-${local.env}-ecs-task"
  requires_compatibilities = ["FARGATE"]
  network_mode             = "awsvpc"
  cpu                      = 1024
  memory                   = 2048
  task_role_arn            = aws_iam_role.task.arn
  execution_role_arn       = aws_iam_role.task_execution.arn

  # lifecycleでcontainer definition changesを無視しているので
  # ecspressoからコンテナ定義をいじってDeployするように移行してもTerraformコードを変更する必要がない
  container_definitions = templatefile("./ecs/container_definitions.json", {
    log_group = aws_cloudwatch_log_group.ecs.name
    region    = local.aws_region
  })

  lifecycle {
    ignore_changes = [
      container_definitions
    ]
  }

  tags = {
    Name = "${local.project}-${local.env}-ecs-task"
  }
}

resource "aws_ecs_service" "service" {
  name            = "${local.project}-${local.env}-ecs-service"
  cluster         = aws_ecs_cluster.cluster.id
  task_definition = aws_ecs_task_definition.task_definition.arn
  desired_count   = 1
  launch_type     = "FARGATE"

  network_configuration {
    subnets          = [for subnet in aws_subnet.private_subnets : subnet.id]
    security_groups  = [aws_security_group.ecs_task_sg.id]
    assign_public_ip = false
  }

  load_balancer {
    target_group_arn = aws_alb_target_group.target_group.arn
    container_name   = "nginx"
    container_port   = local.http_port
  }

  tags = {
    Name                           = "${local.project}-${local.env}-ecs-service"
    HttpsListenerReady             = aws_alb_listener.https_forward_listener.id
    TaskExecutionPolicyAttachment  = aws_iam_role_policy_attachment.task_execution.id
    PrivateSubnetDefaultRouteReady = join(":", [for route in aws_route.private_route : route.id])
  }
}

resource "aws_cloudwatch_log_group" "ecs" {
  name              = "/ecs/${local.project}-${local.env}"
  retention_in_days = 30

  tags = {
    Name = "${local.project}-${local.env}-ecs-log-group"
  }
}

# Assume role policy (信頼ポリシー)
data "aws_iam_policy_document" "ecs_tasks_assume" {
  statement {
    actions = ["sts:AssumeRole"]
    principals {
      type        = "Service"
      identifiers = ["ecs-tasks.amazonaws.com"]
    }
  }
}

# Task Execution Role
resource "aws_iam_role" "task_execution" {
  name               = "${local.project}-${local.env}-task-exec-role"
  assume_role_policy = data.aws_iam_policy_document.ecs_tasks_assume.json
  tags = {
    Name = "${local.project}-${local.env}-task-exec-role"
  }
}

resource "aws_iam_role_policy_attachment" "task_execution" {
  role = aws_iam_role.task_execution.name
  # managed
  policy_arn = "arn:aws:iam::aws:policy/service-role/AmazonECSTaskExecutionRolePolicy"
}

# ECS task definitionにsecretを注入するために必要
# AWS IAM Roleに複数のIAM Role Policyは紐付けできる
resource "aws_iam_role_policy" "task_execution_read_app_db_secret" {
  name = "${local.project}-${local.env}-task-exec-read-app-db-secret"
  role = aws_iam_role.task_execution.id

  policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Effect = "Allow"
        Action = [
          "secretsmanager:GetSecretValue"
        ]
        Resource = aws_secretsmanager_secret.db_app.arn
      }
    ]
  })
}

# Task Role
# 最初はnginxをたてるだけで、ECS exec設定するまでssmの権限もいらない
resource "aws_iam_role" "task" {
  name               = "${local.project}-${local.env}-task-role"
  assume_role_policy = data.aws_iam_policy_document.ecs_tasks_assume.json
  tags = {
    Name = "${local.project}-${local.env}-task-role"
  }
}

resource "aws_iam_role_policy" "task_read_db_secrets" {
  name = "${local.project}-${local.env}-task-read-db-secrets"
  role = aws_iam_role.task.id

  policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Effect = "Allow"
        Action = [
          "secretsmanager:GetSecretValue"
        ]
        Resource = [
          aws_secretsmanager_secret.db_admin.arn,
          aws_secretsmanager_secret.db_app.arn
        ]
      }
    ]
  })
}
