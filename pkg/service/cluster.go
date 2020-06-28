package service

import (
	"fmt"
	"github.com/KubeOperator/KubeOperator/pkg/constant"
	"github.com/KubeOperator/KubeOperator/pkg/model"
	"github.com/KubeOperator/KubeOperator/pkg/repository"
	"github.com/KubeOperator/KubeOperator/pkg/service/dto"
	clusterUtil "github.com/KubeOperator/KubeOperator/pkg/util/cluster"
)

type ClusterService interface {
	Get(name string) (dto.Cluster, error)
	GetStatus(name string) (dto.ClusterStatus, error)
	GetSecrets(name string) (dto.ClusterSecret, error)
	GetMonitor(name string) (dto.ClusterMonitor, error)
	GetSpec(name string) (dto.ClusterSpec, error)
	GetApiServerEndpoint(name string) (dto.Endpoint, error)
	GetRouterEndpoint(name string) (dto.Endpoint, error)
	Delete(name string) error
	Create(creation dto.ClusterCreate) error
	List() ([]dto.Cluster, error)
	Page(num, size int) (dto.ClusterPage, error)
	Batch(batch dto.ClusterBatch) error
}

func NewClusterService() ClusterService {
	return &clusterService{
		clusterRepo:                repository.NewClusterRepository(),
		clusterSpecRepo:            repository.NewClusterSpecRepository(),
		clusterNodeRepo:            repository.NewClusterNodeRepository(),
		clusterStatusRepo:          repository.NewClusterStatusRepository(),
		clusterSecretRepo:          repository.NewClusterSecretRepository(),
		clusterMonitorRepo:         repository.NewClusterMonitorRepository(),
		clusterStatusConditionRepo: repository.NewClusterStatusConditionRepository(),
		hostRepo:                   repository.NewHostRepository(),
		clusterInitService:         NewClusterInitService(),
	}
}

type clusterService struct {
	clusterRepo                repository.ClusterRepository
	clusterSpecRepo            repository.ClusterSpecRepository
	clusterNodeRepo            repository.ClusterNodeRepository
	clusterStatusRepo          repository.ClusterStatusRepository
	clusterSecretRepo          repository.ClusterSecretRepository
	clusterMonitorRepo         repository.ClusterMonitorRepository
	clusterStatusConditionRepo repository.ClusterStatusConditionRepository
	hostRepo                   repository.HostRepository
	clusterInitService         ClusterInitService
}

func (c clusterService) Get(name string) (dto.Cluster, error) {
	var clusterDTO dto.Cluster
	mo, err := c.clusterRepo.Get(name)
	if err != nil {
		return clusterDTO, err
	}
	clusterDTO.Cluster = mo
	clusterDTO.NodeSize = len(mo.Nodes)
	clusterDTO.Status = mo.Status.Phase
	return clusterDTO, nil
}

func (c clusterService) List() ([]dto.Cluster, error) {
	var clusterDTOS []dto.Cluster
	mos, err := c.clusterRepo.List()
	if err != nil {
		return clusterDTOS, nil
	}
	for _, mo := range mos {
		clusterDTOS = append(clusterDTOS, dto.Cluster{
			Cluster:  mo,
			NodeSize: len(mo.Nodes),
			Status:   mo.Status.Phase,
		})
	}
	return clusterDTOS, err
}

func (c clusterService) Page(num, size int) (dto.ClusterPage, error) {
	var page dto.ClusterPage
	total, mos, err := c.clusterRepo.Page(num, size)
	if err != nil {
		return page, nil
	}
	for _, mo := range mos {
		page.Items = append(page.Items, dto.Cluster{
			Cluster:  mo,
			NodeSize: len(mo.Nodes),
			Status:   mo.Status.Phase,
		})
	}
	page.Total = total
	return page, err
}

func (c clusterService) GetSecrets(name string) (dto.ClusterSecret, error) {
	var secret dto.ClusterSecret
	cluster, err := c.clusterRepo.Get(name)
	if err != nil {
		return secret, err
	}
	cs, err := c.clusterSecretRepo.Get(cluster.SecretID)
	if err != nil {
		return secret, err
	}
	secret.ClusterSecret = cs

	return secret, nil
}

func (c clusterService) GetStatus(name string) (dto.ClusterStatus, error) {
	var status dto.ClusterStatus
	cluster, err := c.clusterRepo.Get(name)
	if err != nil {
		return status, err
	}
	cs, err := c.clusterStatusRepo.Get(cluster.StatusID)
	if err != nil {
		return status, err
	}
	status.ClusterStatus = cs
	return status, nil
}

