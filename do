#!/usr/bin/env bash

set -eu

DIR="$( cd "$(dirname "$0")" ; pwd -P )"

function task_run {
    go run "${DIR}/cmd" $@
}

function task_build {

  rm -rf "${DIR}/build"
  mkdir -p "${DIR}/build"

  declare -A targets=(["linux"]="amd64,386,arm64", ["darwin"]="arm64,amd64", ["windows"]="arm64,amd64,386", ["freebsd"]="amd64")

  for platform in "${!targets[@]}"
  do
    local archs=${targets[$platform]}
    for arch in ${archs//,/ }; do
      export GOOS="${platform}"
      export GOARCH="${arch}"
      go build -o "${DIR}/build/snex_${GOOS}_${GOARCH}" "${DIR}/cmd"
    done
  done
}

function task_test {
    go test -v "${DIR}/..."
    chmod +x ${DIR}/build/snex*

    test_return_code 5 replace non-existent-folder
    test_return_code 4 show-templates

    rm -rf "${DIR}/test-output/"
    mkdir -p "${DIR}/test-output/"

    # explicit replace
    cp -r "${DIR}/test/testbed1/input" "${DIR}/test-output/testbed1"
    go run "${DIR}/cmd" replace "${DIR}/test-output/testbed1"
    diff "${DIR}/test-output/testbed1/README.md" "${DIR}/test/testbed1/expected/README.md"

    rm -rf "${DIR}/test-output/"
    mkdir -p "${DIR}/test-output/"

    # default command
    cp -r "${DIR}/test/testbed1/input" "${DIR}/test-output/testbed1"
    go run "${DIR}/cmd" replace "${DIR}/test-output/testbed1"
    diff "${DIR}/test-output/testbed1/README.md" "${DIR}/test/testbed1/expected/README.md"

    cp -r "${DIR}/test/testbed2/input" "${DIR}/test-output/testbed2"
    go run "${DIR}/cmd" replace --template "start\n{{.Content}}\nend" "${DIR}/test-output/testbed2"
    diff "${DIR}/test-output/testbed2/README.md" "${DIR}/test/testbed2/expected/README.md"
}

function test_return_code {

    set +e
    local expected_return_code="${1:-}"
    shift || true

    ${DIR}/build/snex_linux_amd64 $@

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
  test) task_test $@ ;;
  build) task_build $@ ;;
  run) task_run $@ ;;
  *) task_usage ;;
esac