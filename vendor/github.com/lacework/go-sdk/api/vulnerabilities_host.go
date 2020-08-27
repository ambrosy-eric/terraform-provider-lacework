//
// Author:: Salim Afiune Maya (<afiune@lacework.net>)
// Copyright:: Copyright 2020, Lacework Inc.
// License:: Apache License, Version 2.0
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//

package api

import (
	"fmt"
	"strings"
)

// HostVulnerabilityService is a service that interacts with the vulnerabilities
// endpoints for the host space from the Lacework Server
type HostVulnerabilityService struct {
	client *Client
}

// Scan requests an on-demand vulnerability assessment of your software packages
// to determine if the packages contain any common vulnerabilities and exposures
//
// NOTE: Only packages managed by a package manager for supported OS's are reported
func (svc *HostVulnerabilityService) Scan(manifest string) (
	// TODO @afiune add scan response. look at the end of this file
	response map[string]interface{},
	err error,
) {
	err = svc.client.RequestDecoder("POST",
		apiVulnerabilitiesScanPkgManifest,
		strings.NewReader(manifest),
		&response,
	)
	return
}

func (svc *HostVulnerabilityService) ListCves() (
	response hostVulnListCvesResponse,
	err error,
) {
	err = svc.client.RequestDecoder("GET", apiVulnerabilitiesHostListCves, nil, &response)
	return
}

func (svc *HostVulnerabilityService) ListHostsWithCVE(id string) (
	response hostVulnListHostsResponse,
	err error,
) {
	apiPath := fmt.Sprintf(apiVulnerabilitiesListHostsWithCveID, id)
	err = svc.client.RequestDecoder("GET", apiPath, nil, &response)
	return
}

func (svc *HostVulnerabilityService) GetHostAssessment(id string) (
	response hostVulnHostResponse,
	err error,
) {
	apiPath := fmt.Sprintf(apiVulnerabilitiesHostAssessmentFromMachineID, id)
	err = svc.client.RequestDecoder("GET", apiPath, nil, &response)
	return
}

type hostVulnHostResponse struct {
	Assessment HostVulnHostAssessment `json:"data"`
	Ok         bool                   `json:"ok"`
	Message    string                 `json:"message"`
}

type HostVulnHostAssessment struct {
	Host hostVulnHostDetail `json:"host"`
	CVEs []HostVulnCVE      `json:"vulnerabilities"`
}

type hostVulnListHostsResponse struct {
	Hosts   []HostVulnDetail `json:"data"`
	Ok      bool             `json:"ok"`
	Message string           `json:"message"`
}

type HostVulnDetail struct {
	Details  hostVulnHostDetail `json:"host"`
	Packages []HostVulnPackage  `json:"packages"`
	Summary  HostVulnCveSummary `json:"summary"`
}

type hostVulnHostDetail struct {
	Hostname      string      `json:"hostname"`
	MachineID     string      `json:"machine_id"`
	MachineStatus string      `json:"machine_status,omitempty"`
	Tags          hostVulnTag `json:"tags"`
}

type hostVulnTag struct {
	Account        string `json:"Account"`
	AmiID          string `json:"AmiId"`
	ExternalIP     string `json:"ExternalIp"`
	Hostname       string `json:"Hostname"`
	InstanceID     string `json:"InstanceId"`
	InternalIP     string `json:"InternalIp"`
	LwTokenShort   string `json:"LwTokenShort"`
	SubnetID       string `json:"SubnetId"`
	VmInstanceType string `json:"VmInstanceType"`
	VmProvider     string `json:"VmProvider"`
	VpcID          string `json:"VpcId"`
	Zone           string `json:"Zone"`
	Arch           string `json:"arch"`
	Os             string `json:"os"`
}

type hostVulnListCvesResponse struct {
	CVEs    []HostVulnCVE `json:"data"`
	Ok      bool          `json:"ok"`
	Message string        `json:"message"`
}

