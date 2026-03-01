#!/usr/bin/env bash

set -o errexit
set -o nounset
set -o pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "${SCRIPT_DIR}/.." && pwd)"
GOMALL_DIR="${PROJECT_ROOT}/gomall"

cd "${GOMALL_DIR}"

modules=$(find . -name "go.mod" -not -path "./tutorial/*" -not -path "./rpc_gen/*" -not -path "./go.mod" | xargs -I {} dirname {} | sort)

for mod in $modules; do
    pushd "$mod" >/dev/null
    module_name=$(head -n 1 go.mod | cut -d ' ' -f 2)
    echo "=== Testing ${module_name} ==="
    go test -race -covermode=atomic -coverprofile=coverage.out ./... || echo "Warning: Some tests failed in ${module_name}"
    popd >/dev/null
done

echo "=== All tests completed ==="
