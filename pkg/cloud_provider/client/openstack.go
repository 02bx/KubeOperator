package client

import (
	"encoding/json"
	"errors"
	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack"
	"github.com/gophercloud/gophercloud/openstack/blockstorage/v3/volumetypes"
	"github.com/gophercloud/gophercloud/openstack/compute/v2/extensions/availabilityzones"
	"github.com/gophercloud/gophercloud/openstack/compute/v2/extensions/floatingips"
	"github.com/gophercloud/gophercloud/openstack/compute/v2/extensions/secgroups"
	"github.com/gophercloud/gophercloud/openstack/compute/v2/flavors"
	"github.com/gophercloud/gophercloud/openstack/imageservice/v2/images"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/networks"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/subnets"
	"io/ioutil"
	"net/http"
	"strconv"
)

var (
	GetRegionError = "GET_REGION_ERROR"
)

type openStackClient struct {
	Vars map[string]interface{}
}

func NewOpenStackClient(vars map[string]interface{}) *openStackClient {
	return &openStackClient{
		Vars: vars,
	}
}

func (v *openStackClient) ListZones() string {
	return ""
}

func (v *openStackClient) ListDatacenter() ([]string, error) {
	var result []string

	provider, err := v.GetAuth()
	if err != nil {
		return result, err
	}

	client := &http.Client{}
	req, _ := http.NewRequest("GET", v.Vars["identity"].(string)+"/regions", nil)
	req.Header.Add("X-Auth-Token", provider.TokenID)
	resp, _ := client.Do(req)
	body, _ := ioutil.ReadAll(resp.Body)
	m := make(map[string]interface{})
	json.Unmarshal([]byte(body), &m)
	key, exist := m["regions"]
	if exist {
		regions := key.([]interface{})
		for _, r := range regions {
			region := r.(map[string]interface{})
			result = append(result, region["id"].(string))
		}
	} else {
		return result, errors.New(GetRegionError)
	}

	return result, nil
}

