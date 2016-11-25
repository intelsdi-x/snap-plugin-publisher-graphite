#!/bin/bash  

set -e
set -u
set -o pipefail

__dir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
__proj_dir="$(dirname "$__dir")"

. "${__dir}/common.sh"

_info "running the example ${__proj_dir}/examples/tasks/psutil-statistics-graphite.sh"
export PLUGIN_PATH="/etc/snap/path"
source "${__proj_dir}/examples/tasks/psutil-statistics-graphite.sh"

_debug "sleeping for 10 seconds so the task can do some work"
sleep 10

# begin assertions

return_code=0
echo -n "[task is running] "
task_list=$(snaptel task list | tail -1)
if echo $task_list | grep -q Running; then
    echo "ok"
else 
    echo "not ok"
    return_code=-1
fi

echo -n "[task is hitting] "
if [ $(echo $task_list | awk '{print $4}') -gt 0 ]; then
    echo "ok"
else 
    _debug $task_list
    echo "not ok"
    return_code=-1
fi

echo -n "[task has no errors] "
if [ $(echo $task_list | awk '{print $6}') -eq 0 ]; then
    echo "ok"
else 
    echo "not ok"
    return_code=-1
fi

prefix=`curl -X GET -H "application/json" "http://graphite/metrics/find?query=myprefix" | jq -r '.[] | .text'`
echo -n "[graphite is updated] "
if [ "$prefix" = "myprefix" ]; then
    echo "ok"
else 
    echo "not ok"
    return_code=-1
fi

exit $return_code
