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
- Backend application 用 Dockerfile と GitHub Actions workflow を追加済み。workflow は application image を buildx で build/push し、ecspresso deploy で ECS service を更新する。

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
