# aws-log-practice

ECS Fargate と RDS for PostgreSQL を用いたアプリケーション基盤を構築し、CloudWatch Logs と Kinesis Data Firehose を使った Production Ready なログ運用を学ぶためのプロジェクトです。

現在は AWS 上で ALB + ECS Fargate + RDS for PostgreSQL の基盤を構築し、Terraform と ecspresso を使い分けて ECS service / one-shot task を運用する構成にしています。DB credential は Secrets Manager で管理し、batch task から application 用 DB user の作成・権限確認を行います。

## 目的

- ECS Fargate でバックエンドアプリケーションを運用するための基礎を作る
- RDS for PostgreSQL を private subnet に配置し、ECS task から接続する
- CloudWatch Logs を起点に、Kinesis Data Firehose を使ったログ配送を構築する
- Terraform と ecspresso を使い分け、インフラ管理とアプリケーションデプロイの責務を分ける
- GitHub Actions から ECR への image push と ECS deploy を自動化する
- Secrets Manager と one-shot batch task を使い、DB credential を application image や Git に埋め込まない構成にする

## ディレクトリ構成

```text
.
├── .github/workflows/        # CI/CD workflows
├── db/                       # DB migration と SQL query
│   ├── migrations/
│   └── query/
├── docs/                     # 設計メモ、TODO
├── proto/                    # Protocol Buffers 定義
├── server/                   # Go backend
│   ├── cmd/batch/            # ECS one-shot batch entrypoint
│   ├── cmd/server/           # application entrypoint
│   ├── gen/                  # generated code
│   └── internal/
│       ├── adapter/          # DB や transaction など外部接続の adapter
│       ├── domain/           # entity, repository interface, value object
│       ├── infrastructure/   # repository 実装など
│       └── usecase/          # application usecase
├── terraform/dev/aws/        # dev 環境の AWS Terraform
├── ecspresso/                # ECS service / one-shot task の ecspresso 定義
└── web/                      # frontend
```

## Backend

`server/` は Go の backend アプリケーションです。ドメイン層、ユースケース層、インフラ層を分け、DB アクセスは `sqlc` で生成したコードを使う構成です。

主な構成:

- `server/cmd/server`: サーバー起動処理と handler
- `server/cmd/batch`: ECS one-shot task として実行する運用 batch
- `server/internal/domain`: entity、repository interface、value object
- `server/internal/usecase`: application usecase
- `server/internal/infrastructure`: repository 実装
- `server/internal/adapter`: DB、transaction などの adapter
- `server/gen`: protobuf / Connect-Go 生成コード

採用技術:

- Go
- PostgreSQL
- pgx
- sqlc
- Protocol Buffers
- Connect-Go

## Batch Tasks

`server/cmd/batch` には、ECS one-shot task として実行する運用 batch を置いています。

現在の command:

- `create_db_app_user`: Secrets Manager の admin/app credential を読み、admin user で RDS に接続して application 用 DB user を作成または更新します。既存 table / sequence と今後作成される table / sequence に必要な権限を付与します。
- `check_db_app_user`: application 用 DB user で実際に RDS に接続し、database / schema / table / sequence の権限が付いていることを確認します。
- `check_db_migration_drift`: migration 適用前の shadow database と remote RDS の schema に drift がないことを確認します。
- `apply_db_migration`: migration image に同梱された Atlas migration を remote RDS に適用します。

DB の host、port、DB名、subnet、security group、ECS role などの環境依存値は ecspresso が Terraform state から読み込みます。DB password は Docker image や GitHub Actions の build 時には埋め込まず、task 実行時に task role の権限で Secrets Manager から取得します。

batch の ecspresso 定義は `ecspresso/batch/` にあります。
DB migration 用の ecspresso 定義は `ecspresso/migration/` にあります。

```bash
cd ecspresso

IMAGE_URI=<batch-image-uri> \
ecspresso --envfile ../server/cmd/batch/.env run \
  --config ecspresso.config.yaml \
  --task-def batch/ecs-task-def.json \
  --wait
```

権限確認だけを実行する場合:

```bash
cd ecspresso

IMAGE_URI=<batch-image-uri> \
ecspresso --envfile ../server/cmd/batch/.env run \
  --config ecspresso.config.yaml \
  --task-def batch/ecs-check-db-app-user-task-def.json \
  --wait
```

task role には admin / app 両方の DB secret に対する `secretsmanager:GetSecretValue` 権限が必要です。

### DB Migration

DB migration は Atlas を同梱した migration image を ECR に push し、ECS one-shot task として実行します。RDS は private subnet にあるため、GitHub Actions runner から RDS へ直接接続しません。

GitHub Actions の `.github/workflows/migrate-db.yml` は次の順序で動きます。

1. `scripts/detect-new-migrations.sh` で `db/migrations/*.sql` の新規追加を検知する。
2. migration image を build し、`batch-migrate-db-<short-sha>-<yyyymmddHHMMSS>` tag で ECR に push する。
3. `ecspresso/migration/ecs-check-db-migration-drift-task-def.json` で drift check を実行する。
4. drift check が成功した場合だけ、`ecspresso/migration/ecs-apply-db-migration-task-def.json` で `atlas migrate apply` を実行する。

