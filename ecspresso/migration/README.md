# DB migration one-shot task

Atlas migration を ECS one-shot task で実行するための ecspresso 定義です。

RDS は private subnet にあり、GitHub Actions runner から直接接続しません。GitHub Actions は ECR image を build / push し、ecspresso 経由で ECS task を起動します。DB password は GitHub Actions に渡さず、task role が Secrets Manager から admin credential を読みます。

## Task definition

- `ecs-check-db-migration-drift-task-def.json`
  - `check_db_migration_drift` command を実行します。
  - 同じ ECS task 内で `postgres:17` の shadow database sidecar を起動します。
  - `MIGRATION_BASE_VERSION` までの migration を shadow database に適用し、remote RDS の現在 schema と差分がないかを確認します。
- `ecs-apply-db-migration-task-def.json`
  - `apply_db_migration` command を実行します。
  - migration image に同梱された `db/migrations` を使い、remote RDS に `atlas migrate apply` を実行します。

## 実行例

drift check:

```bash
cd ecspresso

IMAGE_URI=<migration-image-uri> \
DB_ADMIN_CREDENTIAL_ID=<secrets-manager-secret-id> \
MIGRATION_BASE_VERSION=<previous-migration-version> \
ecspresso run \
  --config ecspresso.config.yaml \
  --task-def migration/ecs-check-db-migration-drift-task-def.json \
  --wait
```

migration apply:

```bash
cd ecspresso

IMAGE_URI=<migration-image-uri> \
DB_ADMIN_CREDENTIAL_ID=<secrets-manager-secret-id> \
ecspresso run \
  --config ecspresso.config.yaml \
  --task-def migration/ecs-apply-db-migration-task-def.json \
  --wait
```

`DB_ADMIN_CREDENTIAL_ID` は password ではなく Secrets Manager の secret ID です。RDS host、port、DB name、subnet、security group、log group、ECS role は ecspresso が Terraform state から読みます。

初回 bootstrap で RDS にまだ migration が入っていない場合、`MIGRATION_BASE_VERSION` は空にします。空の shadow database と remote RDS を比較し、drift がなければ `apply_db_migration` で pending migration をすべて適用します。

## CI での順序

GitHub Actions から実行する場合は、次の順序にします。

1. `scripts/detect-new-migrations.sh` で新規 migration file を検知する。
2. migration image を build して ECR に push する。
3. `ecs-check-db-migration-drift-task-def.json` で drift check を実行する。
4. drift check が成功した場合だけ、`ecs-apply-db-migration-task-def.json` で migration apply を実行する。

migration apply 後に新規 table / sequence が増えた場合は、application DB user へ権限を付け直す必要があります。その場合は `ecspresso/batch` の `create_db_app_user` task を後続で実行します。
