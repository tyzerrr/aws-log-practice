# Batch Tasks

`create_db_app_user` を ECS の one-shot task として実行するための ecspresso 定義です。

実行は `ecspresso` ディレクトリから行います。

```sh
cd ecspresso

IMAGE_URI=<account-id>.dkr.ecr.ap-northeast-1.amazonaws.com/<repository>:<tag> \
ecspresso --envfile ../server/cmd/batch/.env run \
  --config ecspresso.config.yaml \
  --task-def batch/ecs-task-def.json \
  --wait
```

`DB_ADMIN_CREDENTIAL_ID` と `DB_APP_CREDENTIAL_ID` は `server/cmd/batch/.env` から読み込みます。
RDS の host、port、DB名、subnet、security group、log group、ECS role は既存の ecspresso config 経由で Terraform state から読み込みます。

ECS task には `DATABASE_URL` を渡しません。batch は実行時に RDS の接続情報と Secrets Manager の値から admin DSN をメモリ上で組み立てます。

task role には admin / app 両方の DB secret に対する `secretsmanager:GetSecretValue` 権限が必要です。
