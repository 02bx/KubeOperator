package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/KubeOperator/KubeOperator/pkg/cloud_provider/client"
	"github.com/KubeOperator/KubeOperator/pkg/constant"
	"github.com/KubeOperator/KubeOperator/pkg/model"
	"github.com/KubeOperator/KubeOperator/pkg/model/common"
	"github.com/KubeOperator/KubeOperator/pkg/repository"
	"github.com/KubeOperator/KubeOperator/pkg/util/ipaddr"
	"github.com/KubeOperator/KubeOperator/pkg/util/kotf"
	"github.com/KubeOperator/KubeOperator/pkg/util/lang"
	"strings"
)

type ClusterIaasService interface {
	Init(name string) error
}

func NewClusterIaasService() ClusterIaasService {
	return &clusterIaasService{
		clusterRepo: repository.NewClusterRepository(),
		nodeRepo:    repository.NewClusterNodeRepository(),
		hostRepo:    repository.NewHostRepository(),
		planRepo:    repository.NewPlanRepository(),
	}
}

type clusterIaasService struct {
	clusterRepo repository.ClusterRepository
	hostRepo    repository.HostRepository
	nodeRepo    repository.ClusterNodeRepository
	planRepo    repository.PlanRepository
}

func (c clusterIaasService) Init(name string) error {
	cluster, err := c.clusterRepo.Get(name)
	if err != nil {
		return err
	}
	if cluster.Spec.Provider == constant.ClusterProviderBareMetal || len(cluster.Nodes) > 0 {
		return nil
	}
	plan, err := c.planRepo.GetById(cluster.PlanID)
	if err != nil {
		return err
	}
	hosts, cloudHosts, err := c.createHosts(cluster, plan)
	if err != nil {
		return err
	}
	err = c.hostRepo.BatchSave(hosts)
	if err != nil {
		return err
	}
	k := kotf.NewTerraform(&kotf.Config{Cluster: name})
	err = doInit(k, plan, cloudHosts)
	if err != nil {
		var hs []model.Host
		for _, host := range hosts {
			hs = append(hs, *host)
		}
		_ = c.hostRepo.Batch(constant.BatchOperationDelete, hs)
		return err
	}
	nodes, err := c.createNodes(cluster, hosts)
	if err := c.nodeRepo.BatchSave(nodes); err != nil {
		return err
	}
	return nil
}

func (c clusterIaasService) createNodes(cluster model.Cluster, hosts []*model.Host) ([]*model.ClusterNode, error) {
	masterNum := 0
	workerNum := 0
	var nodes []*model.ClusterNode
	for _, host := range hosts {
		role := getHostRole(host.Name)
		no := 0
		if role == constant.NodeRoleNameMaster {
			masterNum += 1
			no = masterNum
		} else {
			workerNum += 1
			no = workerNum
		}
		node := model.ClusterNode{
			Name:      fmt.Sprintf("%s-%d", role, no),
			HostID:    host.ID,
			ClusterID: cluster.ID,
			Role:      role,
		}
		nodes = append(nodes, &node)
	}
	if err := c.nodeRepo.BatchSave(nodes); err != nil {
		return nil, err
	}
	return nodes, nil
}

func (c clusterIaasService) createHosts(cluster model.Cluster, plan model.Plan) ([]*model.Host, []map[string]interface{}, error) {
	var hosts []*model.Host
	masterAmount := 1
	if plan.DeployTemplate != constant.SINGLE {
		masterAmount = 3
	}
	for i := 0; i < masterAmount; i++ {
		host := model.Host{
			BaseModel: common.BaseModel{},
			Name:      fmt.Sprintf("%s-master-%d", cluster.Name, i+1),
			Port:      22,
			Status:    constant.ClusterWaiting,
			ClusterID: cluster.ID,
		}
		hosts = append(hosts, &host)
	}
	for i := 0; i < cluster.Spec.WorkerAmount; i++ {
		host := model.Host{
			BaseModel: common.BaseModel{},
			Name:      fmt.Sprintf("%s-worker-%d", cluster.Name, i+1),
			Port:      22,
			Status:    constant.ClusterWaiting,
			ClusterID: cluster.ID,
		}
		hosts = append(hosts, &host)
	}
	var selectedIps []string
	group := allocateZone(plan.Zones, hosts)
	for k, v := range group {
		providerVars := map[string]interface{}{}
		providerVars["provider"] = plan.Region.Provider
		_ = json.Unmarshal([]byte(plan.Region.Vars), &providerVars)
		cloudClient := client.NewCloudClient(providerVars)
		err := allocateIpAddr(cloudClient, *k, v, selectedIps)
		if err != nil {
			return nil, nil, err
		}
	}
	return hosts, parseHosts(group, plan), nil
}

