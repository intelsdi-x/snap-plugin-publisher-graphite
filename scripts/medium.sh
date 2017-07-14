#!/bin/bash
set -e
set -u
set -o pipefail

__dir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
__proj_dir="$(dirname "$__dir")"
__proj_name="$(basename $__proj_dir)"

# shellcheck source=scripts/common.sh
. "${__dir}/common.sh"

_verify_docker() {
  type -p docker > /dev/null 2>&1 || _error "docker needs to be installed"
  type -p docker-compose > /dev/null 2>&1 || _error "docker-compose needs to be installed"
  docker version >/dev/null 2>&1 || _error "docker needs to be configured/running"
}

export GOLANGVER=${GOLANGVER:-"go1.7.4"}
export PLUGIN_SRC="${__proj_dir}"
TEST_TYPE="${TEST_TYPE:-"medium"}"
docker_folder="${__proj_dir}/scripts/config"
_docker_org_path="\$GOPATH/src/github.com/intelsdi-x"
_docker_proj_path="${_docker_org_path}/${__proj_name}"

_docker_project () {
  (cd "${docker_folder}" && "$@")
}

_debug "running docker compose images"
_docker_project docker-compose up -d
_debug "running test: ${TEST_TYPE}"
#TODO!!!! 
sleep 3

_docker_project docker-compose exec -T golang gvm-bash.sh "gvm install $GOLANGVER -B; gvm use $GOLANGVER; mkdir -p ${_docker_org_path}; cp -Rf /${__proj_name} ${_docker_org_path}; (cd ${_docker_proj_path} && ./scripts/medium_tests.sh)"
_debug "stopping docker compose images"
_docker_project docker-compose stop
_docker_project docker-compose rm -f
