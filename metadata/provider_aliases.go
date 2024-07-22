// Copyright (c) The OpenTofu Authors
// SPDX-License-Identifier: Apache-2.0

package metadata

import (
	"context"

	provider "github.com/opentofu/libregistry/types/provider"
)

func (r registryDataAPI) ListProviderNamespaceAliases(_ context.Context) (map[string]string, error) {
	// TODO: move this to a JSON file.
	return map[string]string{
		"opentofu": "hashicorp",
	}, nil
}

func (r registryDataAPI) ListProviderAliases(_ context.Context) (map[provider.Addr]provider.Addr, error) {
	// TODO: move this to a JSON file.
	//goland:noinspection GoStructInitializationWithoutFieldNames
	aliases := map[provider.Addr]provider.Addr{
		provider.Addr{"hashicorp", "aci"}:              {"CiscoDevNet", "aci"},
		provider.Addr{"hashicorp", "acme"}:             {"vancluever", "acme"},
		provider.Addr{"hashicorp", "akamai"}:           {"akamai", "akamai"},
		provider.Addr{"hashicorp", "alicloud"}:         {"aliyun", "alicloud"},
		provider.Addr{"hashicorp", "aviatrix"}:         {"AviatrixSystems", "aviatrix"},
		provider.Addr{"hashicorp", "avi"}:              {"vmware", "avi"},
		provider.Addr{"hashicorp", "azuredevops"}:      {"microsoft", "azuredevops"},
		provider.Addr{"hashicorp", "baiducloud"}:       {"baidubce", "baiducloud"},
		provider.Addr{"hashicorp", "bigip"}:            {"F5Networks", "bigip"},
		provider.Addr{"hashicorp", "brightbox"}:        {"brightbox", "brightbox"},
		provider.Addr{"hashicorp", "checkpoint"}:       {"CheckPointSW", "checkpoint"},
		provider.Addr{"hashicorp", "circonus"}:         {"circonus-labs", "circonus"},
		provider.Addr{"hashicorp", "cloudflare"}:       {"cloudflare", "cloudflare"},
		provider.Addr{"hashicorp", "cloudscale"}:       {"cloudscale-ch", "cloudscale"},
		provider.Addr{"hashicorp", "constellix"}:       {"Constellix", "constellix"},
		provider.Addr{"hashicorp", "datadog"}:          {"DataDog", "datadog"},
		provider.Addr{"hashicorp", "digitalocean"}:     {"digitalocean", "digitalocean"},
		provider.Addr{"hashicorp", "dme"}:              {"DNSMadeEasy", "dme"},
		provider.Addr{"hashicorp", "dnsimple"}:         {"dnsimple", "dnsimple"}, //Manually detected from incorrect homepage
		provider.Addr{"hashicorp", "dome9"}:            {"dome9", "dome9"},
		provider.Addr{"hashicorp", "exoscale"}:         {"exoscale", "exoscale"},
		provider.Addr{"hashicorp", "fastly"}:           {"fastly", "fastly"},
		provider.Addr{"hashicorp", "flexibleengine"}:   {"FlexibleEngineCloud", "flexibleengine"},
		provider.Addr{"hashicorp", "fortios"}:          {"fortinetdev", "fortios"},
		provider.Addr{"hashicorp", "github"}:           {"integrations", "github"},
		provider.Addr{"hashicorp", "gitlab"}:           {"gitlabhq", "gitlab"},
		provider.Addr{"hashicorp", "grafana"}:          {"grafana", "grafana"},
		provider.Addr{"hashicorp", "gridscale"}:        {"gridscale", "gridscale"},
		provider.Addr{"hashicorp", "hcloud"}:           {"hetznercloud", "hcloud"},
		provider.Addr{"hashicorp", "heroku"}:           {"heroku", "heroku"},
		provider.Addr{"hashicorp", "huaweicloud"}:      {"huaweicloud", "huaweicloud"},
		provider.Addr{"hashicorp", "huaweicloudstack"}: {"huaweicloud", "huaweicloudstack"},
		provider.Addr{"hashicorp", "icinga2"}:          {"Icinga", "icinga2"},
		provider.Addr{"hashicorp", "launchdarkly"}:     {"launchdarkly", "launchdarkly"},
		provider.Addr{"hashicorp", "linode"}:           {"linode", "linode"},
		provider.Addr{"hashicorp", "logicmonitor"}:     {"logicmonitor", "logicmonitor"}, //Manually detected from incorrect homepage
		provider.Addr{"hashicorp", "mongodbatlas"}:     {"mongodb", "mongodbatlas"},
		provider.Addr{"hashicorp", "ncloud"}:           {"NaverCloudPlatform", "ncloud"},
		provider.Addr{"hashicorp", "newrelic"}:         {"newrelic", "newrelic"},
		provider.Addr{"hashicorp", "ns1"}:              {"ns1-terraform", "ns1"},
		provider.Addr{"hashicorp", "nsxt"}:             {"vmware", "nsxt"},
		provider.Addr{"hashicorp", "nutanix"}:          {"nutanix", "nutanix"},
		provider.Addr{"hashicorp", "oci"}:              {"oracle", "oci"},
		provider.Addr{"hashicorp", "oktaasa"}:          {"oktadeveloper", "oktaasa"},
		provider.Addr{"hashicorp", "okta"}:             {"oktadeveloper", "okta"},
		provider.Addr{"hashicorp", "opennebula"}:       {"OpenNebula", "opennebula"},
		provider.Addr{"hashicorp", "openstack"}:        {"openstack", "openstack"},
		provider.Addr{"hashicorp", "opentelekomcloud"}: {"opentelekomcloud", "opentelekomcloud"},
		provider.Addr{"hashicorp", "opsgenie"}:         {"opsgenie", "opsgenie"},
		provider.Addr{"hashicorp", "ovh"}:              {"ovh", "ovh"},
		provider.Addr{"hashicorp", "packet"}:           {"packethost", "packet"},
		provider.Addr{"hashicorp", "pagerduty"}:        {"PagerDuty", "pagerduty"},
		provider.Addr{"hashicorp", "panos"}:            {"PaloAltoNetworks", "panos"},
		provider.Addr{"hashicorp", "powerdns"}:         {"pan-net", "powerdns"},
		provider.Addr{"hashicorp", "prismacloud"}:      {"PaloAltoNetworks", "prismacloud"},
		provider.Addr{"hashicorp", "profitbricks"}:     {"ionos-cloud", "profitbricks"},
		provider.Addr{"hashicorp", "rancher2"}:         {"rancher", "rancher2"},
		provider.Addr{"hashicorp", "rundeck"}:          {"rundeck", "rundeck"},
		provider.Addr{"hashicorp", "scaleway"}:         {"scaleway", "scaleway"},
		provider.Addr{"hashicorp", "selectel"}:         {"selectel", "selectel"},
		provider.Addr{"hashicorp", "signalfx"}:         {"splunk-terraform", "signalfx"}, // Repo was moved "signalfx", "signalfx",
		provider.Addr{"hashicorp", "skytap"}:           {"skytap", "skytap"},
		provider.Addr{"hashicorp", "spotinst"}:         {"spotinst", "spotinst"},
		provider.Addr{"hashicorp", "stackpath"}:        {"stackpath", "stackpath"},
		provider.Addr{"hashicorp", "statuscake"}:       {"StatusCakeDev", "statuscake"},
		provider.Addr{"hashicorp", "sumologic"}:        {"SumoLogic", "sumologic"},
		provider.Addr{"hashicorp", "tencentcloud"}:     {"tencentcloudstack", "tencentcloud"},
		provider.Addr{"hashicorp", "triton"}:           {"joyent", "triton"},
		provider.Addr{"hashicorp", "turbot"}:           {"turbot", "turbot"},
		provider.Addr{"hashicorp", "ucloud"}:           {"ucloud", "ucloud"},
		provider.Addr{"hashicorp", "vcd"}:              {"vmware", "vcd"},
		provider.Addr{"hashicorp", "venafi"}:           {"Venafi", "venafi"},
		provider.Addr{"hashicorp", "vmc"}:              {"vmware", "vmc"},
		provider.Addr{"hashicorp", "vra7"}:             {"vmware", "vra7"},
		provider.Addr{"hashicorp", "vultr"}:            {"vultr", "vultr"},
		provider.Addr{"hashicorp", "wavefront"}:        {"vmware", "wavefront"},
		provider.Addr{"hashicorp", "yandex"}:           {"yandex-cloud", "yandex"},
	}
	result := make(map[provider.Addr]provider.Addr, len(aliases))
	for from, to := range aliases {
		result[from.Normalize()] = to.Normalize()
	}
	return result, nil
}
