#!/usr/bin/env bash
set -euo pipefail

# GitHub Actions の push event で渡される before/after commit を受け取り、
# db/migrations 配下に「新規追加された SQL migration」があるかだけを判定する。
#
# この script は migration apply そのものは行わない。役割は次の 3 つに限定する。
# 1. 新規 migration が追加されているかを CI に伝える。
# 2. 既存 migration の変更・削除・rename/copy を検知して CI を止める。
# 3. drift check 用に「新規追加分を除いた最後の migration version」を算出する。
base_ref="${1:?base ref is required}"
head_ref="${2:?head ref is required}"

# Git の空 tree object。初回 push などで GitHub Actions の before SHA が
# 000000... になる場合、比較元 commit が存在しないため空 tree と比較する。
empty_tree="4b825dc642cb6eb9a060e54bf8d69288fbee4904"

write_outputs() {
  # GitHub Actions 上では step output として後続 step から参照できるように
  # $GITHUB_OUTPUT に書き込む。ローカル実行時は動作確認しやすいよう標準出力へ出す。
  if [[ -n "${GITHUB_OUTPUT:-}" ]]; then
    cat >> "${GITHUB_OUTPUT}"
  else
    cat
  fi
}

if [[ "${base_ref}" =~ ^0+$ ]] || ! git cat-file -e "${base_ref}^{commit}" 2>/dev/null; then
  base_ref="${empty_tree}"
fi

# --name-status は A/M/D/R/C などの status と path を tab 区切りで返す。
# CI では db/migrations だけを対象にし、アプリケーションコードや Terraform の変更は
# migration 実行可否の判定材料にしない。
diff_output="$(git diff --name-status "${base_ref}" "${head_ref}" -- db/migrations || true)"

declare -a added_migrations=()
declare -a invalid_changes=()

while IFS=$'\t' read -r status path rest; do
  [[ -n "${status:-}" ]] || continue

  # Rename/Copy は "R100 old_path new_path" のように path が 2 つ出る。
  # 判定対象は最終的に置かれる path なので、new_path 側を target_path として扱う。
  target_path="${path}"
  if [[ "${status}" == R* || "${status}" == C* ]]; then
    target_path="${rest:-${path}}"
  fi

  # atlas.sum など SQL 以外のファイル変更はここでは扱わない。
  # 実際の checksum 整合性は ECS task 内で atlas migrate apply が検証する。
  [[ "${target_path}" == db/migrations/*.sql ]] || continue

  # 自動 migration workflow では「追加」だけを許可する。
  # 既存 migration の変更・削除・rename/copy は、すでに適用済みの RDS 履歴と
  # Git 上の migration directory が食い違う原因になるため、手動対応に回す。
  if [[ "${status}" == "A" ]]; then
    added_migrations+=("${target_path}")
  else
    invalid_changes+=("${status} ${path}${rest:+ -> ${rest}}")
  fi
done <<< "${diff_output}"

if (( ${#invalid_changes[@]} > 0 )); then
  {
    echo "Existing migration files must not be modified, deleted, renamed, or copied in the automatic migration workflow."
    printf 'Invalid change: %s\n' "${invalid_changes[@]}"
  } >&2
  exit 1
fi

if (( ${#added_migrations[@]} == 0 )); then
  {
    echo "has_added_migration=false"
    echo "base_version="
    echo "added_count=0"
  } | write_outputs
  exit 0
fi

# drift check では「今回追加された migration をまだ適用していない状態」と
# remote RDS の現在 schema を比較する。そのため、現在の working tree から
# 追加分を差し引き、残った migration の最後の version を base_version とする。
#
# ファイル名は Atlas の timestamp prefix 形式を前提にしている。
# 例: db/migrations/20260527204929_initial_schema.sql -> 20260527204929
base_version="$(
  comm -23 \
    <(git ls-files 'db/migrations/*.sql' | sort) \
    <(printf '%s\n' "${added_migrations[@]}" | sort) |
    sed -E 's#^.*/([0-9]+)_.+$#\1#' |
    sort |
    tail -n 1
)"

{
  echo "has_added_migration=true"
  echo "base_version=${base_version}"
  echo "added_count=${#added_migrations[@]}"
  echo "added_migrations<<EOF"
  printf '%s\n' "${added_migrations[@]}"
  echo "EOF"
} | write_outputs