GitHub Actions には DB password を渡しません。`DB_ADMIN_CREDENTIAL_ID` は Secrets Manager の secret name である `aws-log-practice/dev/db/admin` を workflow 内に固定値として置きます。ECS task role が実行時に Secrets Manager から admin credential を読みます。

初回 bootstrap のように RDS にまだ migration が適用されていない場合は、`workflow_dispatch` で `force_run=true` にし、`migration_base_version` は空のまま実行します。これにより空の shadow database と remote RDS を比較してから、pending migration をすべて適用します。既存 schema に対して追加 migration だけを適用する通常運用では、CI が新規 migration file から base version を自動算出します。

dev 環境では initial migration `20260527204929_initial_schema.sql` を RDS に適用済みです。`MIGRATION_BASE_VERSION=20260527204929` の drift check が成功し、その後 `check_db_app_user` が `checks_count=18` で成功しています。

migration で新しい table / sequence を追加した後は、application 用 DB user の既存 object 権限を同期するために `create_db_app_user` を再実行し、続けて `check_db_app_user` で確認します。

## Frontend

`web/` は Next.js の frontend アプリケーションです。

採用技術:

- Next.js
- React
- TypeScript
- Tailwind CSS
- ShadCN UI style components
- TanStack Query
- ConnectRPC
- Biome
- pnpm

## Local Development

ローカルで動かす場合は、DB、backend、frontend をそれぞれ起動します。
以下のコマンドは、特記がない限りリポジトリルートで実行します。

### 1. Database

PostgreSQL を Docker Compose で起動します。

```bash
make db-up
```

初回、または migration を更新した後は schema を適用します。

```bash
make db-apply
```

DB を停止して volume も削除する場合:

```bash
make db-down
```

### 2. Backend

backend は ConnectRPC の HTTP server として起動します。
local development では `DATABASE_URL` と `SERVER_PORT` を指定します。

```bash
DATABASE_URL="postgres://aws-log-practice:aws-log-practice@localhost:5432/aws-log-practice?sslmode=disable" \
SERVER_PORT=8080 \
go run ./server/cmd/server
```

起動後、backend は `http://localhost:8080` で待ち受けます。

ECS では password 入り `DATABASE_URL` を渡さず、`DB_HOST` / `DB_PORT` / `DB_NAME` / `DB_USER` / `DB_PASSWORD` / `DB_SSLMODE` から backend が DSN を組み立てます。`DB_USER` と `DB_PASSWORD` は Secrets Manager の application DB secret から ECS task definition の `secrets` で注入します。

### 3. Frontend

frontend は `web/` で依存関係を入れてから Next.js dev server を起動します。

```bash
cd web
corepack pnpm install
corepack pnpm run dev
```

起動後、frontend は `http://localhost:3000` で確認できます。

### Generated Code

Protocol Buffers を変更した場合は、backend と frontend の生成コードを更新します。

```bash
make buf-gen

cd web
corepack pnpm run proto:gen
```

## Infrastructure

`terraform/dev/aws/` に dev 環境の AWS リソースを Terraform で定義しています。

現在の主な構成:

- VPC
- public/private subnet
- NAT Gateway
- ALB
- ACM certificate
- Route 53 record
- ECS cluster
- ECS task definition
- ECS service
- ECR repository
- RDS for PostgreSQL
- Secrets Manager
- CloudWatch Logs log group
- GitHub Actions OIDC 用 IAM role

今後追加する主な構成:

- Kinesis Data Firehose
- ログ保存先
- terraform-docs によるドキュメント生成

採用技術:

- Terraform
- AWS ECS Fargate
- AWS ALB
- AWS ECR
- AWS RDS for PostgreSQL
- AWS CloudWatch Logs
- AWS Kinesis Data Firehose
- AWS Route 53
- AWS ACM
- GitHub Actions OIDC
- ecspresso
- terraform-docs

## Deploy And Operations

ECS service と one-shot task の task definition / service definition は ecspresso で管理します。Terraform は VPC、RDS、ECS cluster、IAM role、ECR、CloudWatch Logs などのインフラを管理し、ecspresso は Terraform state の値を参照して ECS の実行定義を組み立てます。

Application image と batch image は同じ ECR repository に push しますが、tag prefix で用途を分けています。

- application: `app-<short-sha>-<yyyymmddHHMMSS>`
- batch: `batch-<batch-name>-<short-sha>-<yyyymmddHHMMSS>`

batch image の例:

```text
aws-log-practice-dev-ecr-repository:batch-create-db-app-user-128f310-20260607114532
aws-log-practice-dev-ecr-repository:batch-migrate-db-56f3535-20260607133000
```

DB credential rotation / sync は、Terraform で RDS master password と admin secret の version を揃えて更新し、その後 `create_db_app_user` batch を実行して application 用 user の password と権限を更新する流れです。application は admin credential を使わず、Secrets Manager の app credential だけで接続します。