func parseHosts(group map[*model.Zone][]*model.Host, plan model.Plan) []map[string]interface{} {
	switch plan.Region.Provider {
	case constant.VSphere:
		return parseVsphereHosts(group, plan)
	case constant.OpenStack:
		return parseOpenstackHosts(group, plan)
	}
	return []map[string]interface{}{}
}

func parseVsphereHosts(group map[*model.Zone][]*model.Host, plan model.Plan) []map[string]interface{} {
	var results []map[string]interface{}
	planVars := map[string]string{}
	_ = json.Unmarshal([]byte(plan.Vars), &planVars)
	for k, v := range group {
		for _, h := range v {
			var zoneVars map[string]interface{}
			_ = json.Unmarshal([]byte(k.Vars), &zoneVars)
			zoneVars["key"] = formatZoneName(k.Name)
			role := getHostRole(h.Name)
			hMap := map[string]interface{}{}
			hMap["name"] = h.Name
			hMap["shortName"] = h.Name
			hMap["cpu"] = constant.VmConfigList[planVars[fmt.Sprintf("%sModel", role)]].Cpu
			hMap["memory"] = constant.VmConfigList[planVars[fmt.Sprintf("%sModel", role)]].Memory * 1024
			hMap["ip"] = h.Ip
			hMap["zone"] = zoneVars
			results = append(results, hMap)
		}
	}
	return results
}

func parseOpenstackHosts(group map[*model.Zone][]*model.Host, plan model.Plan) []map[string]interface{} {
	var results []map[string]interface{}
	planVars := map[string]string{}
	_ = json.Unmarshal([]byte(plan.Vars), &planVars)
	for k, v := range group {
		for _, h := range v {
			role := getHostRole(h.Name)
			hMap := map[string]interface{}{}
			hMap["name"] = h.Name
			hMap["shortName"] = h.Name
			hMap["ip"] = h.Ip
			hMap["model"] = planVars[fmt.Sprintf("%sModel", role)]
			hMap["zone"] = k.Vars
			results = append(results, hMap)
		}
	}
	return results
}

func getHostRole(name string) string {
	if strings.Contains(name, "-master-") {
		return constant.NodeRoleNameMaster
	}
	return constant.NodeRoleNameWorker
}

func doInit(k *kotf.Kotf, plan model.Plan, hosts []map[string]interface{}) error {
	var zonesVars []map[string]interface{}
	for _, zone := range plan.Zones {
		zoneMap := map[string]interface{}{}
		_ = json.Unmarshal([]byte(zone.Vars), &zoneMap)
		zoneMap["key"] = formatZoneName(zone.Name)
		zonesVars = append(zonesVars, zoneMap)
	}
	hostsStr, _ := json.Marshal(&hosts)
	cloudRegion := map[string]interface{}{
		"datacenter": plan.Region.Datacenter,
		"zones":      zonesVars,
	}
	cloudRegionStr, _ := json.Marshal(&cloudRegion)
	res, err := k.Init(plan.Region.Provider, plan.Region.Vars, string(cloudRegionStr), string(hostsStr))
	if err != nil {
		return err
	}
	if !res.Success {
		return errors.New(res.GetMsg())
	}
	return doApply(k)
}

func doApply(k *kotf.Kotf) error {
	res, err := k.Apply()
	if err != nil {
		return err
	}
	if !res.Success {
		return errors.New(res.GetMsg())
	}
	return nil
}

func allocateZone(zones []model.Zone, hosts []*model.Host) map[*model.Zone][]*model.Host {
	groupMap := map[*model.Zone][]*model.Host{}
	for i, _ := range hosts {
		hash := i % len(zones)
		groupMap[&zones[hash]] = append(groupMap[&zones[hash]], hosts[i])
		hosts[i].CredentialID = zones[hash].CredentialID
	}
	return groupMap
}

func allocateIpAddr(p client.CloudClient, zone model.Zone, hosts []*model.Host, selectedIps []string) error {
	ips := ipaddr.GenerateIps("172.16.10.0", 24)
	zoneVars := map[string]string{}
	_ = json.Unmarshal([]byte(zone.Vars), &zoneVars)
	pool, err := p.GetIpInUsed(zoneVars["network"])
	if err != nil {
		return err
	}
end:
	for _, h := range hosts {
		for i, _ := range ips {
			if !exists(ips[i], pool) && !exists(ips[i], selectedIps) {
				h.Ip = ips[i]
				selectedIps = append(selectedIps, h.Ip)
				continue end
			}
		}
	}
	return nil
}

func exists(ip string, pool []string) bool {
	for _, i := range pool {
		if ip == i {
			return true
		}
	}
	return false
}

func formatZoneName(name string) string {
	if lang.CountChinese(name) > 0 {
		return lang.Pinyin(name)
	}
	return name
}
