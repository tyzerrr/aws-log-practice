output "github_iam_role_arn" {
  value       = aws_iam_role.role.arn
  description = "Github OIDC用のIAM Role ARN"
}

output "db_admin_secret_arn" {
  value       = aws_secretsmanager_secret.db_admin.arn
  description = "DB admin credentialを保存するSecrets Manager secretのARN"
}

output "db_app_secret_arn" {
  value       = aws_secretsmanager_secret.db_app.arn
  description = "Application用DB credentialを保存するSecrets Manager secretのARN"
}

output "db_primary_host" {
  value       = aws_db_instance.primary.address
  description = "RDS primary instanceのendpoint hostname"
}

output "db_read_replica_host" {
  value       = aws_db_instance.read_replica.address
  description = "RDS read replica instanceのendpoint hostname"
}

output "db_port" {
  value       = aws_db_instance.primary.port
  description = "RDS PostgreSQLの接続port"
}

output "http_port" {
  value       = local.http_port
  description = "ALB target groupとECS application containerが使うHTTP port"
}

output "primary_db_name" {
  value       = aws_db_instance.primary.db_name
  description = "RDS primary database name"
}
