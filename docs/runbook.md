# Runbook

このドキュメントは、dev 環境の運用手順をまとめる runbook です。
ここに書く手順は、AWS 上の RDS / ECS / Secrets Manager / ecspresso を操作するものです。

## DB Credential Rotation / App User 権限同期

### 目的

RDS admin credential と application 用 DB credential を安全に更新し、application 用 DB user に必要な権限を同期します。

この手順では DB password を Git、Docker image、GitHub Actions の build log、Terraform state に直接置きません。Terraform は Secrets Manager の secret version と RDS password を更新し、ECS one-shot batch が実行時に Secrets Manager から credential を読みます。

### 前提

- AWS profile は dev 環境を操作できるものを使います。
- Terraform backend の state を読める必要があります。
- ecspresso は `ecspresso/ecspresso.config.yaml` から Terraform state を参照します。
- batch image は ECR に push 済みの `batch-create-db-app-user-<short-sha>-<yyyymmddHHMMSS>` tag を使います。
- `server/cmd/batch/.env` には次の secret name が入っています。
  - `DB_ADMIN_CREDENTIAL_ID=aws-log-practice/dev/db/admin`
  - `DB_APP_CREDENTIAL_ID=aws-log-practice/dev/db/app`

### 変更対象

credential rotation で変更する Terraform local は `terraform/dev/aws/local.tf` です。

```hcl
db_admin_password_version = 2
db_app_password_version   = 1
```

- admin password を rotate する場合は `db_admin_password_version` を increment します。
- app password を rotate する場合は `db_app_password_version` を increment します。
- app credential だけを rotate する場合でも、application 用 DB user の password を RDS 側へ反映するために `create_db_app_user` を実行します。

### 手順

1. Terraform の password version を変更します。

   `terraform/dev/aws/local.tf` の `db_admin_password_version` または `db_app_password_version` を 1 つ増やします。

2. Terraform apply で RDS / Secrets Manager を更新します。

   ```bash
   cd terraform/dev/aws

   terraform plan
   terraform apply
   ```

   `db_admin_password_version` を変更した場合、RDS master password と admin secret version が同じ password に更新されます。
   `db_app_password_version` を変更した場合、Secrets Manager の app secret version が更新されます。

3. batch image を用意します。

   GitHub Actions の `build-and-push-batch-image.yml` を実行し、summary に出た image URI を控えます。

   image URI の例:

   ```text
   957573575820.dkr.ecr.ap-northeast-1.amazonaws.com/aws-log-practice-dev-ecr-repository:batch-create-db-app-user-a4b8d7f-20260607121122
   ```

4. `create_db_app_user` を ECS one-shot task として実行します。

   ```bash
   cd ecspresso

   AWS_PROFILE=taichi-aws-log-practice \
   AWS_REGION=ap-northeast-1 \
   IMAGE_URI=<batch-create-db-app-user-image-uri> \
   ecspresso --envfile ../server/cmd/batch/.env run \
     --config ecspresso.config.yaml \
     --task-def batch/ecs-task-def.json \
     --wait
   ```

   成功時は次のような log が出ます。

   ```json
   {"level":"INFO","msg":"created or updated db app user","db_name":"aws_log_practice_primary","schema":"public","username":"aws_log_practice_dev_db_app"}
   ```

5. `check_db_app_user` を ECS one-shot task として実行します。

   ```bash
   cd ecspresso

   AWS_PROFILE=taichi-aws-log-practice \
   AWS_REGION=ap-northeast-1 \
   IMAGE_URI=<batch-create-db-app-user-image-uri> \
   ecspresso --envfile ../server/cmd/batch/.env run \
     --config ecspresso.config.yaml \
     --task-def batch/ecs-check-db-app-user-task-def.json \
     --wait
   ```

   成功時は次のような log が出ます。

   ```json
   {"level":"INFO","msg":"checked db app user privileges","db_name":"aws_log_practice_primary","schema":"public","username":"aws_log_practice_dev_db_app","checks_count":18}
   ```

6. application を再起動または再 deploy します。

   ECS task definition の `secrets` は task 起動時に Secrets Manager から値を注入します。既存の running task は古い app password を持ち続ける可能性があるため、app credential を rotate した後は ECS service を新しい task に入れ替えます。

   backend image の変更がある場合は application deploy workflow を実行します。image の変更がない場合は ECS service の force new deployment で十分です。

   ```bash
   aws ecs update-service \
     --cluster aws-log-practice-dev-ecs-cluster \
     --service aws-log-practice-dev-ecs-service \
     --force-new-deployment
   ```

### 確認項目

- `create_db_app_user` が `created or updated db app user` で終了している。
- `check_db_app_user` が `checked db app user privileges` で終了している。
- `check_db_app_user` の `checks_count` が想定値から減っていない。
- ECS service の running task が rotation 後に起動した task に入れ替わっている。
- application が RDS に接続できている。

### 失敗時の見方

- `AccessDeniedException` が出る場合は、ECS task role が admin / app secret に `secretsmanager:GetSecretValue` できるか確認します。
- `password authentication failed` が出る場合は、Terraform apply 後の RDS password と Secrets Manager の admin secret version がずれている可能性があります。`db_admin_password_version` と `aws_db_instance.primary.password_wo_version`、`aws_secretsmanager_secret_version.db_admin.secret_string_wo_version` が同じ値を参照しているか確認します。
- `check_db_app_user` が権限不足で失敗する場合は、`create_db_app_user` を再実行します。migration で table / sequence が増えた直後は、既存 object への GRANT を同期するために `create_db_app_user` が必要です。
- local で `InvalidAccessKeyId` が出る場合は、壊れた環境変数の AWS credential が `AWS_PROFILE` より優先されている可能性があります。`AWS_ACCESS_KEY_ID` / `AWS_SECRET_ACCESS_KEY` / `AWS_SESSION_TOKEN` を unset してから再実行します。

### セキュリティ上の注意

- admin credential は application runtime では使いません。
- application runtime には app credential だけを渡します。
- `DB_ADMIN_CREDENTIAL_ID` と `DB_APP_CREDENTIAL_ID` は password ではなく Secrets Manager の secret name です。
- Docker build 時に DB password を渡しません。
- GitHub Actions に DB password を secret として渡しません。
