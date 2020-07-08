package client

import (
	"context"
	"github.com/vmware/govmomi"
	"github.com/vmware/govmomi/find"
	"github.com/vmware/govmomi/object"
	"github.com/vmware/govmomi/property"
	"github.com/vmware/govmomi/view"
	"github.com/vmware/govmomi/vim25/mo"
	"github.com/vmware/govmomi/vim25/soap"
	"github.com/vmware/govmomi/vim25/types"
	"net/url"
)

type vSphereClient struct {
	Vars    map[string]interface{}
	Connect Connect
}

type Connect struct {
	Client govmomi.Client
	Ctx    context.Context
}

func NewVSphereClient(vars map[string]interface{}) *vSphereClient {
	return &vSphereClient{
		Vars: vars,
	}
}

func (v *vSphereClient) ListZones() string {
	return ""
}

func (v *vSphereClient) ListDatacenter() ([]string, error) {
	_, err := v.GetConnect()
	if err != nil {
		return nil, err
	}
	client := v.Connect.Client.Client
	var result []string

	var datacenters []*object.Datacenter
	f := find.NewFinder(client, true)
	datacenters, err = f.DatacenterList(v.Connect.Ctx, "*")
	if err != nil {
		return nil, err
	}

	for _, d := range datacenters {
		result = append(result, d.Common.InventoryPath)
	}
	return result, nil
}

func (v *vSphereClient) ListClusters() ([]interface{}, error) {
	_, err := v.GetConnect()
	if err != nil {
		return nil, err
	}
	client := v.Connect.Client.Client
	var result []interface{}

	m := view.NewManager(client)
	view, err := m.CreateContainerView(v.Connect.Ctx, client.ServiceContent.RootFolder, []string{"ClusterComputeResource"}, true)
	if err != nil {
		return result, err
	}

	var clusters []mo.ClusterComputeResource
	err = view.Retrieve(v.Connect.Ctx, []string{"ClusterComputeResource"}, []string{"summary", "name", "resourcePool", "network", "datastore"}, &clusters)
	if err != nil {
		return result, err
	}

	for _, d := range clusters {
		var clusterData map[string]interface{}
		clusterData = make(map[string]interface{})

		clusterData["cluster"] = d.ManagedEntity.Name
		networks, _ := v.GetNetwork(d.ComputeResource.Network)
		clusterData["networks"] = networks
		datastores, _ := v.GetDatastore(d.ComputeResource.Datastore)
		clusterData["datastores"] = datastores
		resourcePools, _ := v.GetResourcePools(*d.ComputeResource.ResourcePool)
		clusterData["resourcePools"] = resourcePools

		result = append(result, clusterData)
	}

	return result, nil
}

func (v *vSphereClient) ListTemplates() ([]interface{}, error) {
	_, err := v.GetConnect()
	if err != nil {
		return nil, err
	}
	client := v.Connect.Client.Client
	var result []interface{}

	m := view.NewManager(client)

	w, err := m.CreateContainerView(v.Connect.Ctx, client.ServiceContent.RootFolder, []string{"VirtualMachine"}, true)
	if err != nil {
		return result, err
	}

	var vms []mo.VirtualMachine
	err = w.Retrieve(v.Connect.Ctx, []string{"VirtualMachine"}, []string{"summary", "name"}, &vms)
	if err != nil {
		return result, err
	}

	for _, vm := range vms {
		var template map[string]string
		template = make(map[string]string)
		if vm.Summary.Config.Template {
			template["imageName"] = vm.Summary.Config.Name
			template["guestId"] = vm.Summary.Config.GuestId
			result = append(result, template)
		}
	}

	return result, nil
}

func (v *vSphereClient) GetNetwork(mos []types.ManagedObjectReference) ([]string, error) {

	pc := property.DefaultCollector(v.Connect.Client.Client)
	rps := []mo.Network{}
	var data []string
	err := pc.Retrieve(v.Connect.Ctx, mos, []string{"summary", "name"}, &rps)
	if err != nil {
		return data, err
	}
	for _, d := range rps {
		data = append(data, d.Name)
	}
	return data, nil
}

func (v *vSphereClient) GetDatastore(mos []types.ManagedObjectReference) ([]string, error) {

	pc := property.DefaultCollector(v.Connect.Client.Client)
	rps := []mo.Datastore{}
	var data []string
	err := pc.Retrieve(v.Connect.Ctx, mos, []string{"summary", "name"}, &rps)
	if err != nil {
		return data, err
	}
	for _, d := range rps {
		data = append(data, d.Name)
	}
	return data, nil
}

func (v *vSphereClient) GetResourcePools(m types.ManagedObjectReference) ([]string, error) {

	pc := property.DefaultCollector(v.Connect.Client.Client)
	rp := mo.ResourcePool{}
	var data []string
	err := pc.RetrieveOne(v.Connect.Ctx, m, []string{"summary", "name", "resourcePool"}, &rp)
	if err != nil {
		return data, err
	}
	data = append(data, rp.Name)

	rps := []mo.ResourcePool{}
	err = pc.Retrieve(v.Connect.Ctx, rp.ResourcePool, []string{"summary", "name"}, &rps)
	if err != nil {
		return data, err
	}
	for _, r := range rps {
		data = append(data, r.Name)
	}

	return data, nil
}

func (v *vSphereClient) GetConnect() (Connect, error) {
	ctx, _ := context.WithCancel(context.Background())
	u, err := soap.ParseURL(v.Vars["host"].(string))
	if err != nil {
		return Connect{}, err
	}
	u.User = url.UserPassword(v.Vars["username"].(string), v.Vars["password"].(string))
	c, err := govmomi.NewClient(ctx, u, true)
	if err != nil {
		return Connect{}, err
	}
	connect := &Connect{
		Client: *c,
		Ctx:    ctx,
	}
	v.Connect = *connect
	return *connect, nil
}
