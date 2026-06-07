resource "aws_db_instance" "primary" {
  identifier = "${local.project}-${local.env}-db-primary"

  engine = "postgres"
  engine_version = "16"
  instance_class = "db.t3.micro"
  allocated_storage = 10

  db_name = replace("${local.project}_primary", "-", "_")
  username = local.db_admin_username
  password_wo = ephemeral.aws_secretsmanager_random_password.db_admin.random_password
  password_wo_version = 1

  storage_encrypted = true
  db_subnet_group_name = aws_db_subnet_group.rds.name

  backup_retention_period = 7

  # Need to define security groups
  # ECS taskがSecret managerから値を読み取って注入するにはIAM role policy (secretsmanager:GetSecretValue) が必要
  vpc_security_group_ids = [aws_security_group.rds_security_group.id]
  publicly_accessible = false
  multi_az = true
  skip_final_snapshot = true
}

resource "aws_db_instance" "read_replica" {
  identifier = "${local.project}-${local.env}-read-replica"
  replicate_source_db = aws_db_instance.primary.arn

  instance_class = "db.t3.micro"
  db_subnet_group_name = aws_db_subnet_group.rds.name
  # Need to define security groups as same as primary
  vpc_security_group_ids = [aws_security_group.rds_security_group.id]
  storage_encrypted = true
  publicly_accessible = false
  skip_final_snapshot = true
}

resource "aws_db_subnet_group" "rds" {
  name = "${local.project}-${local.env}-db"
  subnet_ids = [for subnet in aws_subnet.private_subnets : subnet.id]

  tags = {
    Name = "${local.project}-${local.env}-db-subnet-group"
  }
}

# AWS Secrets Manager のパスワード生成機能を利用してランダムパスワードを生成するが、
# その生成された値を Terraform の state に保存しないためのリソース
# DB admin用: Rotation担当のLambdaのみがつかう
ephemeral "aws_secretsmanager_random_password" "db_admin" {
  password_length = 32
}

# Applicationがコネクション確立時に利用する
ephemeral "aws_secretsmanager_random_password" "db_app" {
  password_length = 32
}
