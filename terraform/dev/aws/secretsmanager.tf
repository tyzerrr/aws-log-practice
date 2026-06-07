resource "aws_secretsmanager_secret" "db_admin" {
  name = "${local.project}/${local.env}/db/admin"

  tags = {
    Name = "${local.project}-${local.env}-db-admin-secrets"
  }
}

resource "aws_secretsmanager_secret" "db_app" {
  name = "${local.project}/${local.env}/db/app"

  tags = {
    Name = "${local.project}-${local.env}-db-app-secrets"
  }
}

resource "aws_secretsmanager_secret_version" "db_admin" {
  secret_id = aws_secretsmanager_secret.db_admin.id

  secret_string_wo = jsonencode({
    username = local.db_admin_username
    password = ephemeral.aws_secretsmanager_random_password.db_admin.random_password
  })

  secret_string_wo_version = 1
}

resource "aws_secretsmanager_secret_version" "db_app" {
  secret_id = aws_secretsmanager_secret.db_app.id

  secret_string_wo = jsonencode({
    username = local.db_app_username
    password = ephemeral.aws_secretsmanager_random_password.db_app.random_password
  })

  secret_string_wo_version = 1
}



