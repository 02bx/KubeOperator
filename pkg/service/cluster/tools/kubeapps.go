package tools

import (
	"encoding/json"
	"fmt"
	"github.com/KubeOperator/KubeOperator/pkg/constant"
	"github.com/KubeOperator/KubeOperator/pkg/model"
)

type Kubeapps struct {
	Tool          *model.ClusterTool
	Cluster       *Cluster
	LocalhostName string
}

func NewKubeapps(cluster *Cluster, localhostName string, tool *model.ClusterTool) (*Kubeapps, error) {
	p := &Kubeapps{
		Tool:          tool,
		Cluster:       cluster,
		LocalhostName: localhostName,
	}
	return p, nil
}

func (k Kubeapps) setDefaultValue() {
	values := map[string]interface{}{}
	_ = json.Unmarshal([]byte(k.Tool.Vars), &values)
	values["global.imageRegistry"] = fmt.Sprintf("%s:%d", k.LocalhostName, constant.LocalDockerRepositoryPort)
	values["apprepository.initialRepos[0].name"] = "kubeoperator"
	values["apprepository.initialRepos[0].url"] = fmt.Sprintf("http://%s:%d:8081/repository/kubeapps", k.LocalhostName, constant.LocalHelmRepositoryPort)
	values["useHelm3"] = true
	values["postgresql.enabled"] = true
	str, _ := json.Marshal(&values)
	k.Tool.Vars = string(str)
}

func (k Kubeapps) Uninstall() error {
	return uninstall(k.Tool, constant.DefaultKubeappsIngress, k.Cluster.HelmClient, k.Cluster.KubeClient)
}

func (c Kubeapps) Install() error {
	c.setDefaultValue()
	if err := installChart(c.Cluster.HelmClient, c.Tool, constant.KubeappsChartName); err != nil {
		return err
	}
	if err := createRoute(constant.DefaultKubeappsIngressName, constant.DefaultKubeappsIngress, constant.DefaultKubeappsServiceName, 80, c.Cluster.KubeClient); err != nil {
		return err
	}
	if err := waitForRunning(constant.DefaultDashboardDeploymentName, 1, c.Cluster.KubeClient); err != nil {
		return err
	}
	return nil
}
