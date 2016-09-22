#!/bin/bash
set -e
set -u
set -o pipefail

"${__dir}/build.sh"

UNIT_TEST="go_test"
test_unit
