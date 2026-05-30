resource "aws_security_group" "alb_sg" {
  vpc_id = aws_vpc.vpc.id
  name   = "${local.project}-${local.env}-alb-sg"
  tags = {
    Name = "${local.project}-${local.env}-alb-sg"
  }
}

resource "aws_security_group_rule" "alb_sg_ingress_rule" {
  security_group_id = aws_security_group.alb_sg.id
  protocol          = "TCP"
  type              = "ingress"
  from_port         = local.https_port
  to_port           = local.https_port
  cidr_blocks       = ["0.0.0.0/0"]
}

resource "aws_security_group_rule" "alb_sg_http_ingress_rule" {
  security_group_id = aws_security_group.alb_sg.id
  protocol          = "TCP"
  type              = "ingress"
  from_port         = local.http_port
  to_port           = local.http_port
  cidr_blocks       = ["0.0.0.0/0"]
}

resource "aws_security_group_rule" "alb_sg_egress_rule" {
  security_group_id = aws_security_group.alb_sg.id
  type              = "egress"
  protocol          = "-1"
  from_port         = 0
  to_port           = 0
  cidr_blocks       = ["0.0.0.0/0"]
}

resource "aws_security_group" "ecs_task_sg" {
  vpc_id = aws_vpc.vpc.id
  name   = "${local.project}-${local.env}-ecs-task-sg"
  tags = {
    Name = "${local.project}-${local.env}-ecs-task-sg"
  }
}

resource "aws_security_group_rule" "ecs_task_sg_ingress_rule" {
  security_group_id        = aws_security_group.ecs_task_sg.id
  protocol                 = "TCP"
  type                     = "ingress"
  from_port                = local.http_port
  to_port                  = local.http_port
  source_security_group_id = aws_security_group.alb_sg.id
}

resource "aws_security_group_rule" "ecs_task_sg_egress_rule" {
  security_group_id = aws_security_group.ecs_task_sg.id
  type              = "egress"
  protocol          = "-1"
  from_port         = 0
  to_port           = 0
  cidr_blocks       = ["0.0.0.0/0"]
}
