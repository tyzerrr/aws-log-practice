output "github_iam_role_arn" {
  value       = aws_iam_role.role.arn
  description = "Github OIDC用のIAM Role ARN"
}

output "db_admin_secret_arn" {
  value = aws_secretsmanager_secret.db_admin.arn
}

output "db_app_secret_arn" {
  value = aws_secretsmanager_secret.db_app.arn
}

output "db_primary_host" {
  value = aws_db_instance.primary.address
}

output "db_read_replica_host" {
  value = aws_db_instance.read_replica.address
}

output "db_port" {
  value = aws_db_instance.primary.port
}

output "http_port" {
  value = local.http_port
}

output "primary_db_name" {
  value = aws_db_instance.primary.db_name
}
