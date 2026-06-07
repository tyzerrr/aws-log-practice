# TODO

## Current State

- Terraform で dev 環境の VPC、private/public subnet、ALB、ECS Fargate、ECR、RDS for PostgreSQL、Secrets Manager、CloudWatch Logs、GitHub Actions OIDC role を管理している。
- ecspresso を導入済み。ECS service と one-shot batch task は ecspresso が Terraform state を参照して実行定義を組み立てる。
- `create_db_app_user` batch は ECS one-shot task として実行済み。Secrets Manager の admin/app credential を読み、application 用 DB user を作成または更新できる。
- `check_db_app_user` batch は ECS one-shot task として実行済み。application 用 DB user で RDS に接続でき、database / schema / table 権限チェックが `checks_count=18` で成功している。
- Batch image は GitHub Actions から ECR に push できる。tag は `batch-<batch-name>-<short-sha>-<yyyymmddHHMMSS>` 形式。
- Atlas migration 用の ECS one-shot task definition と GitHub Actions workflow を追加済み。GitHub Actions から migration image を ECR に push し、ecspresso 経由で drift check と migration apply を実行できる。
- RDS primary database には initial migration `20260527204929_initial_schema.sql` を適用済み。`MIGRATION_BASE_VERSION=20260527204929` の drift check も成功している。
- Backend application は `DATABASE_URL` があればそれを優先し、無ければ `DB_HOST` / `DB_PORT` / `DB_NAME` / `DB_USER` / `DB_PASSWORD` / `DB_SSLMODE` から PostgreSQL DSN を組み立てる。ecspresso の application task definition も `DB_SSLMODE=require` を渡す。

## Done: RDS Migration CI / App User 権限確認

RDS migration を ECS one-shot task として実行する流れは確認済み。

完了したこと:

- GitHub OIDC role に ecspresso run 用の IAM policy を追加した。
- GitHub Actions workflow `.github/workflows/migrate-db.yml` を追加した。
- migration image を `batch-migrate-db-<short-sha>-<yyyymmddHHMMSS>` tag で ECR に push できる。
- `check_db_migration_drift` task で shadow database と remote RDS の schema drift を確認できる。
- `apply_db_migration` task で Atlas migration を RDS に適用できる。
- initial migration 適用後、`MIGRATION_BASE_VERSION=20260527204929` の drift check が成功した。
- migration 適用後に `check_db_app_user` を実行し、`checks_count=18` で application 用 DB user の権限チェックが成功した。

運用上の注意:

- `DB_ADMIN_CREDENTIAL_ID` は DB password ではなく Secrets Manager の secret name なので、CI workflow 内の固定値 `aws-log-practice/dev/db/admin` を使う。
- 初回 bootstrap では `workflow_dispatch` の `force_run=true` を使い、`migration_base_version` は空のままにする。
- 追加 migration の通常運用では、CI が追加 migration file から base version を自動算出する。
- migration で新しい table / sequence が増えた場合は、`create_db_app_user` を再実行して既存 object への GRANT を反映し、その後 `check_db_app_user` で確認する。

## Done: Backend Application の RDS 接続方式を ECS 環境変数に合わせる

Backend application の RDS 接続方式は ECS task definition と噛み合うように修正済み。

完了したこと:

- `DATABASE_URL` が設定されていればそれを優先する。local development では従来どおり `DATABASE_URL` で起動できる。
- `DATABASE_URL` が空なら、`DB_HOST` / `DB_PORT` / `DB_NAME` / `DB_USER` / `DB_PASSWORD` / `DB_SSLMODE` から PostgreSQL DSN を組み立てる。
- password は `net/url` の `url.UserPassword` 経由で URL encode する。
- `DB_PORT` は未指定なら `5432`、`DB_SSLMODE` は未指定なら `require` を使う。
- ecspresso の application task definition に `DB_SSLMODE=require` を追加した。
- unit test で `DATABASE_URL` 優先、`DB_*` からの DSN 組み立て、default port / sslmode、必須 env 不足時の error を確認している。

運用上の注意:

- ECS では application container に password 入り `DATABASE_URL` を直接置かない。
- application runtime には app credential だけを渡す。admin credential は migration / batch task だけで使う。

## P2: Backend Application Image / Deploy を整備する

Batch image の CI は動いているが、application image の Dockerfile / deploy はまだ整理が必要。

やること:

- application 用 Dockerfile を作成する。
- GitHub Actions で application image を build し、ECR に push する。
- tag は `app-<short-sha>-<yyyymmddHHMMSS>` 形式にする。
- ecspresso で ECS service を deploy できるようにする。
- ALB 経由で backend application に HTTPS 疎通できることを確認する。

完了条件:

- ECR に application image が push される。
- ecspresso deploy で ECS service が新しい image に更新される。
- backend が RDS に接続した状態で起動する。

## P3: Frontend 開発

Frontend は backend / DB / deploy が安定してから API 結合すると進めやすい。ただし、画面設計や mock API 前提の UI 開発は P1/P2 と並行して進められる。

やること:

- 商品、在庫、注文、取引など現在の backend domain に対応する画面を整理する。
- ConnectRPC client 経由で backend API を呼ぶ。
- local backend と AWS backend の接続先切り替えを整理する。
- 最低限の一覧、作成、詳細、状態確認の導線を作る。

完了条件:

- local で frontend から backend API を呼べる。
- AWS backend に向けた接続設定が整理されている。
- 主要 domain の基本操作を画面から確認できる。

## P4: DB Credential Rotation / Batch 運用を手順化する

現状は手動で以下の流れを実行できる。

1. Terraform で RDS master password と admin secret の version を揃えて更新する。
2. `create_db_app_user` batch を ECS one-shot task として実行する。
3. `check_db_app_user` batch を ECS one-shot task として実行する。

やること:

- rotation 手順を runbook 化する。
- 必要なら workflow_dispatch で batch image build と ecspresso run を実行できるようにする。
- admin credential は application runtime から分離し、batch/task role だけが読める状態を保つ。
- app credential の rotation 時に application deploy / restart が必要か整理する。

完了条件:

- credential 更新手順が README または docs に明文化されている。
- operator が手順どおりに `create` -> `check` を実行できる。
- application が admin credential を一切使わない。

## P5: Logging / Observability

当初目的である CloudWatch Logs 起点のログ配送は、ECS/RDS/application deploy が安定してから進める。

やること:

- CloudWatch Logs の log group / stream prefix 設計を整理する。
- Kinesis Data Firehose を追加する。
- ECS application logs と batch logs の保存先を設計する。
- retention、暗号化、権限、監視を整理する。

完了条件:

- application logs と batch logs を用途別に追跡できる。
- CloudWatch Logs から保存先へ配送できる。
- Production Ready なログ運用に必要な retention / encryption / permission が定義されている。

## P6: Terraform / ecspresso / Documentation 整理

機能実装が落ち着いた後に保守性を上げる。

やること:

- terraform-docs を導入する。
- Terraform module / environment の README を生成する。
- 必要であれば CI で terraform-docs の差分チェックを行う。
- ecspresso の application / batch task definition の重複を整理する。
- task role を application 用と batch 用で分けるか検討する。

完了条件:

- Terraform の inputs / outputs が docs で追える。
- ecspresso 定義の重複が許容範囲に収まっている。
- IAM role の責務が明確になっている。
