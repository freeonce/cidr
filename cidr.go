// Copyright 2021 helloshaohua <wu.shaohua@foxmail.com>;
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

package httpclient

import (
	"fmt"
	"math"
	"net"
	"strconv"
)

// CIDRClassType type definition.
type CIDRClassType string

// String is implemented by any value that has a String method,
// which defines the ``native'' format for that value.
// The String method is used to print values passed as an operand
// to any format that accepts a string or to an unformatted printer
// such as Print.
func (C CIDRClassType) String() string {
	return string(C)
}

// CIDR class type related constant definition.
const (
	CIDRClassTypeOfAMask       CIDRClassType = "ANet"
	CIDRClassTypeOfBMask       CIDRClassType = "BNet"
	CIDRClassTypeOfCMask       CIDRClassType = "CNet"
	CIDRClassTypeOfUnknownMask CIDRClassType = "UnknownNet"
)

// Default route masks for IPv4.
var (
	classAMask = net.IPv4Mask(0xff, 0, 0, 0)
	classBMask = net.IPv4Mask(0xff, 0xff, 0, 0)
	classCMask = net.IPv4Mask(0xff, 0xff, 0xff, 0)
)

// CIDRMaskType returns the CIDR mask class type.
func CIDRMaskType(cidr string) (CIDRClassType, error) {
	ip, _, err := net.ParseCIDR(cidr)
	if err != nil {
		return CIDRClassTypeOfUnknownMask, err
	}

	var r CIDRClassType

	switch ip.DefaultMask().String() {
	case classAMask.String():
		r = CIDRClassTypeOfAMask
	case classBMask.String():
		r = CIDRClassTypeOfBMask
	case classCMask.String():
		r = CIDRClassTypeOfCMask
	default:
		r = CIDRClassTypeOfUnknownMask
	}

	return r, err
}

// IPRangeToCIDR Convert IPv4 range into CIDR.
func IPRangeToCIDR(ipStart string, ipEnd string) ([]string, error) {
	var (
		cidr2mask = []uint32{
			0x00000000, 0x80000000, 0xC0000000,
			0xE0000000, 0xF0000000, 0xF8000000,
			0xFC000000, 0xFE000000, 0xFF000000,
			0xFF800000, 0xFFC00000, 0xFFE00000,
			0xFFF00000, 0xFFF80000, 0xFFFC0000,
			0xFFFE0000, 0xFFFF0000, 0xFFFF8000,
			0xFFFFC000, 0xFFFFE000, 0xFFFFF000,
			0xFFFFF800, 0xFFFFFC00, 0xFFFFFE00,
			0xFFFFFF00, 0xFFFFFF80, 0xFFFFFFC0,
			0xFFFFFFE0, 0xFFFFFFF0, 0xFFFFFFF8,
			0xFFFFFFFC, 0xFFFFFFFE, 0xFFFFFFFF,
		}
		ipStartUint32 = uint32(IPv4ToInt64(ipStart))
		ipEndUint32   = uint32(IPv4ToInt64(ipEnd))
		cidrs         = make([]string, 0)
	)

	if ipStartUint32 > ipEndUint32 {
		return nil, fmt.Errorf("start ip %s must be less than end ip %s", ipStart, ipEnd)
	}

	for ipEndUint32 >= ipStartUint32 {
		maxSize := 32
		for maxSize > 0 {
			maskedBase := ipStartUint32 & cidr2mask[maxSize-1]
			if maskedBase != ipStartUint32 {
				break
			}
			maxSize--
		}

		x := math.Log(float64(ipEndUint32-ipStartUint32+1)) / math.Log(2)
		maxDiff := 32 - int(math.Floor(x))
		if maxSize < maxDiff {
			maxSize = maxDiff
		}

		cidrs = append(cidrs, Int64ToIPv4(int64(ipStartUint32))+"/"+strconv.Itoa(maxSize))

		ipStartUint32 += uint32(math.Exp2(float64(32 - maxSize)))
	}

	return cidrs, nil
}
