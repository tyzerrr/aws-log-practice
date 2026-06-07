# scripts

CI/CD から利用する補助 script を置くディレクトリです。

ここに置く script は、GitHub Actions の YAML に長い shell 処理を直接書かないためのものです。特に migration workflow では「差分検知」「既存 migration の保護」「後続 step への値の受け渡し」が必要になり、workflow file だけに書くと意図が読み取りづらくなります。

## detect-new-migrations.sh

`db/migrations/*.sql` の差分を調べ、DB migration workflow を実行する必要があるかを判定します。

使い方:

```bash
scripts/detect-new-migrations.sh <base-ref> <head-ref>
```

GitHub Actions では push event の `before` と `sha` を渡します。

```bash
scripts/detect-new-migrations.sh "${{ github.event.before }}" "${{ github.sha }}"
```

この script が行うこと:

- `db/migrations/*.sql` に新規追加された migration があるかを判定する。
- 既存 migration file の変更、削除、rename、copy を検知したら失敗する。
- 今回追加された migration を除いた、直前の migration version を `base_version` として出力する。
- GitHub Actions 上では `$GITHUB_OUTPUT` に書き込み、後続 step から `steps.<id>.outputs.*` として参照できるようにする。
- ローカル実行時は標準出力に同じ内容を出す。

出力:

```text
has_added_migration=true
base_version=20260527204929
added_count=1
added_migrations<<EOF
db/migrations/20260607123000_add_example.sql
EOF
```

`has_added_migration` が `false` の場合、migration image build、drift check、ECS one-shot task による migration apply はスキップします。

## CI で達成したいこと

DB migration workflow の目的は、Git に追加された migration を RDS に安全に適用することです。

想定フロー:

1. `detect-new-migrations.sh` で、新規 migration が追加された push かを判定する。
2. 新規 migration がなければ workflow を終了する。
3. 新規 migration があれば migration image を build し、ECR に push する。
4. ECS one-shot task で drift check を実行する。
5. drift check が通った場合だけ、ECS one-shot task で `atlas migrate apply` を実行する。

drift check では、今回追加された migration をまだ適用していない状態の shadow database と、remote RDS の現在 schema を比較します。差分がある場合、RDS に手動変更や未管理変更が入っている可能性があるため、自動 migration apply は止めます。

GitHub Actions runner から RDS へ直接接続しない設計にしています。RDS は private subnet にあり、security group も ECS task からの接続だけを許可しているためです。DB password も GitHub Actions には渡さず、ECS task role が Secrets Manager から admin credential を読み、task 実行時に DSN を組み立てます。

## 既存 migration を変更しない理由

Atlas は migration directory と RDS 側の migration 履歴を照合します。過去に適用済みの migration file を後から変更すると、Git 上の migration file と RDS に記録された履歴が一致しなくなります。

そのため、自動 workflow では既存 migration の変更を許可しません。過去 migration の修正が必要な場合は、基本的に新しい migration file を追加して forward-only に修正します。
