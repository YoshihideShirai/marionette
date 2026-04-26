#!/usr/bin/env bash
set -euo pipefail

base_ref="${1:-}"
zero_sha="0000000000000000000000000000000000000000"

if [[ -z "$base_ref" || "$base_ref" == "$zero_sha" ]] || ! git rev-parse --verify "$base_ref" >/dev/null 2>&1; then
  if git rev-parse --verify HEAD~1 >/dev/null 2>&1; then
    base_ref="HEAD~1"
  else
    echo "No comparable base commit available; skipping TS/TSX new-file check."
    exit 0
  fi
fi

added_files="$(git diff --name-status "$base_ref"...HEAD -- '*.ts' '*.tsx' | awk '$1=="A" {print $2}')"

if [[ -n "$added_files" ]]; then
  echo "Detected newly added TypeScript files (.ts/.tsx):"
  echo "$added_files"
  echo "This repository is Go-only; adding TypeScript files is blocked by CI."
  exit 1
fi

echo "No newly added .ts/.tsx files detected."
