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
		provider.Addr{Namespace: "opentofu", Name: "aci"}:              {Namespace: "CiscoDevNet", Name: "aci"},
		provider.Addr{Namespace: "opentofu", Name: "acme"}:             {Namespace: "vancluever", Name: "acme"},
		provider.Addr{Namespace: "opentofu", Name: "akamai"}:           {Namespace: "akamai", Name: "akamai"},
		provider.Addr{Namespace: "opentofu", Name: "alicloud"}:         {Namespace: "aliyun", Name: "alicloud"},
		provider.Addr{Namespace: "opentofu", Name: "aviatrix"}:         {Namespace: "AviatrixSystems", Name: "aviatrix"},
		provider.Addr{Namespace: "opentofu", Name: "avi"}:              {Namespace: "vmware", Name: "avi"},
		provider.Addr{Namespace: "opentofu", Name: "azuredevops"}:      {Namespace: "microsoft", Name: "azuredevops"},
		provider.Addr{Namespace: "opentofu", Name: "baiducloud"}:       {Namespace: "baidubce", Name: "baiducloud"},
		provider.Addr{Namespace: "opentofu", Name: "bigip"}:            {Namespace: "F5Networks", Name: "bigip"},
		provider.Addr{Namespace: "opentofu", Name: "brightbox"}:        {Namespace: "brightbox", Name: "brightbox"},
		provider.Addr{Namespace: "opentofu", Name: "checkpoint"}:       {Namespace: "CheckPointSW", Name: "checkpoint"},
		provider.Addr{Namespace: "opentofu", Name: "circonus"}:         {Namespace: "circonus-labs", Name: "circonus"},
		provider.Addr{Namespace: "opentofu", Name: "cloudflare"}:       {Namespace: "cloudflare", Name: "cloudflare"},
		provider.Addr{Namespace: "opentofu", Name: "cloudscale"}:       {Namespace: "cloudscale-ch", Name: "cloudscale"},
		provider.Addr{Namespace: "opentofu", Name: "constellix"}:       {Namespace: "Constellix", Name: "constellix"},
		provider.Addr{Namespace: "opentofu", Name: "datadog"}:          {Namespace: "DataDog", Name: "datadog"},
		provider.Addr{Namespace: "opentofu", Name: "digitalocean"}:     {Namespace: "digitalocean", Name: "digitalocean"},
		provider.Addr{Namespace: "opentofu", Name: "dme"}:              {Namespace: "DNSMadeEasy", Name: "dme"},
		provider.Addr{Namespace: "opentofu", Name: "dnsimple"}:         {Namespace: "dnsimple", Name: "dnsimple"}, //Manually detected from incorrect homepage
		provider.Addr{Namespace: "opentofu", Name: "dome9"}:            {Namespace: "dome9", Name: "dome9"},
		provider.Addr{Namespace: "opentofu", Name: "exoscale"}:         {Namespace: "exoscale", Name: "exoscale"},
		provider.Addr{Namespace: "opentofu", Name: "fastly"}:           {Namespace: "fastly", Name: "fastly"},
		provider.Addr{Namespace: "opentofu", Name: "flexibleengine"}:   {Namespace: "FlexibleEngineCloud", Name: "flexibleengine"},
		provider.Addr{Namespace: "opentofu", Name: "fortios"}:          {Namespace: "fortinetdev", Name: "fortios"},
		provider.Addr{Namespace: "opentofu", Name: "github"}:           {Namespace: "integrations", Name: "github"},
		provider.Addr{Namespace: "opentofu", Name: "gitlab"}:           {Namespace: "gitlabhq", Name: "gitlab"},
		provider.Addr{Namespace: "opentofu", Name: "grafana"}:          {Namespace: "grafana", Name: "grafana"},
		provider.Addr{Namespace: "opentofu", Name: "gridscale"}:        {Namespace: "gridscale", Name: "gridscale"},
		provider.Addr{Namespace: "opentofu", Name: "hcloud"}:           {Namespace: "hetznercloud", Name: "hcloud"},
		provider.Addr{Namespace: "opentofu", Name: "heroku"}:           {Namespace: "heroku", Name: "heroku"},
		provider.Addr{Namespace: "opentofu", Name: "huaweicloud"}:      {Namespace: "huaweicloud", Name: "huaweicloud"},
		provider.Addr{Namespace: "opentofu", Name: "huaweicloudstack"}: {Namespace: "huaweicloud", Name: "huaweicloudstack"},
		provider.Addr{Namespace: "opentofu", Name: "icinga2"}:          {Namespace: "Icinga", Name: "icinga2"},
		provider.Addr{Namespace: "opentofu", Name: "launchdarkly"}:     {Namespace: "launchdarkly", Name: "launchdarkly"},
		provider.Addr{Namespace: "opentofu", Name: "linode"}:           {Namespace: "linode", Name: "linode"},
		provider.Addr{Namespace: "opentofu", Name: "logicmonitor"}:     {Namespace: "logicmonitor", Name: "logicmonitor"}, //Manually detected from incorrect homepage
		provider.Addr{Namespace: "opentofu", Name: "mongodbatlas"}:     {Namespace: "mongodb", Name: "mongodbatlas"},
		provider.Addr{Namespace: "opentofu", Name: "ncloud"}:           {Namespace: "NaverCloudPlatform", Name: "ncloud"},
		provider.Addr{Namespace: "opentofu", Name: "newrelic"}:         {Namespace: "newrelic", Name: "newrelic"},
		provider.Addr{Namespace: "opentofu", Name: "ns1"}:              {Namespace: "ns1-terraform", Name: "ns1"},
		provider.Addr{Namespace: "opentofu", Name: "nsxt"}:             {Namespace: "vmware", Name: "nsxt"},
		provider.Addr{Namespace: "opentofu", Name: "nutanix"}:          {Namespace: "nutanix", Name: "nutanix"},
		provider.Addr{Namespace: "opentofu", Name: "oci"}:              {Namespace: "oracle", Name: "oci"},
		provider.Addr{Namespace: "opentofu", Name: "oktaasa"}:          {Namespace: "oktadeveloper", Name: "oktaasa"},
		provider.Addr{Namespace: "opentofu", Name: "okta"}:             {Namespace: "oktadeveloper", Name: "okta"},
		provider.Addr{Namespace: "opentofu", Name: "opennebula"}:       {Namespace: "OpenNebula", Name: "opennebula"},
		provider.Addr{Namespace: "opentofu", Name: "openstack"}:        {Namespace: "openstack", Name: "openstack"},
		provider.Addr{Namespace: "opentofu", Name: "opentelekomcloud"}: {Namespace: "opentelekomcloud", Name: "opentelekomcloud"},
		provider.Addr{Namespace: "opentofu", Name: "opsgenie"}:         {Namespace: "opsgenie", Name: "opsgenie"},
		provider.Addr{Namespace: "opentofu", Name: "ovh"}:              {Namespace: "ovh", Name: "ovh"},
		provider.Addr{Namespace: "opentofu", Name: "packet"}:           {Namespace: "packethost", Name: "packet"},
		provider.Addr{Namespace: "opentofu", Name: "pagerduty"}:        {Namespace: "PagerDuty", Name: "pagerduty"},
		provider.Addr{Namespace: "opentofu", Name: "panos"}:            {Namespace: "PaloAltoNetworks", Name: "panos"},
		provider.Addr{Namespace: "opentofu", Name: "powerdns"}:         {Namespace: "pan-net", Name: "powerdns"},
		provider.Addr{Namespace: "opentofu", Name: "prismacloud"}:      {Namespace: "PaloAltoNetworks", Name: "prismacloud"},
		provider.Addr{Namespace: "opentofu", Name: "profitbricks"}:     {Namespace: "ionos-cloud", Name: "profitbricks"},
		provider.Addr{Namespace: "opentofu", Name: "rancher2"}:         {Namespace: "rancher", Name: "rancher2"},
		provider.Addr{Namespace: "opentofu", Name: "rundeck"}:          {Namespace: "rundeck", Name: "rundeck"},
		provider.Addr{Namespace: "opentofu", Name: "scaleway"}:         {Namespace: "scaleway", Name: "scaleway"},
		provider.Addr{Namespace: "opentofu", Name: "selectel"}:         {Namespace: "selectel", Name: "selectel"},
		provider.Addr{Namespace: "opentofu", Name: "signalfx"}:         {Namespace: "splunk-terraform", Name: "signalfx"}, // Repo was moved "signalfx", "signalfx",
		provider.Addr{Namespace: "opentofu", Name: "skytap"}:           {Namespace: "skytap", Name: "skytap"},
		provider.Addr{Namespace: "opentofu", Name: "spotinst"}:         {Namespace: "spotinst", Name: "spotinst"},
		provider.Addr{Namespace: "opentofu", Name: "stackpath"}:        {Namespace: "stackpath", Name: "stackpath"},
		provider.Addr{Namespace: "opentofu", Name: "statuscake"}:       {Namespace: "StatusCakeDev", Name: "statuscake"},
		provider.Addr{Namespace: "opentofu", Name: "sumologic"}:        {Namespace: "SumoLogic", Name: "sumologic"},
		provider.Addr{Namespace: "opentofu", Name: "tencentcloud"}:     {Namespace: "tencentcloudstack", Name: "tencentcloud"},
		provider.Addr{Namespace: "opentofu", Name: "triton"}:           {Namespace: "joyent", Name: "triton"},
		provider.Addr{Namespace: "opentofu", Name: "turbot"}:           {Namespace: "turbot", Name: "turbot"},
		provider.Addr{Namespace: "opentofu", Name: "ucloud"}:           {Namespace: "ucloud", Name: "ucloud"},
		provider.Addr{Namespace: "opentofu", Name: "vcd"}:              {Namespace: "vmware", Name: "vcd"},
		provider.Addr{Namespace: "opentofu", Name: "venafi"}:           {Namespace: "Venafi", Name: "venafi"},
		provider.Addr{Namespace: "opentofu", Name: "vmc"}:              {Namespace: "vmware", Name: "vmc"},
		provider.Addr{Namespace: "opentofu", Name: "vra7"}:             {Namespace: "vmware", Name: "vra7"},
		provider.Addr{Namespace: "opentofu", Name: "vultr"}:            {Namespace: "vultr", Name: "vultr"},
		provider.Addr{Namespace: "opentofu", Name: "wavefront"}:        {Namespace: "vmware", Name: "wavefront"},
		provider.Addr{Namespace: "opentofu", Name: "yandex"}:           {Namespace: "yandex-cloud", Name: "yandex"},
	}
	result := make(map[provider.Addr]provider.Addr, len(aliases))
	for from, to := range aliases {
		result[from.Normalize()] = to.Normalize()
	}

	return result, nil
}
