#!/usr/bin/env bash
set -euo pipefail

# Integration test script for TachiCrypt
# - Builds the CLI
# - Generates sample files/directories
# - Runs encrypt/decrypt with different configurations
# - Verifies decrypted output matches originals

SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
ROOT_DIR="$(cd "${SCRIPT_DIR}/.." && pwd)"
BIN_PATH="${ROOT_DIR}/cmd/client/tachicrypt"

TEST_ROOT="${ROOT_DIR}/test/cases"
PWD_VALUE="test-password-123"

log() {
  echo "[test] $*"
}

cleanup() {
  rm -rf "${TEST_ROOT}"
}

make_inputs() {
  local in_dir="${TEST_ROOT}/inputs"
  mkdir -p "${in_dir}"

  # Single text file
  echo "Hello TachiCrypt $(date -u +%s)" > "${in_dir}/hello.txt"
  # Another text with multiple lines
  cat > "${in_dir}/lorem.txt" <<'EOF'
Lorem ipsum dolor sit amet, consectetur adipiscing elit.
Sed do eiusmod tempor incididunt ut labore et dolore magna aliqua.
Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat.
EOF

  # Binary-ish file with random bytes
  head -c 4096 </dev/urandom > "${in_dir}/random.bin" || dd if=/dev/urandom of="${in_dir}/random.bin" bs=4096 count=1 status=none

  # Nested directory tree
  local tree="${in_dir}/tree"
  mkdir -p "${tree}/a/b" "${tree}/c"
  echo "alpha" > "${tree}/a/file1.txt"
  echo "beta" > "${tree}/a/b/file2.txt"
  echo "gamma" > "${tree}/c/file3.txt"
}

build_cli() {
  log "Building CLI binary"
  (cd "${ROOT_DIR}/cmd/client" && go build -o "${BIN_PATH}")
}

encrypt() {
  local src="$1" parts="$2" outdir="$3"
  TACHICRYPT_PASSWORD="${PWD_VALUE}" "${BIN_PATH}" --hide --parts "${parts}" --data "${src}" --output "${outdir}" >/dev/null
}

decrypt() {
  local encdir="$1" outdir="$2"
  TACHICRYPT_PASSWORD="${PWD_VALUE}" "${BIN_PATH}" --unhide --data "${encdir}" --output "${outdir}" >/dev/null
}

expect_match() {
  local src="$1" outdir="$2"
  local base
  base="$(basename "$src")"
  local restored_path="${outdir}/${base}"

  if [ -d "$src" ]; then
    diff -r "$src" "$restored_path"
  else
    diff "$src" "$restored_path"
  fi
}

run_case() {
  local src="$1" parts="$2" name="$3"
  local encdir="${TEST_ROOT}/runs/${name}/enc"
  local outdir="${TEST_ROOT}/runs/${name}/out"
  mkdir -p "$encdir" "$outdir"
  log "Encrypting '${src}' with parts=${parts}"
  encrypt "$src" "$parts" "$encdir"
  log "Decrypting to '${outdir}'"
  decrypt "$encdir" "$outdir"
  log "Verifying round-trip"
  expect_match "$src" "$outdir"
  log "OK: ${name}"
}

main() {
  cleanup || true
  mkdir -p "${TEST_ROOT}"
  make_inputs
  build_cli

  local IN="${TEST_ROOT}/inputs"

  # File cases
  run_case "${IN}/hello.txt" 3 "file_hello_p3"
  run_case "${IN}/random.bin" 5 "file_random_p5"

  # Directory case
  run_case "${IN}/tree" 4 "dir_tree_p4"

  log "All integration tests passed."
}

main "$@"

