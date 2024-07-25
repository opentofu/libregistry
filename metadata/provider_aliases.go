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
		"opentofu": "hashicorp",
	}, nil
}

func (r registryDataAPI) ListProviderAliases(_ context.Context) (map[provider.Addr]provider.Addr, error) {
	// TODO: move this to a JSON file.
	aliases := map[provider.Addr]provider.Addr{
		provider.Addr{Namespace: "hashicorp", Name: "aci"}:              {Namespace: "CiscoDevNet", Name: "aci"},
		provider.Addr{Namespace: "hashicorp", Name: "acme"}:             {Namespace: "vancluever", Name: "acme"},
		provider.Addr{Namespace: "hashicorp", Name: "akamai"}:           {Namespace: "akamai", Name: "akamai"},
		provider.Addr{Namespace: "hashicorp", Name: "alicloud"}:         {Namespace: "aliyun", Name: "alicloud"},
		provider.Addr{Namespace: "hashicorp", Name: "aviatrix"}:         {Namespace: "AviatrixSystems", Name: "aviatrix"},
		provider.Addr{Namespace: "hashicorp", Name: "avi"}:              {Namespace: "vmware", Name: "avi"},
		provider.Addr{Namespace: "hashicorp", Name: "azuredevops"}:      {Namespace: "microsoft", Name: "azuredevops"},
		provider.Addr{Namespace: "hashicorp", Name: "baiducloud"}:       {Namespace: "baidubce", Name: "baiducloud"},
		provider.Addr{Namespace: "hashicorp", Name: "bigip"}:            {Namespace: "F5Networks", Name: "bigip"},
		provider.Addr{Namespace: "hashicorp", Name: "brightbox"}:        {Namespace: "brightbox", Name: "brightbox"},
		provider.Addr{Namespace: "hashicorp", Name: "checkpoint"}:       {Namespace: "CheckPointSW", Name: "checkpoint"},
		provider.Addr{Namespace: "hashicorp", Name: "circonus"}:         {Namespace: "circonus-labs", Name: "circonus"},
		provider.Addr{Namespace: "hashicorp", Name: "cloudflare"}:       {Namespace: "cloudflare", Name: "cloudflare"},
		provider.Addr{Namespace: "hashicorp", Name: "cloudscale"}:       {Namespace: "cloudscale-ch", Name: "cloudscale"},
		provider.Addr{Namespace: "hashicorp", Name: "constellix"}:       {Namespace: "Constellix", Name: "constellix"},
		provider.Addr{Namespace: "hashicorp", Name: "datadog"}:          {Namespace: "DataDog", Name: "datadog"},
		provider.Addr{Namespace: "hashicorp", Name: "digitalocean"}:     {Namespace: "digitalocean", Name: "digitalocean"},
		provider.Addr{Namespace: "hashicorp", Name: "dme"}:              {Namespace: "DNSMadeEasy", Name: "dme"},
		provider.Addr{Namespace: "hashicorp", Name: "dnsimple"}:         {Namespace: "dnsimple", Name: "dnsimple"}, //Manually detected from incorrect homepage
		provider.Addr{Namespace: "hashicorp", Name: "dome9"}:            {Namespace: "dome9", Name: "dome9"},
		provider.Addr{Namespace: "hashicorp", Name: "exoscale"}:         {Namespace: "exoscale", Name: "exoscale"},
		provider.Addr{Namespace: "hashicorp", Name: "fastly"}:           {Namespace: "fastly", Name: "fastly"},
		provider.Addr{Namespace: "hashicorp", Name: "flexibleengine"}:   {Namespace: "FlexibleEngineCloud", Name: "flexibleengine"},
		provider.Addr{Namespace: "hashicorp", Name: "fortios"}:          {Namespace: "fortinetdev", Name: "fortios"},
		provider.Addr{Namespace: "hashicorp", Name: "github"}:           {Namespace: "integrations", Name: "github"},
		provider.Addr{Namespace: "hashicorp", Name: "gitlab"}:           {Namespace: "gitlabhq", Name: "gitlab"},
		provider.Addr{Namespace: "hashicorp", Name: "grafana"}:          {Namespace: "grafana", Name: "grafana"},
		provider.Addr{Namespace: "hashicorp", Name: "gridscale"}:        {Namespace: "gridscale", Name: "gridscale"},
		provider.Addr{Namespace: "hashicorp", Name: "hcloud"}:           {Namespace: "hetznercloud", Name: "hcloud"},
		provider.Addr{Namespace: "hashicorp", Name: "heroku"}:           {Namespace: "heroku", Name: "heroku"},
		provider.Addr{Namespace: "hashicorp", Name: "huaweicloud"}:      {Namespace: "huaweicloud", Name: "huaweicloud"},
		provider.Addr{Namespace: "hashicorp", Name: "huaweicloudstack"}: {Namespace: "huaweicloud", Name: "huaweicloudstack"},
		provider.Addr{Namespace: "hashicorp", Name: "icinga2"}:          {Namespace: "Icinga", Name: "icinga2"},
		provider.Addr{Namespace: "hashicorp", Name: "launchdarkly"}:     {Namespace: "launchdarkly", Name: "launchdarkly"},
		provider.Addr{Namespace: "hashicorp", Name: "linode"}:           {Namespace: "linode", Name: "linode"},
		provider.Addr{Namespace: "hashicorp", Name: "logicmonitor"}:     {Namespace: "logicmonitor", Name: "logicmonitor"}, //Manually detected from incorrect homepage
		provider.Addr{Namespace: "hashicorp", Name: "mongodbatlas"}:     {Namespace: "mongodb", Name: "mongodbatlas"},
		provider.Addr{Namespace: "hashicorp", Name: "ncloud"}:           {Namespace: "NaverCloudPlatform", Name: "ncloud"},
		provider.Addr{Namespace: "hashicorp", Name: "newrelic"}:         {Namespace: "newrelic", Name: "newrelic"},
		provider.Addr{Namespace: "hashicorp", Name: "ns1"}:              {Namespace: "ns1-terraform", Name: "ns1"},
		provider.Addr{Namespace: "hashicorp", Name: "nsxt"}:             {Namespace: "vmware", Name: "nsxt"},
		provider.Addr{Namespace: "hashicorp", Name: "nutanix"}:          {Namespace: "nutanix", Name: "nutanix"},
		provider.Addr{Namespace: "hashicorp", Name: "oci"}:              {Namespace: "oracle", Name: "oci"},
		provider.Addr{Namespace: "hashicorp", Name: "oktaasa"}:          {Namespace: "oktadeveloper", Name: "oktaasa"},
		provider.Addr{Namespace: "hashicorp", Name: "okta"}:             {Namespace: "oktadeveloper", Name: "okta"},
		provider.Addr{Namespace: "hashicorp", Name: "opennebula"}:       {Namespace: "OpenNebula", Name: "opennebula"},
		provider.Addr{Namespace: "hashicorp", Name: "openstack"}:        {Namespace: "openstack", Name: "openstack"},
		provider.Addr{Namespace: "hashicorp", Name: "opentelekomcloud"}: {Namespace: "opentelekomcloud", Name: "opentelekomcloud"},
		provider.Addr{Namespace: "hashicorp", Name: "opsgenie"}:         {Namespace: "opsgenie", Name: "opsgenie"},
		provider.Addr{Namespace: "hashicorp", Name: "ovh"}:              {Namespace: "ovh", Name: "ovh"},
		provider.Addr{Namespace: "hashicorp", Name: "packet"}:           {Namespace: "packethost", Name: "packet"},
		provider.Addr{Namespace: "hashicorp", Name: "pagerduty"}:        {Namespace: "PagerDuty", Name: "pagerduty"},
		provider.Addr{Namespace: "hashicorp", Name: "panos"}:            {Namespace: "PaloAltoNetworks", Name: "panos"},
		provider.Addr{Namespace: "hashicorp", Name: "powerdns"}:         {Namespace: "pan-net", Name: "powerdns"},
		provider.Addr{Namespace: "hashicorp", Name: "prismacloud"}:      {Namespace: "PaloAltoNetworks", Name: "prismacloud"},
		provider.Addr{Namespace: "hashicorp", Name: "profitbricks"}:     {Namespace: "ionos-cloud", Name: "profitbricks"},
		provider.Addr{Namespace: "hashicorp", Name: "rancher2"}:         {Namespace: "rancher", Name: "rancher2"},
		provider.Addr{Namespace: "hashicorp", Name: "rundeck"}:          {Namespace: "rundeck", Name: "rundeck"},
		provider.Addr{Namespace: "hashicorp", Name: "scaleway"}:         {Namespace: "scaleway", Name: "scaleway"},
		provider.Addr{Namespace: "hashicorp", Name: "selectel"}:         {Namespace: "selectel", Name: "selectel"},
		provider.Addr{Namespace: "hashicorp", Name: "signalfx"}:         {Namespace: "splunk-terraform", Name: "signalfx"}, // Repo was moved "signalfx", "signalfx",
		provider.Addr{Namespace: "hashicorp", Name: "skytap"}:           {Namespace: "skytap", Name: "skytap"},
		provider.Addr{Namespace: "hashicorp", Name: "spotinst"}:         {Namespace: "spotinst", Name: "spotinst"},
		provider.Addr{Namespace: "hashicorp", Name: "stackpath"}:        {Namespace: "stackpath", Name: "stackpath"},
		provider.Addr{Namespace: "hashicorp", Name: "statuscake"}:       {Namespace: "StatusCakeDev", Name: "statuscake"},
		provider.Addr{Namespace: "hashicorp", Name: "sumologic"}:        {Namespace: "SumoLogic", Name: "sumologic"},
		provider.Addr{Namespace: "hashicorp", Name: "tencentcloud"}:     {Namespace: "tencentcloudstack", Name: "tencentcloud"},
		provider.Addr{Namespace: "hashicorp", Name: "triton"}:           {Namespace: "joyent", Name: "triton"},
		provider.Addr{Namespace: "hashicorp", Name: "turbot"}:           {Namespace: "turbot", Name: "turbot"},
		provider.Addr{Namespace: "hashicorp", Name: "ucloud"}:           {Namespace: "ucloud", Name: "ucloud"},
		provider.Addr{Namespace: "hashicorp", Name: "vcd"}:              {Namespace: "vmware", Name: "vcd"},
		provider.Addr{Namespace: "hashicorp", Name: "venafi"}:           {Namespace: "Venafi", Name: "venafi"},
		provider.Addr{Namespace: "hashicorp", Name: "vmc"}:              {Namespace: "vmware", Name: "vmc"},
		provider.Addr{Namespace: "hashicorp", Name: "vra7"}:             {Namespace: "vmware", Name: "vra7"},
		provider.Addr{Namespace: "hashicorp", Name: "vultr"}:            {Namespace: "vultr", Name: "vultr"},
		provider.Addr{Namespace: "hashicorp", Name: "wavefront"}:        {Namespace: "vmware", Name: "wavefront"},
		provider.Addr{Namespace: "hashicorp", Name: "yandex"}:           {Namespace: "yandex-cloud", Name: "yandex"},
	}
	result := make(map[provider.Addr]provider.Addr, len(aliases))
	for from, to := range aliases {
		result[from.Normalize()] = to.Normalize()
	}
	return result, nil
}
