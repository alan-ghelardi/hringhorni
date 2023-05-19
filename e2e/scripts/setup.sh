#!/usr/bin/env bash

set -euo pipefail
cur_dir="$( cd "$( dirname "${BASH_SOURCE[0]}" )" > /dev/null && pwd )"

kind create cluster --config ${cur_dir}/../kind.yaml

${cur_dir}/../config/prometheus/install.sh
