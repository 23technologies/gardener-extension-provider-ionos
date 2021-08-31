#!/bin/bash
#
# Copyright (c) 2019 SAP SE or an SAP affiliate company. All rights reserved. This file is licensed under the Apache Software License, v. 2 except as noted otherwise in the LICENSE file
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#      http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

set -o errexit
set -o nounset
set -o pipefail

rm -f "${GOPATH}/bin/*-gen"
rm -fR "${GOPATH}/src/github.com/23technologies/gardener-extension-provider-ionos"

PROJECT_ROOT="$(dirname $0)/.."

mkdir -p "${GOPATH}/src/github.com/23technologies"
ln -s "$(realpath -L ${PROJECT_ROOT})" "${GOPATH}/src/github.com/23technologies/gardener-extension-provider-ionos"

bash "${PROJECT_ROOT}/vendor/k8s.io/code-generator/generate-internal-groups.sh" \
  deepcopy,defaulter,conversion \
  github.com/23technologies/gardener-extension-provider-ionos/pkg/client \
  github.com/23technologies/gardener-extension-provider-ionos/pkg/ionos \
  github.com/23technologies/gardener-extension-provider-ionos/pkg/ionos \
  "apis:v1alpha1" \
  --go-header-file "${PROJECT_ROOT}/vendor/github.com/gardener/gardener/hack/LICENSE_BOILERPLATE.txt"

bash "${PROJECT_ROOT}/vendor/k8s.io/code-generator/generate-internal-groups.sh" \
  deepcopy,defaulter,conversion \
  github.com/23technologies/gardener-extension-provider-ionos/pkg/client/config \
  github.com/23technologies/gardener-extension-provider-ionos/pkg/ionos/apis \
  github.com/23technologies/gardener-extension-provider-ionos/pkg/ionos/apis \
  "config:v1alpha1" \
  --go-header-file "${PROJECT_ROOT}/vendor/github.com/gardener/gardener/hack/LICENSE_BOILERPLATE.txt"