func (c clusterService) GetSpec(name string) (dto.ClusterSpec, error) {
	var spec dto.ClusterSpec
	cluster, err := c.clusterRepo.Get(name)
	if err != nil {
		return spec, err
	}
	cs, err := c.clusterSpecRepo.Get(cluster.SpecID)
	if err != nil {
		return spec, err
	}
	spec.ClusterSpec = cs
	return spec, nil
}

func (c clusterService) GetMonitor(name string) (dto.ClusterMonitor, error) {
	var monitor dto.ClusterMonitor
	cluster, err := c.clusterRepo.Get(name)
	if err != nil {
		return monitor, err
	}
	cm, err := c.clusterMonitorRepo.Get(cluster.MonitorID)
	if err != nil {
		return monitor, err
	}
	monitor.ClusterMonitor = cm
	return monitor, nil
}

func (c clusterService) Create(creation dto.ClusterCreate) error {
	cluster := model.Cluster{
		Name:   creation.Name,
		Source: constant.ClusterSourceLocal,
	}
	spec := model.ClusterSpec{
		RuntimeType:          creation.RuntimeType,
		DockerStorageDir:     creation.DockerStorageDIr,
		ContainerdStorageDir: creation.ContainerdStorageDIr,
		NetworkType:          creation.NetworkType,
		ClusterCIDR:          creation.ClusterCIDR,
		ServiceCIDR:          creation.ServiceCIDR,
		Version:              creation.Version,
		AppDomain:            creation.AppDomain,
		KubeApiServerPort:    constant.DefaultApiServerPort,
	}
	status := model.ClusterStatus{Phase: constant.ClusterWaiting}
	secret := model.ClusterSecret{
		KubeadmToken: clusterUtil.GenerateKubeadmToken(),
	}
	monitor := model.ClusterMonitor{Enable: false}
	cluster.Spec = spec
	cluster.Status = status
	cluster.Secret = secret
	cluster.Monitor = monitor
	workerNo := 1
	masterNo := 1
	for _, nc := range creation.Nodes {
		node := model.ClusterNode{
			ClusterID: cluster.ID,
			Role:      nc.Role,
		}
		switch node.Role {
		case constant.NodeRoleNameMaster:
			node.Name = fmt.Sprintf("%s-%d", constant.NodeRoleNameMaster, masterNo)
			masterNo++
		case constant.NodeRoleNameWorker:
			node.Name = fmt.Sprintf("%s-%d", constant.NodeRoleNameWorker, workerNo)
			workerNo++
		}
		host, err := c.hostRepo.Get(nc.HostName)
		if err != nil {
			return err
		}
		node.HostID = host.ID
		node.Host = host
		cluster.Nodes = append(cluster.Nodes, node)
	}
	if len(cluster.Nodes) > 0 {
		cluster.Spec.KubeRouter = cluster.Nodes[0].Host.Ip
	}
	if err := c.clusterRepo.Save(&cluster); err != nil {
		return err
	}
	//if err := c.clusterInitService.Init(cluster.Name); err != nil {
	//	return err
	//}
	return nil
}

func (c clusterService) GetApiServerEndpoint(name string) (dto.Endpoint, error) {
	cluster, err := c.clusterRepo.Get(name)
	var endpoint dto.Endpoint
	if err != nil {
		return endpoint, err
	}
	endpoint.Port = cluster.Spec.KubeApiServerPort
	if cluster.Spec.LbKubeApiserverIp != "" {
		endpoint.Address = cluster.Spec.LbKubeApiserverIp
		return endpoint, nil
	}
	master, err := c.clusterNodeRepo.FistMaster(cluster.ID)
	if err != nil {
		return endpoint, err
	}
	endpoint.Address = master.Host.Ip
	return endpoint, nil
}

func (c clusterService) GetRouterEndpoint(name string) (dto.Endpoint, error) {
	cluster, err := c.clusterRepo.Get(name)
	var endpoint dto.Endpoint
	if err != nil {
		return endpoint, err
	}
	endpoint.Address = cluster.Spec.KubeRouter
	return endpoint, nil
}

func (c clusterService) Delete(name string) error {
	return c.clusterRepo.Delete(name)
}

func (c clusterService) Batch(batch dto.ClusterBatch) error {
	switch batch.Operation {
	case constant.BatchOperationDelete:
		var names []string
		for _, item := range batch.Items {
			names = append(names, item.Name)
		}
		if err := c.clusterRepo.BatchDelete(names...); err != nil {
			return err
		}
	}
	return nil
}
