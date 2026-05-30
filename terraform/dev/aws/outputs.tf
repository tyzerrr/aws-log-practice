output "github_iam_role_arn" {
  value       = aws_iam_role.role.arn
  description = "Github OIDC用のIAM Role ARN"
}