type HostVulnCVE struct {
	ID       string             `json:"cve_id"`
	Packages []HostVulnPackage  `json:"packages"`
	Summary  HostVulnCveSummary `json:"summary"`
}

type HostVulnPackage struct {
	Name                string `json:"name"`
	Namespace           string `json:"namespace"`
	Severity            string `json:"severity"`
	Status              string `json:"status,omitempty"`
	VulnerabilityStatus string `json:"vulnerabiliy_status,omitempty"` // @afiune typo
	Version             string `json:"version"`
	HostCount           string `json:"host_count"`
	PackageStatus       string `json:"package_status"`
	CveLink             string `json:"cve_link"`
	CvssScore           string `json:"cvss_score"`
	CvssV2Score         string `json:"cvss_v_2_score"`
	CvssV3Score         string `json:"cvss_v_3_score"`
	//FirstSeenTime time.Time `json:"first_seen_time"`
	FixAvailable string `json:"fix_available"`
	FixedVersion string `json:"fixed_version"`
}

func (assessment *HostVulnHostAssessment) VulnerabilityCounts() HostVulnCounts {
	var hostCounts = HostVulnCounts{}

	for _, cve := range assessment.CVEs {
		for _, pkg := range cve.Packages {

			switch strings.ToLower(pkg.Severity) {
			case "critical":
				hostCounts.Critical++
				if pkg.FixedVersion != "" {
					hostCounts.CritFixable++
				}
			case "high":
				hostCounts.High++
				if pkg.FixedVersion != "" {
					hostCounts.HighFixable++
				}
			case "medium":
				hostCounts.Medium++
				if pkg.FixedVersion != "" {
					hostCounts.MedFixable++
				}
			case "low":
				hostCounts.Low++
				if pkg.FixedVersion != "" {
					hostCounts.LowFixable++
				}
			default:
				hostCounts.Negligible++
				if pkg.FixedVersion != "" {
					hostCounts.NegFixable++
				}
			}
		}
	}

	return hostCounts
}

type HostVulnCounts struct {
	Critical     int32
	CritFixable  int32
	High         int32
	HighFixable  int32
	Medium       int32
	MedFixable   int32
	Low          int32
	LowFixable   int32
	Negligible   int32
	NegFixable   int32
	Total        int32
	TotalFixable int32
}

type HostVulnSeverityCounts struct {
	Critical   *hostVulnSeverityCountsDetails `json:"Critical"`
	High       *hostVulnSeverityCountsDetails `json:"High"`
	Medium     *hostVulnSeverityCountsDetails `json:"Medium"`
	Low        *hostVulnSeverityCountsDetails `json:"Low"`
	Negligible *hostVulnSeverityCountsDetails `json:"Negligible"`
}

func (counts *HostVulnSeverityCounts) VulnerabilityCounts() HostVulnCounts {
	var hostCounts = HostVulnCounts{}

	if counts.Critical != nil {
		hostCounts.Critical += counts.Critical.Vulnerabilities
		hostCounts.CritFixable += counts.Critical.Fixable
		hostCounts.Total += counts.Critical.Vulnerabilities
		hostCounts.TotalFixable += counts.Critical.Fixable
	}

	if counts.High != nil {
		hostCounts.High += counts.High.Vulnerabilities
		hostCounts.HighFixable += counts.High.Fixable
		hostCounts.Total += counts.High.Vulnerabilities
		hostCounts.TotalFixable += counts.High.Fixable
	}

	if counts.Medium != nil {
		hostCounts.Medium += counts.Medium.Vulnerabilities
		hostCounts.MedFixable += counts.Medium.Fixable
		hostCounts.Total += counts.Medium.Vulnerabilities
		hostCounts.TotalFixable += counts.Medium.Fixable
	}

	if counts.Low != nil {
		hostCounts.Low += counts.Low.Vulnerabilities
		hostCounts.LowFixable += counts.Low.Fixable
		hostCounts.Total += counts.Low.Vulnerabilities
		hostCounts.TotalFixable += counts.Low.Fixable
	}

	if counts.Negligible != nil {
		hostCounts.Negligible += counts.Negligible.Vulnerabilities
		hostCounts.NegFixable += counts.Negligible.Fixable
		hostCounts.Total += counts.Negligible.Vulnerabilities
		hostCounts.TotalFixable += counts.Negligible.Fixable
	}

	return hostCounts
}

