#!/usr/bin/env bash
set -euo pipefail

# GitHub Actions の push event で渡される before/after commit を受け取り、
# db/migrations 配下に「新規追加された SQL migration」があるかだけを判定する。
#
# この script は migration apply そのものは行わない。役割は次の 3 つに限定する。
# 1. 新規 migration が追加されているかを CI に伝える。
# 2. 既存 migration の変更・削除・rename/copy を検知して CI を止める。
# 3. drift check 用に「新規追加分を除いた最後の migration version」を算出する。

# Git の空 tree object。初回 push などで GitHub Actions の before SHA が
# 000000... になる場合、比較元 commit が存在しないため空 tree と比較する。
readonly EMPTY_TREE="4b825dc642cb6eb9a060e54bf8d69288fbee4904"

declare -a ADDED_MIGRATIONS=()
declare -a INVALID_CHANGES=()

main() {
  local base_ref="${1:?base ref is required}"
  local head_ref="${2:?head ref is required}"
  local diff_output

  base_ref="$(normalize_base_ref "${base_ref}")"
  diff_output="$(migration_diff "${base_ref}" "${head_ref}")"

  collect_migration_changes "${diff_output}"
  fail_if_invalid_migration_changes

  if (( ${#ADDED_MIGRATIONS[@]} == 0 )); then
    output_no_added_migration
    return
  fi

  output_added_migrations "$(previous_migration_version)"
}

normalize_base_ref() {
  local base_ref="${1}"

  if [[ "${base_ref}" =~ ^0+$ ]] || ! git cat-file -e "${base_ref}^{commit}" 2>/dev/null; then
    echo "${EMPTY_TREE}"
    return
  fi

  echo "${base_ref}"
}

migration_diff() {
  local base_ref="${1}"
  local head_ref="${2}"

  # --name-status は A/M/D/R/C などの status と path を tab 区切りで返す。
  # CI では db/migrations だけを対象にし、アプリケーションコードや Terraform の変更は
  # migration 実行可否の判定材料にしない。
  git diff --name-status "${base_ref}" "${head_ref}" -- db/migrations || true
}

collect_migration_changes() {
  local diff_output="${1}"
  local status
  local path
  local rest
  local target_path

  while IFS=$'\t' read -r status path rest; do
    [[ -n "${status:-}" ]] || continue

    target_path="$(changed_target_path "${status}" "${path}" "${rest:-}")"

    # atlas.sum など SQL 以外のファイル変更はここでは扱わない。
    # 実際の checksum 整合性は ECS task 内で atlas migrate apply が検証する。
    [[ "${target_path}" == db/migrations/*.sql ]] || continue

    record_migration_change "${status}" "${path}" "${rest:-}"
  done <<< "${diff_output}"
}

changed_target_path() {
  local status="${1}"
  local path="${2}"
  local rest="${3}"

  # Rename/Copy は "R100 old_path new_path" のように path が 2 つ出る。
  # 判定対象は最終的に置かれる path なので、new_path 側を target path として返す。
  if [[ "${status}" == R* || "${status}" == C* ]]; then
    echo "${rest:-${path}}"
    return
  fi

  echo "${path}"
}

record_migration_change() {
  local status="${1}"
  local path="${2}"
  local rest="${3}"

  # 自動 migration workflow では「追加」だけを許可する。
  # 既存 migration の変更・削除・rename/copy は、すでに適用済みの RDS 履歴と
  # Git 上の migration directory が食い違う原因になるため、手動対応に回す。
  if [[ "${status}" == "A" ]]; then
    ADDED_MIGRATIONS+=("${path}")
    return
  fi

  INVALID_CHANGES+=("${status} ${path}${rest:+ -> ${rest}}")
}

fail_if_invalid_migration_changes() {
  if (( ${#INVALID_CHANGES[@]} == 0 )); then
    return
  fi

  {
    echo "Existing migration files must not be modified, deleted, renamed, or copied in the automatic migration workflow."
    printf 'Invalid change: %s\n' "${INVALID_CHANGES[@]}"
  } >&2
  exit 1
}

previous_migration_version() {
  # drift check では「今回追加された migration をまだ適用していない状態」と
  # remote RDS の現在 schema を比較する。そのため、現在の working tree から
  # 追加分を差し引き、残った migration の最後の version を base_version とする。
  #
  # ファイル名は Atlas の timestamp prefix 形式を前提にしている。
  # 例: db/migrations/20260527204929_initial_schema.sql -> 20260527204929
  comm -23 \
    <(git ls-files 'db/migrations/*.sql' | sort) \
    <(printf '%s\n' "${ADDED_MIGRATIONS[@]}" | sort) |
    sed -E 's#^.*/([0-9]+)_.+$#\1#' |
    sort |
    tail -n 1
}

output_no_added_migration() {
  {
    echo "has_added_migration=false"
    echo "base_version="
    echo "added_count=0"
  } | write_outputs
}

output_added_migrations() {
  local base_version="${1}"

  {
    echo "has_added_migration=true"
    echo "base_version=${base_version}"
    echo "added_count=${#ADDED_MIGRATIONS[@]}"
    echo "added_migrations<<EOF"
    printf '%s\n' "${ADDED_MIGRATIONS[@]}"
    echo "EOF"
  } | write_outputs
}

write_outputs() {
  # GitHub Actions 上では step output として後続 step から参照できるように
  # $GITHUB_OUTPUT に書き込む。ローカル実行時は動作確認しやすいよう標準出力へ出す。
  if [[ -n "${GITHUB_OUTPUT:-}" ]]; then
    cat >> "${GITHUB_OUTPUT}"
    return
  fi

  cat
}

main "$@"
