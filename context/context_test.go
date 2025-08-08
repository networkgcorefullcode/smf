// SPDX-FileCopyrightText: 2025 Canonical Ltd
// SPDX-License-Identifier: Apache-2.0
//

package context

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/omec-project/openapi/nfConfigApi"
)

func makeSessionConfig(
	sliceName, mcc, mnc, sst string,
	sd string, dnnName, ueSubnet, hostname string, port int32,
) nfConfigApi.SessionManagement {
	sstUint64, err := strconv.ParseUint(sst, 10, 8)
	if err != nil {
		sstUint64 = 0
	}
	sstint := int32(sstUint64)

	return nfConfigApi.SessionManagement{
		SliceName: sliceName,
		PlmnId: nfConfigApi.PlmnId{
			Mcc: mcc,
			Mnc: mnc,
		},
		Snssai: nfConfigApi.Snssai{
			Sst: sstint,
			Sd:  &sd,
		},
		IpDomain: []nfConfigApi.IpDomain{
			{
				DnnName:  dnnName,
				DnsIpv4:  "8.8.8.8",
				UeSubnet: ueSubnet,
				Mtu:      1400,
			},
		},
		Upf: &nfConfigApi.Upf{
			Hostname: hostname,
			Port:     &port,
		},
		GnbNames: []string{"gnb1", "gnb2"},
	}
}

func TestUpdateSmfContext(t *testing.T) {
	validSingleSliceConfig := makeSessionConfig("slice1", "222", "02", "2", "1", "internet", "192.168.1.0/24", "upf-1", 38415)

	multiSliceConfig := []nfConfigApi.SessionManagement{
		makeSessionConfig("slice1", "111", "01", "1", "1", "internet", "192.168.1.0/24", "upf-1", 38412),
		makeSessionConfig("slice2", "111", "01", "1", "1", "fast", "192.168.2.0/24", "upf-2", 38412),
	}

	tests := []struct {
		name     string
		config   []nfConfigApi.SessionManagement
		validate func(*SMFContext, error) (bool, string)
	}{
		{
			name:   "Empty config should clear context",
			config: nil,
			validate: func(smCtx *SMFContext, err error) (bool, string) {
				if err != nil {
					return false, err.Error()
				}
				if len(smCtx.SnssaiInfos) != 0 {
					return false, "expected SnssaiInfos to be cleared"
				}
				if smCtx.UserPlaneInformation != nil && len(smCtx.UserPlaneInformation.UPNodes) != 0 {
					return false, "expected UPNodes to be cleared"
				}
				return true, ""
			},
		},
		{
			name:   "Valid single slice config",
			config: []nfConfigApi.SessionManagement{validSingleSliceConfig},
			validate: func(smCtx *SMFContext, err error) (bool, string) {
				if err != nil {
					return false, err.Error()
				}
				if len(smCtx.SnssaiInfos) != 1 {
					return false, fmt.Sprintf("expected 1 SnssaiInfo, got %d", len(smCtx.SnssaiInfos))
				}
				if smCtx.UserPlaneInformation == nil || smCtx.UserPlaneInformation.DefaultUserPlanePath == nil {
					return false, "UserPlaneInformation or DefaultUserPlanePath is nil"
				}
				if _, ok := smCtx.UserPlaneInformation.UPNodes["upf-1"]; !ok {
					return false, "expected UPNode for upf-1 to exist"
				}
				if _, ok := smCtx.UserPlaneInformation.AccessNetwork["gnb1"]; !ok {
					return false, "expected gnb1 in AccessNetwork"
				}
				if _, ok := smCtx.UserPlaneInformation.AccessNetwork["gnb2"]; !ok {
					return false, "expected gnb2 in AccessNetwork"
				}
				if len(smCtx.UserPlaneInformation.UPFIPToName) == 0 {
					return false, "expected UPFIPToName to be populated"
				}
				return true, ""
			},
		},
		{
			name:   "Multiple slice config",
			config: multiSliceConfig,
			validate: func(smCtx *SMFContext, err error) (bool, string) {
				if err != nil {
					return false, err.Error()
				}
				if len(smCtx.SnssaiInfos) != 2 {
					return false, fmt.Sprintf("expected 2 SnssaiInfos, got %d", len(smCtx.SnssaiInfos))
				}
				if _, ok := smCtx.UserPlaneInformation.UPNodes["upf-1"]; !ok {
					return false, "expected UPNode for upf-1"
				}
				if _, ok := smCtx.UserPlaneInformation.UPNodes["upf-2"]; !ok {
					return false, "expected UPNode for upf-2"
				}
				return true, ""
			},
		},
		{
			name: "Duplicate UPF should merge DNNs",
			config: []nfConfigApi.SessionManagement{
				makeSessionConfig("slice1", "111", "01", "1", "010101", "internet", "192.168.1.0/24", "upf-1", 38412),
				makeSessionConfig("slice1", "112", "01", "1", "010101", "iot", "192.168.2.0/24", "upf-1", 38412),
			},
			validate: func(smCtx *SMFContext, err error) (bool, string) {
				if len(smCtx.UserPlaneInformation.UPFs) != 1 {
					return false, "expected 1 UPF for duplicate hostname"
				}
				if len(smCtx.SnssaiInfos) != 2 {
					return false, "expected 1 SnssaiInfo for same slice"
				}
				return true, ""
			},
		},
		{
			name: "Invalid UPF IP should fallback to FQDN",
			config: []nfConfigApi.SessionManagement{
				makeSessionConfig("slice1", "112", "01", "1", "010101", "internet", "192.168.1.0/24", "bad:ip", 38412),
			},
			validate: func(smCtx *SMFContext, err error) (bool, string) {
				upf := smCtx.UserPlaneInformation.UPFs["bad:ip"]
				if upf == nil {
					return false, "expected UPF to be created with bad IP hostname"
				}
				return true, ""
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			smCtx := &SMFContext{}
			err := UpdateSmfContext(smCtx, tt.config)
			if ok, msg := tt.validate(smCtx, err); !ok {
				t.Errorf("validation failed: %s", msg)
			}
		})
	}
}
