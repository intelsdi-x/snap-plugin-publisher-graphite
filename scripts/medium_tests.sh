#!/bin/bash

set -e
set -u
set -o pipefail

__dir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
__proj_dir="$(dirname "$__dir")"

. "${__dir}/common.sh"
. "${__dir}/deps.sh"

_go_get github.com/smartystreets/goconvey/convey
_go_get github.com/smartystreets/assertions

export TEST_TYPE="${TEST_TYPE:-"medium"}"
export SNAP_GRAPHITE_HOST="0.0.0.0"
_go_test