type hostVulnSeverityCountsDetails struct {
	Fixable         int32 `json:"fixable"`
	Vulnerabilities int32 `json:"vulnerabilities"`
}

type HostVulnCveSummary struct {
	Severity             HostVulnSeverityCounts `json:"severity"`
	TotalVulnerabilities int                    `json:"total_vulnerabilities"`
	LastEvaluationTime   Json16DigitTime        `json:"last_evaluation_time"`
}

// TODO @afiune add scan response
//    {
//      "CVE_PROPS": {
//        "cve_batch_id": "d4cca68d-8850-4f77-8ce4-554a434dbbf9",
//        "description": "The OpenSSL ECDSA signature algorithm has been shown to be vulnerable to a timing side channel attack. An attacker could use variations in the signing algorithm to recover the private key. Fixed in OpenSSL 1.1.0j (Affected 1.1.0-1.1.0i). Fixed in OpenSSL 1.1.1a (Affected 1.1.1).",
//        "link": "http://people.ubuntu.com/~ubuntu-security/cve/CVE-2018-0735",
//        "metadata": {
//          "NVD": {
//            "CVSSv2": {
//              "PublishedDateTime": "2018-10-29T13:29Z",
//              "Score": 4.3,
//              "Vectors": "AV:N/AC:M/Au:N/C:P/I:N/A:N"
//            },
//            "CVSSv3": {
//              "ExploitabilityScore": 2.2,
//              "ImpactScore": 3.6,
//              "Score": 5.9,
//              "Vectors": "CVSS:3.0/AV:N/AC:H/PR:N/UI:N/S:U/C:H/I:N/A:N"
//            }
//          }
//        }
//      },
//      "FEATURE_KEY": {
//        "name": "openssl",
//        "namespace": "ubuntu:18.04"
//      },
//      "FIX_INFO": {
//        "compare_result": -1,
//        "eval_status": "GOOD",
//        "fix_available": 0,
//        "fixed_version": "0:1.1.0g-2ubuntu4.3",
//        "fixed_version_comparison_infos": [
//          {
//            "curr_fix_ver": "1.1.0g-2ubuntu4.3",
//            "is_curr_fix_ver_greater_than_other_fix_ver": "0",
//            "other_fix_ver": "1.1.0g-2ubuntu4.3"
//          }
//        ],
//        "fixed_version_comparison_score": 0,
//        "max_prefix_matching_len_score": 6,
//        "version_installed": "0:1.1.1-1ubuntu2.1~18.04.5"
//      },
//      "OS_PKG_INFO": {
//        "namespace": "ubuntu:18.04",
//        "os": "Ubuntu",
//        "os_ver": "18.04",
//        "pkg": "openssl",
//        "pkg_ver": "1.1.1-1ubuntu2.1~18.04.5",
//        "version_format": "dpkg"
//      },
//      "PROPS": {
//        "eval_algo": "1001"
//      },
//      "SEVERITY": "Low",
//      "SUMMARY": {
//        "eval_created_time": "Mon, 17 Aug 2020 06:27:42 -0700",
//        "eval_status": "MATCH_VULN",
//        "num_fixable_vuln": 4,
//        "num_fixable_vuln_by_severity": {
//          "1": 0,
//          "2": 0,
//          "3": 4,
//          "4": 0,
//          "5": 0
//        },
//        "num_total": 64,
//        "num_vuln": 4,
//        "num_vuln_by_severity": {
//          "1": 0,
//          "2": 0,
//          "3": 4,
//          "4": 0,
//          "5": 0
//        }
//      },
//      "VULN_ID": "CVE-2018-0735"
//    },
