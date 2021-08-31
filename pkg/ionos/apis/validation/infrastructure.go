/*
Copyright (c) 2021 SAP SE or an SAP affiliate company. All rights reserved.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

// Package validation contains functions to validate controller specifications
package validation

import (
	"fmt"
	"github.com/23technologies/gardener-extension-provider-ionos/pkg/ionos/apis"
)

// ValidateInfrastructureConfigSpec validates provider specification to check if all fields are present and valid
//
// PARAMETERS
// spec *apis.InfrastructureConfig Provider specification to validate
func ValidateInfrastructureConfigSpec(spec *apis.InfrastructureConfig) []error {
	var allErrs []error

	if nil != spec.FloatingPool {
		if "" == spec.FloatingPool.Name {
			allErrs = append(allErrs, fmt.Errorf("floatingPool.name is required field"))
		}

		if spec.FloatingPool.Size < 1 {
			allErrs = append(allErrs, fmt.Errorf("floatingPool.size is required field"))
		}
	}

	if nil != spec.Networks && "" == spec.Networks.Workers {
		allErrs = append(allErrs, fmt.Errorf("networks.workers is required field"))
	}

	return allErrs
}
