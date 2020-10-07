#!/usr/bin/env bash

set -eu

DIR="$( cd "$(dirname "$0")" ; pwd -P )"

function task_build {
  go build
}

function task_test {
    go test
    task_build

    rm -rf "${DIR}/test/target"
    mkdir -p "${DIR}/test/target"

    test_return_code 1
    test_return_code 1 -source non_empty -target non_empty
    test_return_code 1 -snippets non_empty -target non_empty
    test_return_code 1 -snippets non_empty -source non_empty
    test_return_code 2 -source non_existing -target non_existing -snippets non_existing
    test_return_code 2 -source ./test/source -target non_existing -snippets non_existing
    test_return_code 3 -source ${DIR}/test/source -target ${DIR}/test/snippets -snippets ${DIR}/test/snippets
    test_return_code 3 -source ${DIR}/test/source -target ${DIR}/test/snippets/nested_folder -snippets ${DIR}/test/snippets
    test_return_code 0 -source ${DIR}/test/source -target ${DIR}/test/target -snippets ${DIR}/test/snippets
    diff "${DIR}/test/target" "${DIR}/test/expected"
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