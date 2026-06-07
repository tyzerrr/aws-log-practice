# TODO

## Current State

- Terraform で dev 環境の VPC、private/public subnet、ALB、ECS Fargate、ECR、RDS for PostgreSQL、Secrets Manager、CloudWatch Logs、GitHub Actions OIDC role を管理している。
- ecspresso を導入済み。ECS service と one-shot batch task は ecspresso が Terraform state を参照して実行定義を組み立てる。
- `create_db_app_user` batch は ECS one-shot task として実行済み。Secrets Manager の admin/app credential を読み、application 用 DB user を作成または更新できる。
- `check_db_app_user` batch は ECS one-shot task として実行済み。application 用 DB user で RDS に接続でき、現時点で確認可能な database/schema 権限は通っている。
- Batch image は GitHub Actions から ECR に push できる。tag は `batch-<batch-name>-<short-sha>-<yyyymmddHHMMSS>` 形式。
- Atlas migration 用の ECS one-shot task definition と GitHub Actions workflow を追加済み。GitHub Actions は migration file を検知し、migration image を ECR に push し、ecspresso 経由で drift check と migration apply を実行する。

## P0: RDS Migration CI を適用して確認する

現状の `check_db_app_user` は `checks_count=2` で成功している。これは application 用 DB user が DB に接続でき、`public` schema を使えることは確認できているという意味。

一方で、`products` / `stocks` / `transactions` / `orders` などの table 権限チェックはまだ出ていない。RDS 上に application schema が未適用、または対象 table がまだ存在しない状態と考えられる。

migration one-shot task と CI workflow の実装は入っているが、AWS 側へ Terraform の IAM 変更を apply し、実際の GitHub Actions / ecspresso 実行で確認する作業が残っている。

やること:

- Terraform apply で GitHub OIDC role に ecspresso run 用の権限を反映する。
- GitHub repository secret に `DB_ADMIN_CREDENTIAL_ID` を設定する。これは DB password ではなく Secrets Manager の secret ID。
- migration workflow を実行し、drift check と `atlas migrate apply` が成功することを確認する。
- migration 実行後に `create_db_app_user` を再実行し、既存 table / sequence への GRANT を反映する。
- その後 `check_db_app_user` を再実行し、table 権限の check が増えてすべて成功することを確認する。

完了条件:

- RDS に application table が作成されている。
- GitHub Actions の migration workflow が `batch-migrate-db-<short-sha>-<yyyymmddHHMMSS>` image を ECR に push し、ECS one-shot task を完走できる。
- `check_db_app_user` の `checks_count` が database/schema だけでなく table 権限分も含む数になる。
- `SELECT` / `INSERT` / `UPDATE` / `DELETE` が対象 table ですべて成功している。

## P1: Backend Application の RDS 接続方式を ECS 環境変数に合わせる

ecspresso の application task definition は、RDS 接続情報を分解して渡している。

- `DB_HOST`: Terraform state の primary RDS endpoint
- `DB_PORT`: Terraform state の RDS port
- `DB_NAME`: Terraform state の DB name
- `DB_USER`: Secrets Manager の app secret から注入
- `DB_PASSWORD`: Secrets Manager の app secret から注入

しかし backend application code は現状 `DATABASE_URL` だけを読んでいる。

```go
db.NewDBPool(ctx, logger, os.Getenv("DATABASE_URL"))
```

つまり、ECS task definition が渡している env と application が期待している env が噛み合っていない。

やること:

- backend 起動時の DB 接続文字列生成を整理する。
- `DATABASE_URL` が設定されていればそれを優先する。これは local development を壊さないため。
- `DATABASE_URL` が空なら、`DB_HOST` / `DB_PORT` / `DB_NAME` / `DB_USER` / `DB_PASSWORD` / `DB_SSLMODE` から PostgreSQL DSN を組み立てる。
- password は URL encode する。
- ECS では app credential だけを使う。admin credential は application container に渡さない。

完了条件:

- local では従来どおり `DATABASE_URL` で起動できる。
- ECS では `DB_*` env と Secrets Manager 注入値だけで RDS に接続できる。
- application container に password 入り `DATABASE_URL` を直接置かない。

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
