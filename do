#!/usr/bin/env bash

set -eu

DIR="$( cd "$(dirname "$0")" ; pwd -P )"

function task_build {
  go build
}

function task_test {
    go test
    task_build

    rm -rf "${DIR}/test/distinct_folders/target"
    mkdir -p "${DIR}/test/distinct_folders/target"

    test_return_code 1
    test_return_code 1 -source non_empty -target non_empty
    test_return_code 1 -snippets non_empty -target non_empty
    test_return_code 1 -snippets non_empty -source non_empty
    test_return_code 2 -source non_existing -target non_existing -snippets non_existing
    test_return_code 2 -source ${DIR}/test/distinct_folders/source -target non_existing -snippets non_existing
    test_return_code 3 -source ${DIR}/test/distinct_folders/source -target ${DIR}/test/distinct_folders/snippets -snippets ${DIR}/test/distinct_folders/snippets
    test_return_code 3 -source ${DIR}/test/distinct_folders/source -target ${DIR}/test/distinct_folders/snippets/nested_folder -snippets ${DIR}/test/distinct_folders/snippets
    test_return_code 0 -source ${DIR}/test/distinct_folders/source -target ${DIR}/test/distinct_folders/target -snippets ${DIR}/test/distinct_folders/snippets
    diff "${DIR}/test/distinct_folders/target" "${DIR}/test/distinct_folders/expected"

    rm -rf "${DIR}/test/single_source_inside_snippets_no_target/snippets"
    cp -rv "${DIR}/test/single_source_inside_snippets_no_target/template" "${DIR}/test/single_source_inside_snippets_no_target/snippets"

    test_return_code 0 -source ${DIR}/test/single_source_inside_snippets_no_target/snippets/source1.txt -snippets ${DIR}/test/single_source_inside_snippets_no_target/snippets
    diff "${DIR}/test/single_source_inside_snippets_no_target/snippets" "${DIR}/test/single_source_inside_snippets_no_target/expected"

    rm -rf "${DIR}/test/template_in_file/target"
    mkdir -p "${DIR}/test/template_in_file/target"
    test_return_code 0 -template-file ${DIR}/test/template_in_file/test.template -source ${DIR}/test/template_in_file/source -target ${DIR}/test/template_in_file/target -snippets ${DIR}/test/template_in_file/snippets
    diff "${DIR}/test/template_in_file/target" "${DIR}/test/template_in_file/expected"

    mkdir -p "${DIR}/test/file_include_does_not_exist/target"
    mkdir -p "${DIR}/test/file_include_does_not_exist/source"
    test_return_code 4 -source ${DIR}/test/file_include_does_not_exist/source -target ${DIR}/test/file_include_does_not_exist/target -snippets ${DIR}/test/file_include_does_not_exist/snippets
}

function test_return_code {

    set +e
    local expected_return_code="${1:-}"
    shift || true
    "${DIR}/snex" "$@"

    local return_code=$?
    if [ ${return_code} -ne ${expected_return_code} ]; then
        echo "expected return code ${expected_return_code} but got return code ${return_code} for command line '$@'"
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