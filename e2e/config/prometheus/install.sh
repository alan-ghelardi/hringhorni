#!/usr/bin/env bash

set -euo pipefail
cur_dir="$( cd "$( dirname "${BASH_SOURCE[0]}" )" > /dev/null && pwd )"

kubectl create namespace monitoring

helm install -f $cur_dir/values.yaml h8i oci://registry-1.docker.io/bitnamicharts/kube-prometheus
