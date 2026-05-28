resource "aws_amplify_app" "web" {
  name        = "${local.project}-${local.env}-web"
  repository  = local.github_repository
  platform    = "WEB_COMPUTE"
  description = "${local.project} ${local.env} frontend"

  build_spec = <<-YAML
    version: 1
    applications:
      - appRoot: web
        frontend:
          phases:
            preBuild:
              commands:
                - corepack enable
                - pnpm install --frozen-lockfile
            build:
              commands:
                - pnpm run build
          artifacts:
            baseDirectory: .next
            files:
              - '**/*'
          cache:
            paths:
              - node_modules/**/*
              - .next/cache/**/*
  YAML

  environment_variables = {
    AMPLIFY_MONOREPO_APP_ROOT = "web"
    AMPLIFY_DIFF_DEPLOY       = "false"
  }

  custom_rule {
    source = "/<*>"
    target = "/index.html"
    status = "404-200"
  }

  tags = {
    Name = "${local.project}-${local.env}-web"
  }

  lifecycle {
    ignore_changes = [
      iam_service_role_arn,
    ]
  }
}

resource "aws_amplify_branch" "develop" {
  app_id            = aws_amplify_app.web.id
  branch_name       = "develop"
  stage             = "DEVELOPMENT"
  framework         = "Next.js - SSR"
  enable_auto_build = true
}
