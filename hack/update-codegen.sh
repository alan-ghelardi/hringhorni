#!/usr/bin/env bash

# Copyright 2023 The hringhorni Authors
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

set -o errexit
set -o nounset
set -o pipefail

source $(dirname $0)/../vendor/knative.dev/hack/codegen-library.sh
export PATH="$GOBIN:$PATH"
K8S_CODEGEN="./vendor/k8s.io/code-generator/cmd"

function run_yq() {
    run_go_tool github.com/mikefarah/yq/v4@v4.23.1 yq "$@"
}

echo "=== Update Codegen for ${MODULE_NAME}"

group "Kubernetes Codegen"

# generate the code with:
# --output-base    because this script should also be able to run inside the vendor dir of
#                  k8s.io/kubernetes. The output-base is needed for the generators to output into the vendor dir
#                  instead of the $GOPATH directly. For normal projects this can be dropped.
${CODEGEN_PKG}/generate-groups.sh "deepcopy,client,informer,lister" \
              github.com/nubank/hringhorni/pkg/client github.com/nubank/hringhorni/pkg/apis \
              "rollouts:v1alpha1" \
              --go-header-file ${REPO_ROOT_DIR}/hack/boilerplate/boilerplate.go.txt

# Generate deep copy functions for other packages.
go run ${K8S_CODEGEN}/deepcopy-gen/main.go \
   -O zz_generated.deepcopy \
   --go-header-file ${REPO_ROOT_DIR}/hack/boilerplate/boilerplate.go.txt \
   --input-dirs $(echo \
                      github.com/nubank/hringhorni/pkg/apis/config \
                      | sed "s/ /,/g")

group "Knative Codegen"

# Knative Injection
${KNATIVE_CODEGEN_PKG}/hack/generate-knative.sh "injection" \
                      github.com/nubank/hringhorni/pkg/client github.com/nubank/hringhorni/pkg/apis \
                      "rollouts:v1alpha1" \
                      --go-header-file ${REPO_ROOT_DIR}/hack/boilerplate/boilerplate.go.txt

group "Update CRD Schema"

go run $(dirname $0)/../cmd/schema/ dump Analysis \
    | run_yq eval-all --header-preprocess=false --inplace 'select(fileIndex == 0).spec.versions[0].schema.openAPIV3Schema = select(fileIndex == 1) | select(fileIndex == 0)' \
             $(dirname $0)/../config/base/crd/analyses.yaml -

go run $(dirname $0)/../cmd/schema/ dump Rollout \
    | run_yq eval-all --header-preprocess=false --inplace 'select(fileIndex == 0).spec.versions[0].schema.openAPIV3Schema = select(fileIndex == 1) | select(fileIndex == 0)' \
             $(dirname $0)/../config/base/crd/rollouts.yaml -

group "Update deps post-codegen"

# Make sure our dependencies are up-to-date
${REPO_ROOT_DIR}/hack/update-deps.sh
