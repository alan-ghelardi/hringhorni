#!/usr/bin/env bash

set -euo pipefail
cur_dir="$( cd "$( dirname "${BASH_SOURCE[0]}" )" > /dev/null && pwd )"

# Doc: Prints the provided message to stderr, shows the script usage and exits
# with an error code.
#
# Arguments:
# $*: message to be shown.
function usage() {
    >&2  echo -e "$*"
    echo "install-api.sh <broken|initial|stable|unstable>"
    exit 1
}


export KO_DOCKER_REPO=kind.local

case "$1" in
    initial ) kustomize build ${cur_dir}/../config/ | ko apply -f - ;;
    broken | stable | unstable ) ko apply -f ${cur_dir}/../config/${1}.yaml ;;
    * ) usage "Unknown version $1" ;;
esac
