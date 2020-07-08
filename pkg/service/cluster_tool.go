package service

import (
	"encoding/json"
	"github.com/KubeOperator/KubeOperator/pkg/constant"
	"github.com/KubeOperator/KubeOperator/pkg/dto"
	"github.com/KubeOperator/KubeOperator/pkg/model"
	"github.com/KubeOperator/KubeOperator/pkg/repository"
	"github.com/KubeOperator/KubeOperator/pkg/service/cluster/tools"
)

type ClusterToolService interface {
	List(clusterName string) ([]dto.ClusterTool, error)
	Enable(clusterName string, tool dto.ClusterTool) (dto.ClusterTool, error)
}

func NewClusterToolService() *clusterToolService {
	return &clusterToolService{
		toolRepo:       repository.NewClusterToolRepository(),
		clusterService: NewClusterService(),
	}
}

type clusterToolService struct {
	toolRepo       repository.ClusterToolRepository
	clusterService ClusterService
}

func (c clusterToolService) List(clusterName string) ([]dto.ClusterTool, error) {
	var items []dto.ClusterTool
	ms, err := c.toolRepo.List(clusterName)
	if err != nil {
		return items, err
	}
	for _, m := range ms {
		d := dto.ClusterTool{ClusterTool: m}
		d.Vars = map[string]interface{}{}
		_ = json.Unmarshal([]byte(m.Vars), &d.Vars)
		items = append(items, d)
	}
	return items, nil
}

func (c clusterToolService) Enable(clusterName string, tool dto.ClusterTool) (dto.ClusterTool, error) {
	cluster, err := c.clusterService.Get(clusterName)
	if err != nil {
		return tool, err
	}
	tool.ClusterID = cluster.ID
	mo := tool.ClusterTool
	buf, _ := json.Marshal(&tool.Vars)
	mo.Vars = string(buf)
	mo.Status = constant.ClusterInitializing
	err = c.toolRepo.Save(&mo)
	if err != nil {
		return tool, err
	}
	tool.ClusterTool = mo

	endpoint, err := c.clusterService.GetApiServerEndpoint(clusterName)
	if err != nil {
		return tool, err
	}
	secret, err := c.clusterService.GetSecrets(clusterName)
	if err != nil {
		return tool, err
	}

	ct, err := tools.NewClusterTool(&tool.ClusterTool, cluster.Cluster, endpoint, secret.ClusterSecret)
	if err != nil {
		return tool, err
	}
	go c.do(ct, &tool.ClusterTool)
	return tool, nil
}

func (c clusterToolService) do(p tools.Interface, tool *model.ClusterTool) {
	err := p.Install()
	if err != nil {
		tool.Status = constant.ClusterFailed
		tool.Message = err.Error()
	} else {
		tool.Status = constant.ClusterRunning
	}
	_ = c.toolRepo.Save(tool)
}