func (v *openStackClient) ListClusters() ([]interface{}, error) {
	var result []interface{}

	provider, err := v.GetAuth()
	if err != nil {
		return result, err
	}

	client, err := openstack.NewComputeV2(provider, gophercloud.EndpointOpts{
		Region: v.Vars["datacenter"].(string),
	})
	if err != nil {
		return result, err
	}

	pager, err := availabilityzones.List(client).AllPages()
	if err != nil {
		return result, err
	}
	zones, err := availabilityzones.ExtractAvailabilityZones(pager)
	if err != nil {
		return result, err
	}

	allPages, err := floatingips.List(client).AllPages()
	if err != nil {
		return result, err
	}

	allFloatingIPs, err := floatingips.ExtractFloatingIPs(allPages)
	if err != nil {
		return result, err
	}

	sgPages, err := secgroups.List(client).AllPages()
	if err != nil {
		return result, err
	}
	allSecurityGroups, err := secgroups.ExtractSecurityGroups(sgPages)
	if err != nil {
		return result, err
	}

	iPages, err := images.List(client, images.ListOpts{}).AllPages()
	if err != nil {
		return result, err
	}

	allImages, err := images.ExtractImages(iPages)
	if err != nil {
		panic(err)
	}

	networkClient, err := openstack.NewNetworkV2(provider, gophercloud.EndpointOpts{
		Region: v.Vars["datacenter"].(string),
	})

	networkPager, err := networks.List(networkClient, networks.ListOpts{}).AllPages()
	if err != nil {
		return result, err
	}
	allnetworks, err := networks.ExtractNetworks(networkPager)
	if err != nil {
		return result, err
	}

	blockStorageClient, err := openstack.NewBlockStorageV3(provider, gophercloud.EndpointOpts{
		Region: v.Vars["datacenter"].(string),
	})

	vPages, err := volumetypes.List(blockStorageClient, volumetypes.ListOpts{}).AllPages()
	if err != nil {
		return result, err
	}
	allVPages, err := volumetypes.ExtractVolumeTypes(vPages)
	if err != nil {
		return result, err
	}

	var ipTypes []string
	ipTypes = append(ipTypes, "private")
	ipTypes = append(ipTypes, "floating")

	for _, z := range zones {
		clusterData := make(map[string]interface{})
		clusterData["cluster"] = z.ZoneName

		var networkList []interface{}
		for _, n := range allnetworks {
			networkData := make(map[string]interface{})
			networkData["name"] = n.Name
			networkData["id"] = n.ID
			subnetPages, err := subnets.List(networkClient, subnets.ListOpts{
				NetworkID: n.ID,
			}).AllPages()
			if err != nil {
				continue
			}
			allSubnets, err := subnets.ExtractSubnets(subnetPages)
			if err != nil {
				continue
			}
			var subnetList []interface{}
			for _, s := range allSubnets {
				subnetData := make(map[string]interface{})
				subnetData["id"] = s.ID
				subnetData["name"] = s.Name
				subnetList = append(subnetList, subnetData)
			}
			networkData["subnetList"] = subnetList
			networkList = append(networkList, networkData)
		}
		clusterData["networkList"] = networkList

		var floatingNetworkList []interface{}
		for _, n := range allFloatingIPs {
			floatingNetworkData := make(map[string]interface{})
			floatingNetworkData["id"] = n.ID
			floatingNetworkList = append(floatingNetworkList, floatingNetworkData)
		}
		clusterData["floatingNetworkList"] = floatingNetworkList

		var securityGroups []string
		for _, s := range allSecurityGroups {
			securityGroups = append(securityGroups, s.Name)
		}
		clusterData["securityGroups"] = securityGroups

		var volumeTypes []interface{}
		for _, d := range allVPages {
			volumeData := make(map[string]interface{})
			volumeData["name"] = d.Name
			volumeData["id"] = d.ID
			volumeTypes = append(volumeTypes, volumeData)
		}
		clusterData["storages"] = volumeTypes

		var imageList []interface{}
		for _, i := range allImages {
			imageData := make(map[string]interface{})
			imageData["name"] = i.Name
			imageData["id"] = i.ID
			imageList = append(imageList, imageData)
		}
		clusterData["imageList"] = imageList
		clusterData["ipTypes"] = ipTypes

		result = append(result, clusterData)
	}

	return result, nil
}
func (v *openStackClient) ListTemplates() ([]interface{}, error) {
	return []interface{}{}, nil
}

func (v *openStackClient) GetAuth() (*gophercloud.ProviderClient, error) {

	scope := gophercloud.AuthScope{
		ProjectID: v.Vars["projectId"].(string),
	}

	opts := gophercloud.AuthOptions{
		IdentityEndpoint: v.Vars["identity"].(string),
		Username:         v.Vars["username"].(string),
		Password:         v.Vars["password"].(string),
		DomainName:       v.Vars["domainName"].(string),
		Scope:            &scope,
	}

	provider, err := openstack.AuthenticatedClient(opts)
	if err != nil {
		return nil, err
	}
	return provider, nil
}

func (v *openStackClient) ListFlavors() ([]interface{}, error) {

	var result []interface{}

	provider, err := v.GetAuth()
	if err != nil {
		return result, err
	}

	client, err := openstack.NewComputeV2(provider, gophercloud.EndpointOpts{
		Region: v.Vars["datacenter"].(string),
	})
	if err != nil {
		return result, err
	}
	pager, err := flavors.ListDetail(client, flavors.ListOpts{}).AllPages()
	if err != nil {
		return result, err
	}
	allPages, err := flavors.ExtractFlavors(pager)
	if err != nil {
		return result, err
	}

	for _, f := range allPages {

		if f.RAM > 1024 {
			vmConfig := make(map[string]interface{})
			vmConfig["name"] = f.Name

			config := make(map[string]interface{})
			config["id"], _ = strconv.Atoi(f.ID)
			config["disk"] = f.Disk
			config["cpu"] = f.VCPUs
			config["memory"] = f.RAM / 1024

			vmConfig["config"] = config
			result = append(result, vmConfig)
		}
	}

	return result, nil
}

func (v *openStackClient) GetIpInUsed(network string) ([]string, error) {
	return nil, nil
}
