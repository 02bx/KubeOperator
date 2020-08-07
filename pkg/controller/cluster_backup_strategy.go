package controller

import (
	"github.com/KubeOperator/KubeOperator/pkg/dto"
	"github.com/KubeOperator/KubeOperator/pkg/service"
	"github.com/kataras/iris/v12/context"
)

type ClusterBackupStrategyController struct {
	Ctx                          context.Context
	CLusterBackupStrategyService service.CLusterBackupStrategyService
}

func NewClusterBackupStrategyController() *ClusterBackupStrategyController {
	return &ClusterBackupStrategyController{
		CLusterBackupStrategyService: service.NewCLusterBackupStrategyService(),
	}
}

func (c ClusterBackupStrategyController) GetStrategyBy(clusterName string) (*dto.ClusterBackupStrategy, error) {
	cb, err := c.CLusterBackupStrategyService.Get(clusterName)
	if err != nil {
		return nil, err
	}
	return cb, nil
}

func (c ClusterBackupStrategyController) PostStrategy() (*dto.ClusterBackupStrategy, error) {
	var req dto.ClusterBackupStrategyRequest
	err := c.Ctx.ReadJSON(&req)
	if err != nil {
		return nil, err
	}
	cb, err := c.CLusterBackupStrategyService.Save(req)
	if err != nil {
		return nil, err
	}
	return cb, nil
}
