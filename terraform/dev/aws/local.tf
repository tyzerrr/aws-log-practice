locals {
  aws_region        = "ap-northeast-1"
  env               = "dev"
  project           = "aws-log-practice"
  github_repository = "https://github.com/tyzerrr/aws-log-practice"
  http_port         = 80
  https_port        = 443
  db_port = 5432
  db_admin_username = replace("${local.project}-${local.env}_db_admin", "-", "_")
  db_app_username =replace("${local.project}-${local.env}_db_app", "-", "_")
}
