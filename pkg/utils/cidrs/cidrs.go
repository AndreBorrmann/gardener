// Copyright (c) 2018 SAP SE or an SAP affiliate company. All rights reserved.
// This file is licensed under the Apache Software License, v. 2 except as noted otherwise in the LICENSE file
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

package cidrs

import (
	"fmt"
	"regexp"

	"inet.af/netaddr"
)

// The CidrPair stores the IPv4 and IPv6 CIDRS if present
type CidrPair struct {
	V4Cidr netaddr.IPPrefix
	V6Cidr netaddr.IPPrefix
}

// Parse the CidrPair from a String
func ParseCidrs(s string) (CidrPair, error) {
	split := regexp.MustCompile(", *")
	cidrs := split.Split(s, -1)
	cidr_count := len(cidrs)
	if cidr_count > 2 {
		return CidrPair{}, fmt.Errorf("cidrpair.ParseCidrs(%q): more then 2 CIDRS given", s)
	}
	cidr1, err := netaddr.ParseIPPrefix(cidrs[0])
	if err != nil {
		return CidrPair{}, err
	}

	if cidr_count > 1 {
		cidr2, err := netaddr.ParseIPPrefix(cidrs[1])
		if err != nil {
			return CidrPair{}, err
		}
		if (cidr1.IP().Is4() && cidr2.IP().Is4()) || (cidr1.IP().Is6() && cidr2.IP().Is6()) {
			return CidrPair{}, fmt.Errorf("cidrpair.ParseCidrs(%q): both CIDRS cann't have the same type (IPv4 or IPv6)", s)
		}

		if cidr1.IP().Is4() {
			return CidrPair{V4Cidr: cidr1, V6Cidr: cidr2}, nil
		} else {
			return CidrPair{V4Cidr: cidr2, V6Cidr: cidr1}, nil
		}
	} else {
		if cidr1.IP().Is4() {
			return CidrPair{V4Cidr: cidr1}, nil
		} else {
			return CidrPair{V6Cidr: cidr1}, nil
		}
	}
}

// Parse the CidrPair from a String
func MustParseCidrs(s string) CidrPair {
	cp, err := ParseCidrs(s)
	if err != nil {
		panic(err)
	}

	return cp
}

func (cp CidrPair) Cidr4() *netaddr.IPPrefix { return &cp.V4Cidr }

func (cp CidrPair) Cidr6() *netaddr.IPPrefix { return &cp.V6Cidr }

func (cp CidrPair) IsDualStack() bool { return cp.V4Cidr.IsValid() && cp.V6Cidr.IsValid() }

func (cp CidrPair) Is4() bool { return cp.V4Cidr.IsValid() }

func (cp CidrPair) Is6() bool { return cp.V6Cidr.IsValid() }

func (cp CidrPair) String() string {
	if cp.IsDualStack() {
		return fmt.Sprintf("%s,%s", cp.V4Cidr, cp.V6Cidr)
	}
	if cp.Is4() {
		return fmt.Sprint(cp.V4Cidr)
	}
	return fmt.Sprint(cp.V6Cidr)
}
