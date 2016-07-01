//
//Copyright [2016] [SnapRoute Inc]
//
//Licensed under the Apache License, Version 2.0 (the "License");
//you may not use this file except in compliance with the License.
//You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
//	 Unless required by applicable law or agreed to in writing, software
//	 distributed under the License is distributed on an "AS IS" BASIS,
//	 WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//	 See the License for the specific language governing permissions and
//	 limitations under the License.
//
// _______  __       __________   ___      _______.____    __    ____  __  .___________.  ______  __    __
// |   ____||  |     |   ____\  \ /  /     /       |\   \  /  \  /   / |  | |           | /      ||  |  |  |
// |  |__   |  |     |  |__   \  V  /     |   (----` \   \/    \/   /  |  | `---|  |----`|  ,----'|  |__|  |
// |   __|  |  |     |   __|   >   <       \   \      \            /   |  |     |  |     |  |     |   __   |
// |  |     |  `----.|  |____ /  .  \  .----)   |      \    /\    /    |  |     |  |     |  `----.|  |  |  |
// |__|     |_______||_______/__/ \__\ |_______/        \__/  \__/     |__|     |__|      \______||__|  |__|
//
package rx

import (
	"fmt"
	"infra/sysd/sysdCommonDefs"
	"l3/ndp/debug"
	"log/syslog"
	"net"
	"reflect"
	"testing"
	"utils/logging"
)

var testPkt = []byte{
	0xfe, 0x80, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0xc0, 0x00, 0x54, 0xff, 0xfe, 0xf5, 0x00, 0x01,
}

func NDPTestNewLogger(name string, tag string, listenToConfig bool) (*logging.Writer, error) {
	var err error
	srLogger := new(logging.Writer)
	srLogger.MyComponentName = name

	srLogger.SysLogger, err = syslog.New(syslog.LOG_INFO|syslog.LOG_DAEMON, tag)
	if err != nil {
		fmt.Println("Failed to initialize syslog - ", err)
		return srLogger, err
	}

	srLogger.GlobalLogging = true
	srLogger.MyLogLevel = sysdCommonDefs.INFO
	return srLogger, err
}

// Test ND Solicitation message Decoder
func TestNDSolicitationDecoder(t *testing.T) {
	var err error
	logger, err := NDPTestNewLogger("ndpd", "NDPTEST", true)
	if err != nil {
		t.Error("creating logger failed")
	}
	debug.NDPSetLogger(logger)
	nds := &NDSolicitation{}
	err = DecodeNDSolicitation(testPkt, nds)
	if err != nil {
		t.Error("Decoding ipv6 and icmpv6 header failed", err)
	}
	ndWant := &NDSolicitation{
		TargetAddress: net.IP{0xfe, 0x80, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0xc0, 0x00, 0x54, 0xff, 0xfe, 0xf5, 0x00, 0x01},
	}
	if !reflect.DeepEqual(nds, ndWant) {
		t.Error("Decoding NDS Failed")
	}
}

// Test ND Solicitation multicast Address Validation
func TestNDSMulticast(t *testing.T) {
	b := net.IP{0xfe, 0x80, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0xc0, 0x00, 0x54, 0xff, 0xfe, 0xf5, 0x00, 0x01}

	// b is not multicast address, fail the test case if true is returned
	if IsNDSolicitationMulticastAddr(b) {
		t.Error("byte is not ipv6 muticast address", b)
	}

	b[0] = 0xff
	// b is multicast address, fail the test case if false is returned
	if !IsNDSolicitationMulticastAddr(b) {
		t.Error("byte is ipv6 muticast address", b)
	}
}

// Test ND Solicitation src ip Address Validation
func TestNDSIpAddress(t *testing.T) {
	srcIP := net.IP{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}
	dstIP := net.IP{0xff, 0x02, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01, 0xff, 0x10, 0x78, 0x2e}
	t.Log("SrcIP->", srcIP.String(), "DstIP->", dstIP.String())
	err := ValidateIpAddrs(srcIP, dstIP)
	if err != nil {
		t.Error("Validation of ip address failed with error", err)
	}

	srcIP = net.IP{0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}
	t.Log("SrcIP->", srcIP.String(), "DstIP->", dstIP.String())
	err = ValidateIpAddrs(srcIP, dstIP)
	if err != nil {
		t.Error("Validation of ip address", srcIP, "failed with error", err)
	}
	dstIP = net.IP{0xfe, 0x80, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0xc0, 0x00, 0x54, 0xff, 0xfe, 0xf5, 0x00, 0x01}
	t.Log("SrcIP->", srcIP.String(), "DstIP->", dstIP.String())
	err = ValidateIpAddrs(srcIP, dstIP)
	if err != nil {
		t.Error("Validation of ip address", srcIP, "dst Ip", dstIP, "failed with error", err)
	}
}
