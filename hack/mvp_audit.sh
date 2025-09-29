#!/usr/bin/env bash
set -euo pipefail

red() { printf "\033[31m%s\033[0m\n" "$*"; }
grn() { printf "\033[32m%s\033[0m\n" "$*"; }
ylw() { printf "\033[33m%s\033[0m\n" "$*"; }

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "$ROOT"

echo "=== Rancher Multi-Compute MVP Audit ==="
echo "Repo: $ROOT"
echo

# ---------- File & layout checks ----------
REQ_PATHS=(
  "api/multi.suse.io/v1alpha1"
  "controllers/compute-auto-operator-controller"
  "fleet/overlays/stable/VERSION.yaml"
  "fleet/overlays/lts/VERSION.yaml"
  "fleet/overlays/canary/VERSION.yaml"
  "policies"
  "ui/rancher-extension"
  "hack/e2e-kind.sh"
  "docs/operations.md"
  "docs/support-matrix.md"
)
MISS=()
for p in "${REQ_PATHS[@]}"; do
  [[ -e "$p" ]] || MISS+=("$p")
done
if ((${#MISS[@]})); then
  ylw "MISSING paths:"
  printf '  - %s\n' "${MISS[@]}"
else
  grn "OK: core layout present."
fi
echo

# ---------- Go build & tests ----------
echo "-> Checking Go build/tests..."
if make generate >/dev/null 2>&1; then grn "OK: make generate"; else red "FAIL: make generate"; fi
if make test >/dev/null 2>&1; then grn "OK: make test"; else red "FAIL: make test"; fi
if golangci-lint run >/dev/null 2>&1; then grn "OK: golangci-lint"; else ylw "SKIP: golangci-lint not available"; fi
echo

# ---------- Coverage ----------
if [[ -f coverage.out ]]; then
  TOTAL=$(go tool cover -func=coverage.out | grep total: | awk '{print $3}')
  echo "Coverage: $TOTAL"
else
  ylw "Coverage report missing (run 'make test' with coverage)."
fi
echo

# ---------- KIND e2e smoke ----------
if [[ -x hack/e2e-kind.sh ]]; then
  echo "Run: make e2e-kind"
else
  ylw "No e2e-kind.sh executable."
fi
echo

# ---------- Docs ----------
for f in docs/operations.md docs/support-matrix.md; do
  if [[ -s "$f" ]]; then grn "OK: $f exists"; else ylw "MISSING or empty: $f"; fi
done

echo
echo "=== Audit complete. ==="