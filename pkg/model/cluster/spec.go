package cluster

import (
	"github.com/KubeOperator/KubeOperator/pkg/model/common"
	uuid "github.com/satori/go.uuid"
)

type Spec struct {
	common.BaseModel
	ID                   string
	Version              string
	Provider             string
	NetworkType          string
	RuntimeType          string
	DockerStorageDir     string
	ContainerdStorageDir string
	ClusterCIDR          string `gorm:"column:cluster_cidr"`
	ServiceCIDR          string `gorm:"column:service_cidr"`
}

func (s *Spec) BeforeCreate() (err error) {
	s.ID = uuid.NewV4().String()
	return nil
}

func (s Spec) TableName() string {
	return "ko_cluster_spec"
}
