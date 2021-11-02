// Copyright (c) 2020 SAP SE or an SAP affiliate company. All rights reserved. This file is licensed under the Apache Software License, v. 2 except as noted otherwise in the LICENSE file
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cidr

import (
	"fmt"

	v1beta1constants "github.com/gardener/gardener/pkg/apis/core/v1beta1/constants"
	"github.com/gardener/gardener/pkg/utils/cidrs"
	"k8s.io/apimachinery/pkg/util/validation/field"
)

// ValidateNetworkDisjointedness validates that the given <seedNetworks> and <k8sNetworks> are disjoint.
func ValidateNetworkDisjointedness(fldPath *field.Path, shootNodes, shootPods, shootServices, seedNodes *string, seedPods, seedServices string) field.ErrorList {
	var (
		allErrs = field.ErrorList{}

		pathNodes    = fldPath.Child("nodes")
		pathServices = fldPath.Child("services")
		pathPods     = fldPath.Child("pods")
	)

	if shootNodes != nil && seedNodes != nil && NetworksIntersect(*shootNodes, *seedNodes) {
		allErrs = append(allErrs, field.Invalid(pathNodes, *shootNodes, "shoot node network intersects with seed node network"))
	}
	if shootNodes != nil && NetworksIntersect(*shootNodes, v1beta1constants.DefaultVpnRange) {
		allErrs = append(allErrs, field.Invalid(pathNodes, *shootNodes, fmt.Sprintf("shoot node network intersects with default vpn network (%s)", v1beta1constants.DefaultVpnRange)))
	}

	if shootServices != nil {
		if NetworksIntersect(seedServices, *shootServices) {
			allErrs = append(allErrs, field.Invalid(pathServices, *shootServices, "shoot service network intersects with seed service network"))
		}
		if NetworksIntersect(seedPods, *shootServices) {
			allErrs = append(allErrs, field.Invalid(pathServices, *shootServices, "shoot service network intersects with seed pod network"))
		}
		if NetworksIntersect(v1beta1constants.DefaultVpnRange, *shootServices) {
			allErrs = append(allErrs, field.Invalid(pathServices, *shootServices, fmt.Sprintf("shoot service network intersects with default vpn network (%s)", v1beta1constants.DefaultVpnRange)))
		}
	} else {
		allErrs = append(allErrs, field.Required(pathServices, "services is required"))
	}

	if shootPods != nil {
		if NetworksIntersect(seedPods, *shootPods) {
			allErrs = append(allErrs, field.Invalid(pathPods, *shootPods, "shoot pod network intersects with seed pod network"))
		}
		if NetworksIntersect(seedServices, *shootPods) {
			allErrs = append(allErrs, field.Invalid(pathPods, *shootPods, "shoot pod network intersects with seed service network"))
		}
		if NetworksIntersect(v1beta1constants.DefaultVpnRange, *shootPods) {
			allErrs = append(allErrs, field.Invalid(pathPods, *shootPods, fmt.Sprintf("shoot pod network intersects with default vpn network (%s)", v1beta1constants.DefaultVpnRange)))
		}
	} else {
		allErrs = append(allErrs, field.Required(pathPods, "pods is required"))
	}

	return allErrs
}

// NetworksIntersect returns true if the given network CIDRs intersect.
func NetworksIntersect(cidr1, cidr2 string) bool {
	cidr_pair1, err1 := cidrs.ParseCidrs(cidr1)
	cidr_pair2, err2 := cidrs.ParseCidrs(cidr2)

	if err1 != nil || err2 != nil {
		return true
	}

	if cidr_pair1.IsDualStack() && cidr_pair2.IsDualStack() {
		return cidr_pair2.Cidr4().Contains(cidr_pair1.Cidr4().IP) || cidr_pair1.Cidr4().Contains(cidr_pair2.Cidr4().IP) || cidr_pair2.Cidr6().Contains(cidr_pair1.Cidr6().IP) || cidr_pair1.Cidr6().Contains(cidr_pair2.Cidr6().IP)
	}

	if cidr_pair1.Is4() && (cidr_pair2.Is4() || cidr_pair2.IsDualStack()) {
		return cidr_pair2.Cidr4().Contains(cidr_pair1.Cidr4().IP) || cidr_pair1.Cidr4().Contains(cidr_pair2.Cidr4().IP)
	}

	if cidr_pair1.Is6() && (cidr_pair2.Is6() || cidr_pair2.IsDualStack()) {
		return cidr_pair2.Cidr6().Contains(cidr_pair1.Cidr6().IP) || cidr_pair1.Cidr6().Contains(cidr_pair2.Cidr6().IP)
	}

	if (cidr_pair1.Is4() || cidr_pair1.IsDualStack()) && cidr_pair2.Is4() {
		return cidr_pair2.Cidr4().Contains(cidr_pair1.Cidr4().IP) || cidr_pair1.Cidr4().Contains(cidr_pair2.Cidr4().IP)
	}

	if (cidr_pair1.Is6() || cidr_pair1.IsDualStack()) && cidr_pair2.Is6() {
		return cidr_pair2.Cidr6().Contains(cidr_pair1.Cidr6().IP) || cidr_pair1.Cidr6().Contains(cidr_pair2.Cidr6().IP)
	}

	return false
}
