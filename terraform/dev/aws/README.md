# Terraform Dev AWS

dev 環境の AWS infrastructure を管理する Terraform root module です。

この module は VPC、subnet、ALB、ECS Fargate、ECR、RDS for PostgreSQL、Secrets Manager、CloudWatch Logs、Route 53、Amplify、GitHub Actions OIDC role を管理します。

`Requirements`、`Providers`、`Inputs`、`Outputs` は terraform-docs で生成します。手動で編集する場合は `BEGIN_TF_DOCS` / `END_TF_DOCS` の外側だけを変更してください。

リポジトリルートから更新する場合:

```bash
make terraform-docs
```

<!-- BEGIN_TF_DOCS -->
## Requirements

| Name | Version |
|------|---------|
| <a name="requirement_terraform"></a> [terraform](#requirement\_terraform) | ~> 1.15.0 |
| <a name="requirement_aws"></a> [aws](#requirement\_aws) | ~> 6.46 |
| <a name="requirement_http"></a> [http](#requirement\_http) | ~> 3.6 |
| <a name="requirement_tls"></a> [tls](#requirement\_tls) | ~> 4.3 |

## Providers

| Name | Version |
|------|---------|
| <a name="provider_aws"></a> [aws](#provider\_aws) | 6.47.0 |
| <a name="provider_http"></a> [http](#provider\_http) | 3.6.0 |
| <a name="provider_tls"></a> [tls](#provider\_tls) | 4.3.0 |

## Inputs

| Name | Description | Type | Default | Required |
|------|-------------|------|---------|:--------:|
| <a name="input_az_count"></a> [az\_count](#input\_az\_count) | Availability zone count | `number` | `2` | no |
| <a name="input_domain_name"></a> [domain\_name](#input\_domain\_name) | domain | `string` | `"aws-log-practice.com"` | no |
| <a name="input_github_organization"></a> [github\_organization](#input\_github\_organization) | Github organization name | `string` | `"tyzerrr"` | no |
| <a name="input_github_repository"></a> [github\_repository](#input\_github\_repository) | Github repository name | `string` | `"aws-log-practice"` | no |
| <a name="input_subnet_cidr_allocated_bit"></a> [subnet\_cidr\_allocated\_bit](#input\_subnet\_cidr\_allocated\_bit) | New allocated bit from VPC | `number` | `8` | no |
| <a name="input_vpc_cidr_block"></a> [vpc\_cidr\_block](#input\_vpc\_cidr\_block) | VPCに割り当てるCIDR Block | `string` | `"10.0.0.0/16"` | no |

## Outputs

| Name | Description |
|------|-------------|
| <a name="output_db_admin_secret_arn"></a> [db\_admin\_secret\_arn](#output\_db\_admin\_secret\_arn) | DB admin credentialを保存するSecrets Manager secretのARN |
| <a name="output_db_app_secret_arn"></a> [db\_app\_secret\_arn](#output\_db\_app\_secret\_arn) | Application用DB credentialを保存するSecrets Manager secretのARN |
| <a name="output_db_port"></a> [db\_port](#output\_db\_port) | RDS PostgreSQLの接続port |
| <a name="output_db_primary_host"></a> [db\_primary\_host](#output\_db\_primary\_host) | RDS primary instanceのendpoint hostname |
| <a name="output_db_read_replica_host"></a> [db\_read\_replica\_host](#output\_db\_read\_replica\_host) | RDS read replica instanceのendpoint hostname |
| <a name="output_github_iam_role_arn"></a> [github\_iam\_role\_arn](#output\_github\_iam\_role\_arn) | Github OIDC用のIAM Role ARN |
| <a name="output_http_port"></a> [http\_port](#output\_http\_port) | ALB target groupとECS application containerが使うHTTP port |
| <a name="output_primary_db_name"></a> [primary\_db\_name](#output\_primary\_db\_name) | RDS primary database name |
<!-- END_TF_DOCS -->
