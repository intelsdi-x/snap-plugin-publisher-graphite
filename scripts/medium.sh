#!/usr/bin/env bash
# File managed by pluginsync

set -e
set -u
set -o pipefail

__dir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
__proj_dir="$(dirname "$__dir")"
__proj_name="$(basename "$__proj_dir")"

# shellcheck source=scripts/common.sh
. "${__dir}/common.sh"

_verify_docker() {
  type -p docker > /dev/null 2>&1 || _error "docker needs to be installed"
  docker version >/dev/null 2>&1 || _error "docker needs to be configured/running"
}

_verify_docker

[[ -f "${__proj_dir}/build/linux/x86_64/${__proj_name}" ]] || (cd "${__proj_dir}" && make)
SNAP_VERSION=${SNAP_VERSION:-"latest"}
OS=${OS:-"alpine"}
PLUGIN_PATH=${PLUGIN_PATH:-"${__proj_dir}"}
DEMO=${DEMO:-"false"}
TASK=${TASK:-""}

if [[ ${DEBUG:-} == "true" ]]; then
  cmd="cd /plugin/scripts && rescue rspec ./test/*_spec.rb"
else
  cmd="cd /plugin/scripts && rspec ./test/*_spec.rb"
fi

_info "running medium test"
#Starting docker with graphite
_docker_ps_id="$(docker run -d -p 80:80 -p 2003:2003 -p 2004:2004 -p 4444:4444/udp -p 8126:8126 -p 8086:8086 hopsoft/graphite-statsd)"
if ./scripts/medium_tests.sh ; then
    _info "medium test ended: succeeded"
else
    _info "medium test ended: failed"
fi
docker kill $_docker_ps_id >/dev/null && docker rm $_docker_ps_id >/dev/null
