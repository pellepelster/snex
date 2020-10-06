#!/usr/bin/env bash

set -eu

DIR="$( cd "$(dirname "$0")" ; pwd -P )"

function task_build {
  go build
}

function task_test {
    go test
    task_build

    test_return_code 1
    test_return_code 1 -source non_empty -target non_empty
    test_return_code 1 -snippets non_empty -target non_empty
    test_return_code 1 -snippets non_empty -source non_empty
    test_return_code 2 -source non_existing -target non_existing -snippets non_existing
    test_return_code 2 -source ./test/source -target non_existing -snippets non_existing
    test_return_code 0 -source ./test/source -target ./test/target -snippets ./test/snippets
}

function test_return_code {

    set +e
    local expected_return_code="${1:-}"
    shift || true
    "${DIR}/snex" "$@"

    if [ $? -ne ${expected_return_code} ]; then
        echo "expected return code ${expected_return_code} but got return code $? for command line '$@'"
        exit 1
    fi
    set -e
}

function task_usage {
  echo "Usage: $0 build | test"
  exit 1
}

arg=${1:-}
shift || true
case ${arg} in
  test) task_test "$@" ;;
  build) task_build "$@" ;;
  *) task_usage ;;
esac