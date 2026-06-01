# aws-log-practice

ECS Fargate と RDS for PostgreSQL を用いたアプリケーション基盤を構築し、CloudWatch Logs と Kinesis Data Firehose を使った Production Ready なログ運用を学ぶためのプロジェクトです。

現在は AWS 上で ALB + ECS Fargate の HTTPS 疎通まで確認済みで、今後 RDS 接続、ecspresso による ECS デプロイ、CI/CD、ログ配送基盤を順に整備していきます。

## 目的

- ECS Fargate でバックエンドアプリケーションを運用するための基礎を作る
- RDS for PostgreSQL を private subnet に配置し、ECS task から接続する
- CloudWatch Logs を起点に、Kinesis Data Firehose を使ったログ配送を構築する
- Terraform と ecspresso を使い分け、インフラ管理とアプリケーションデプロイの責務を分ける
- GitHub Actions から ECR への image push と ECS deploy を自動化する

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
│   ├── cmd/server/           # application entrypoint
│   ├── gen/                  # generated code
│   └── internal/
│       ├── adapter/          # DB や transaction など外部接続の adapter
│       ├── domain/           # entity, repository interface, value object
│       ├── infrastructure/   # repository 実装など
│       └── usecase/          # application usecase
├── terraform/dev/aws/        # dev 環境の AWS Terraform
└── web/                      # frontend
```

## Backend

`server/` は Go の backend アプリケーションです。ドメイン層、ユースケース層、インフラ層を分け、DB アクセスは `sqlc` で生成したコードを使う構成です。

主な構成:

- `server/cmd/server`: サーバー起動処理と handler
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
`DATABASE_URL` と `SERVER_PORT` が必要です。

```bash
DATABASE_URL="postgres://aws-log-practice:aws-log-practice@localhost:5432/aws-log-practice?sslmode=disable" \
SERVER_PORT=8080 \
go run ./server/cmd/server
```

起動後、backend は `http://localhost:8080` で待ち受けます。

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
- CloudWatch Logs log group
- GitHub Actions OIDC 用 IAM role

今後追加する主な構成:

- RDS for PostgreSQL
- Secrets Manager または SSM Parameter Store
- Kinesis Data Firehose
- ログ保存先
- ecspresso 設定
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
