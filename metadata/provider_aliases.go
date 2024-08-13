// Copyright (c) The OpenTofu Authors
// SPDX-License-Identifier: MPL-2.0

package metadata

import (
	"context"

	provider "github.com/opentofu/libregistry/types/provider"
)

func (r registryDataAPI) ListProviderNamespaceAliases(_ context.Context) (map[string]string, error) {
	// TODO: move this to a JSON file.
	return map[string]string{
		"hashicorp": "opentofu",
	}, nil
}

func (r registryDataAPI) ListProviderAliases(_ context.Context) (map[provider.Addr]provider.Addr, error) {
	// TODO: move this to a JSON file.
	aliases := map[provider.Addr]provider.Addr{
		provider.Addr{Namespace: "", Name: "aci"}:              {Namespace: "CiscoDevNet", Name: "aci"},
		provider.Addr{Namespace: "", Name: "acme"}:             {Namespace: "vancluever", Name: "acme"},
		provider.Addr{Namespace: "", Name: "akamai"}:           {Namespace: "akamai", Name: "akamai"},
		provider.Addr{Namespace: "", Name: "alicloud"}:         {Namespace: "aliyun", Name: "alicloud"},
		provider.Addr{Namespace: "", Name: "aviatrix"}:         {Namespace: "AviatrixSystems", Name: "aviatrix"},
		provider.Addr{Namespace: "", Name: "avi"}:              {Namespace: "vmware", Name: "avi"},
		provider.Addr{Namespace: "", Name: "azuredevops"}:      {Namespace: "microsoft", Name: "azuredevops"},
		provider.Addr{Namespace: "", Name: "baiducloud"}:       {Namespace: "baidubce", Name: "baiducloud"},
		provider.Addr{Namespace: "", Name: "bigip"}:            {Namespace: "F5Networks", Name: "bigip"},
		provider.Addr{Namespace: "", Name: "brightbox"}:        {Namespace: "brightbox", Name: "brightbox"},
		provider.Addr{Namespace: "", Name: "checkpoint"}:       {Namespace: "CheckPointSW", Name: "checkpoint"},
		provider.Addr{Namespace: "", Name: "circonus"}:         {Namespace: "circonus-labs", Name: "circonus"},
		provider.Addr{Namespace: "", Name: "cloudflare"}:       {Namespace: "cloudflare", Name: "cloudflare"},
		provider.Addr{Namespace: "", Name: "cloudscale"}:       {Namespace: "cloudscale-ch", Name: "cloudscale"},
		provider.Addr{Namespace: "", Name: "constellix"}:       {Namespace: "Constellix", Name: "constellix"},
		provider.Addr{Namespace: "", Name: "datadog"}:          {Namespace: "DataDog", Name: "datadog"},
		provider.Addr{Namespace: "", Name: "digitalocean"}:     {Namespace: "digitalocean", Name: "digitalocean"},
		provider.Addr{Namespace: "", Name: "dme"}:              {Namespace: "DNSMadeEasy", Name: "dme"},
		provider.Addr{Namespace: "", Name: "dnsimple"}:         {Namespace: "dnsimple", Name: "dnsimple"}, //Manually detected from incorrect homepage
		provider.Addr{Namespace: "", Name: "dome9"}:            {Namespace: "dome9", Name: "dome9"},
		provider.Addr{Namespace: "", Name: "exoscale"}:         {Namespace: "exoscale", Name: "exoscale"},
		provider.Addr{Namespace: "", Name: "fastly"}:           {Namespace: "fastly", Name: "fastly"},
		provider.Addr{Namespace: "", Name: "flexibleengine"}:   {Namespace: "FlexibleEngineCloud", Name: "flexibleengine"},
		provider.Addr{Namespace: "", Name: "fortios"}:          {Namespace: "fortinetdev", Name: "fortios"},
		provider.Addr{Namespace: "", Name: "github"}:           {Namespace: "integrations", Name: "github"},
		provider.Addr{Namespace: "", Name: "gitlab"}:           {Namespace: "gitlabhq", Name: "gitlab"},
		provider.Addr{Namespace: "", Name: "grafana"}:          {Namespace: "grafana", Name: "grafana"},
		provider.Addr{Namespace: "", Name: "gridscale"}:        {Namespace: "gridscale", Name: "gridscale"},
		provider.Addr{Namespace: "", Name: "hcloud"}:           {Namespace: "hetznercloud", Name: "hcloud"},
		provider.Addr{Namespace: "", Name: "heroku"}:           {Namespace: "heroku", Name: "heroku"},
		provider.Addr{Namespace: "", Name: "huaweicloud"}:      {Namespace: "huaweicloud", Name: "huaweicloud"},
		provider.Addr{Namespace: "", Name: "huaweicloudstack"}: {Namespace: "huaweicloud", Name: "huaweicloudstack"},
		provider.Addr{Namespace: "", Name: "icinga2"}:          {Namespace: "Icinga", Name: "icinga2"},
		provider.Addr{Namespace: "", Name: "launchdarkly"}:     {Namespace: "launchdarkly", Name: "launchdarkly"},
		provider.Addr{Namespace: "", Name: "linode"}:           {Namespace: "linode", Name: "linode"},
		provider.Addr{Namespace: "", Name: "logicmonitor"}:     {Namespace: "logicmonitor", Name: "logicmonitor"}, //Manually detected from incorrect homepage
		provider.Addr{Namespace: "", Name: "mongodbatlas"}:     {Namespace: "mongodb", Name: "mongodbatlas"},
		provider.Addr{Namespace: "", Name: "ncloud"}:           {Namespace: "NaverCloudPlatform", Name: "ncloud"},
		provider.Addr{Namespace: "", Name: "newrelic"}:         {Namespace: "newrelic", Name: "newrelic"},
		provider.Addr{Namespace: "", Name: "ns1"}:              {Namespace: "ns1-terraform", Name: "ns1"},
		provider.Addr{Namespace: "", Name: "nsxt"}:             {Namespace: "vmware", Name: "nsxt"},
		provider.Addr{Namespace: "", Name: "nutanix"}:          {Namespace: "nutanix", Name: "nutanix"},
		provider.Addr{Namespace: "", Name: "oci"}:              {Namespace: "oracle", Name: "oci"},
		provider.Addr{Namespace: "", Name: "oktaasa"}:          {Namespace: "oktadeveloper", Name: "oktaasa"},
		provider.Addr{Namespace: "", Name: "okta"}:             {Namespace: "oktadeveloper", Name: "okta"},
		provider.Addr{Namespace: "", Name: "opennebula"}:       {Namespace: "OpenNebula", Name: "opennebula"},
		provider.Addr{Namespace: "", Name: "openstack"}:        {Namespace: "openstack", Name: "openstack"},
		provider.Addr{Namespace: "", Name: "opentelekomcloud"}: {Namespace: "opentelekomcloud", Name: "opentelekomcloud"},
		provider.Addr{Namespace: "", Name: "opsgenie"}:         {Namespace: "opsgenie", Name: "opsgenie"},
		provider.Addr{Namespace: "", Name: "ovh"}:              {Namespace: "ovh", Name: "ovh"},
		provider.Addr{Namespace: "", Name: "packet"}:           {Namespace: "packethost", Name: "packet"},
		provider.Addr{Namespace: "", Name: "pagerduty"}:        {Namespace: "PagerDuty", Name: "pagerduty"},
		provider.Addr{Namespace: "", Name: "panos"}:            {Namespace: "PaloAltoNetworks", Name: "panos"},
		provider.Addr{Namespace: "", Name: "powerdns"}:         {Namespace: "pan-net", Name: "powerdns"},
		provider.Addr{Namespace: "", Name: "prismacloud"}:      {Namespace: "PaloAltoNetworks", Name: "prismacloud"},
		provider.Addr{Namespace: "", Name: "profitbricks"}:     {Namespace: "ionos-cloud", Name: "profitbricks"},
		provider.Addr{Namespace: "", Name: "rancher2"}:         {Namespace: "rancher", Name: "rancher2"},
		provider.Addr{Namespace: "", Name: "rundeck"}:          {Namespace: "rundeck", Name: "rundeck"},
		provider.Addr{Namespace: "", Name: "scaleway"}:         {Namespace: "scaleway", Name: "scaleway"},
		provider.Addr{Namespace: "", Name: "selectel"}:         {Namespace: "selectel", Name: "selectel"},
		provider.Addr{Namespace: "", Name: "signalfx"}:         {Namespace: "splunk-terraform", Name: "signalfx"}, // Repo was moved "signalfx", "signalfx",
		provider.Addr{Namespace: "", Name: "skytap"}:           {Namespace: "skytap", Name: "skytap"},
		provider.Addr{Namespace: "", Name: "spotinst"}:         {Namespace: "spotinst", Name: "spotinst"},
		provider.Addr{Namespace: "", Name: "stackpath"}:        {Namespace: "stackpath", Name: "stackpath"},
		provider.Addr{Namespace: "", Name: "statuscake"}:       {Namespace: "StatusCakeDev", Name: "statuscake"},
		provider.Addr{Namespace: "", Name: "sumologic"}:        {Namespace: "SumoLogic", Name: "sumologic"},
		provider.Addr{Namespace: "", Name: "tencentcloud"}:     {Namespace: "tencentcloudstack", Name: "tencentcloud"},
		provider.Addr{Namespace: "", Name: "triton"}:           {Namespace: "joyent", Name: "triton"},
		provider.Addr{Namespace: "", Name: "turbot"}:           {Namespace: "turbot", Name: "turbot"},
		provider.Addr{Namespace: "", Name: "ucloud"}:           {Namespace: "ucloud", Name: "ucloud"},
		provider.Addr{Namespace: "", Name: "vcd"}:              {Namespace: "vmware", Name: "vcd"},
		provider.Addr{Namespace: "", Name: "venafi"}:           {Namespace: "Venafi", Name: "venafi"},
		provider.Addr{Namespace: "", Name: "vmc"}:              {Namespace: "vmware", Name: "vmc"},
		provider.Addr{Namespace: "", Name: "vra7"}:             {Namespace: "vmware", Name: "vra7"},
		provider.Addr{Namespace: "", Name: "vultr"}:            {Namespace: "vultr", Name: "vultr"},
		provider.Addr{Namespace: "", Name: "wavefront"}:        {Namespace: "vmware", Name: "wavefront"},
		provider.Addr{Namespace: "", Name: "yandex"}:           {Namespace: "yandex-cloud", Name: "yandex"},
	}
	result := make(map[provider.Addr]provider.Addr, len(aliases)*2)
	for from, to := range aliases {
		result[provider.Addr{
			Namespace: "hashicorp",
			Name:      from.Name,
		}] = to.Normalize()
		result[provider.Addr{
			Namespace: "opentofu",
			Name:      from.Name,
		}] = to.Normalize()
	}

	return result, nil
}
