package cluster

import (
	"github.com/KubeOperator/KubeOperator/pkg/model/common"
	"github.com/KubeOperator/KubeOperator/pkg/model/host"
	"github.com/KubeOperator/kobe/api"
	uuid "github.com/satori/go.uuid"
	"time"
)

type Node struct {
	common.BaseModel
	Host   host.Host
	HostID string
	Labels []Label
}

func (n Node) LabelValue(name string) string {
	result := ""
	for _, label := range n.Labels {
		if n.Name == name {
			result = label.Value
		}
	}
	return result
}

func (n Node) ToKobeHost() *api.Host {
	return &api.Host{
		Ip:       n.Host.Ip,
		Name:     n.Host.Name,
		Port:     int32(n.Host.Port),
		User:     n.Host.User,
		Password: n.Host.Password,
	}
}

func (n *Node) BeforeCreate() error {
	n.ID = uuid.NewV4().String()
	n.CreatedDate = time.Now()
	n.UpdatedDate = time.Now()
	return nil
}

func (n *Node) BeforeUpdate() error {
	n.UpdatedDate = time.Now()
	return nil
}

func (n Node) TableName() string {
	return "ko_cluster_node"
}
