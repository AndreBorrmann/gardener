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
	"net"
	"regexp"

	k8s "k8s.io/utils/net"
)

// The CidrPair stores the IPv4 and IPv6 CIDRS if present
type CidrPair struct {
	DualStack bool
	IpNets    []*net.IPNet
}

// Parse the CidrPair from a String
func ParseCidrs(s string) (CidrPair, error) {
	split := regexp.MustCompile(", *")
	cidrs := split.Split(s, -1)
	dualstack, err := k8s.IsDualStackCIDRStrings(cidrs)
	if err != nil {
		return CidrPair{}, err
	}
	ipnets, err := k8s.ParseCIDRs(cidrs)
	if err != nil {
		return CidrPair{}, err
	}

	return CidrPair{
		DualStack: dualstack, IpNets: ipnets,
	}, nil
}

// Parse the CidrPair from a String
func MustParseCidrs(s string) CidrPair {
	cp, err := ParseCidrs(s)
	if err != nil {
		panic(err)
	}

	return cp
}

func (cp CidrPair) Cidr4() *net.IPNet {
	if k8s.IsIPv4CIDR(cp.IpNets[0]) {
		return cp.IpNets[0]
	}

	return cp.IpNets[1]
}

func (cp CidrPair) Cidr6() *net.IPNet {
	if k8s.IsIPv6CIDR(cp.IpNets[0]) {
		return cp.IpNets[0]
	}

	return cp.IpNets[1]
}

func (cp CidrPair) IsDualStack() bool { return cp.DualStack }

func (cp CidrPair) Is4() bool { return len(cp.IpNets) == 1 && k8s.IsIPv4CIDR(cp.IpNets[0]) }

func (cp CidrPair) Is6() bool { return len(cp.IpNets) == 1 && k8s.IsIPv6CIDR(cp.IpNets[0]) }

func (cp CidrPair) String() string {
	if cp.IsDualStack() {
		return fmt.Sprintf("%s,%s", cp.IpNets[0], cp.IpNets[1])
	}
	return fmt.Sprint(cp.IpNets[0])
}
