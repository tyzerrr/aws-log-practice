# Batch Tasks

DB application user 関連の one-shot ECS task を実行するための ecspresso 定義です。

実行は `ecspresso` ディレクトリから行います。

## create_db_app_user

application 用 DB user を作成または更新し、必要な権限を付与します。

```sh
cd ecspresso

IMAGE_URI=<account-id>.dkr.ecr.ap-northeast-1.amazonaws.com/<repository>:<tag> \
ecspresso --envfile ../server/cmd/batch/.env run \
  --config ecspresso.config.yaml \
  --task-def batch/ecs-task-def.json \
  --wait
```

## check_db_app_user

application 用 DB user で実際に DB 接続し、必要な権限が付いていることを確認します。

```sh
cd ecspresso

IMAGE_URI=<account-id>.dkr.ecr.ap-northeast-1.amazonaws.com/<repository>:<tag> \
ecspresso --envfile ../server/cmd/batch/.env run \
  --config ecspresso.config.yaml \
  --task-def batch/ecs-check-db-app-user-task-def.json \
  --wait
```

`DB_ADMIN_CREDENTIAL_ID` と `DB_APP_CREDENTIAL_ID` は `server/cmd/batch/.env` から読み込みます。
RDS の host、port、DB名、subnet、security group、log group、ECS role は既存の ecspresso config 経由で Terraform state から読み込みます。

ECS task には `DATABASE_URL` を渡しません。batch は実行時に RDS の接続情報と Secrets Manager の値から DSN をメモリ上で組み立てます。

task role には admin / app 両方の DB secret に対する `secretsmanager:GetSecretValue` 権限が必要です。
