import json

from rest_framework.views import APIView
from kubeops_api.models.item_resource import ItemResource
from kubeops_api.models.item import Item
from kubeops_api.models.cluster import Cluster
from kubeops_api.models.item_resource_dto import Resource
from kubeops_api.models.backup_storage import BackupStorage
from cloud_provider.models import Plan
from kubeops_api.models.host import Host
from storage.models import NfsStorage, CephStorage
from django.http import HttpResponse
from kubeops_api.utils.json_resource_encoder import JsonResourceEncoder
from rest_framework import viewsets
from kubeops_api.serializers.item import ItemSerializer
from kubeops_api.models.node import Node
from storage.models import ClusterCephStorage
from kubeops_api.models.cluster_backup import ClusterBackup


__all__ = ["ItemResourceView"]


class ItemViewSet(viewsets.ModelViewSet):
    queryset = Item.objects.all()
    serializer_class = ItemSerializer

    lookup_field = 'name'
    lookup_url_kwarg = 'name'


class ItemResourceView(APIView):

    def get(self, request, *args, **kwargs):
        item_name = kwargs['item_name']
        item = Item.objects.get(name=item_name)
        resource_ids = ItemResource.objects.filter(item_id=item.id).values_list('resource_id', flat=True)
        resources = []
        clusters = Cluster.objects.filter(id__in=resource_ids)
        for c in clusters:
            resource = Resource(resource_id=c.id, resource_type=ItemResource.RESOURCE_TYPE_CLUSTER, data=c,
                                checked=True)
            resources.append(resource.__dict__)
        hosts = Host.objects.filter(id__in=resource_ids)
        for h in hosts:
            resource = Resource(resource_id=h.id, resource_type=ItemResource.RESOURCE_TYPE_HOST, data=h, checked=True)
            resources.append(resource.__dict__)
        backup_storage = BackupStorage.objects.filter(id__in=resource_ids)
        for b in backup_storage:
            resource = Resource(resource_id=b.id, resource_type=ItemResource.RESOURCE_TYPE_BACKUP_STORAGE, data=b,
                                checked=True)
            resources.append(resource.__dict__)
        plan = Plan.objects.filter(id__in=resource_ids)
        for p in plan:
            resource = Resource(resource_id=p.id, resource_type=ItemResource.RESOURCE_TYPE_PLAN, data=p, checked=True)
            resources.append(resource.__dict__)
        nfs = NfsStorage.objects.filter(id__in=resource_ids)
        for n in nfs:
            resource = Resource(resource_id=n.id, resource_type=ItemResource.RESOURCE_TYPE_STORAGE, data=n,
                                checked=True)
            resources.append(resource.__dict__)
        ceph = CephStorage.objects.filter(id__in=resource_ids)
        for c in ceph:
            resource = Resource(resource_id=c.id, resource_type=ItemResource.RESOURCE_TYPE_STORAGE, data=n,
                                checked=True)
            resources.append(resource)

        response = HttpResponse(content_type='application/json')
        response.write(json.dumps(resources, cls=JsonResourceEncoder))
        return response

    def delete(self, request, *args, **kwargs):
        resource_id = kwargs['item_name']
        ItemResource.objects.get(resource_id=resource_id).delete()
        response = HttpResponse(content_type='application/json')
        response.write(json.dumps({'msg': '取消成功'}))
        return response

    def post(self, request, *args, **kwargs):
        item_resources = request.data
        objs = []
        for item_resource in item_resources:
            obj = ItemResource(resource_type=item_resource['resource_type'], resource_id=item_resource['resource_id'],
                               item_id=item_resource['item_id'])
            objs.append(obj)
            if item_resource['resource_type'] == ItemResource.RESOURCE_TYPE_CLUSTER:
                cluster = Cluster.objects.get(id=item_resource['resource_id'])
                nodes = Node.objects.filter(project_id=cluster.id)
                for node in nodes:
                    if node.host_id is not None:
                        node_obj = ItemResource(resource_type=ItemResource.RESOURCE_TYPE_HOST, resource_id=node.host_id,
                                     item_id=item_resource['item_id'])
                        objs.append(node_obj)
                if cluster.persistent_storage == 'nfs':
                    nfs_name = cluster.configs['nfs']
                    nfs = NfsStorage.objects.get(name=nfs_name)
                    nfs_obj = ItemResource(resource_type=ItemResource.RESOURCE_TYPE_STORAGE, resource_id=nfs.id,
                                           item_id=item_resource['item_id'])
                    objs.append(nfs_obj)
                if cluster.persistent_storage == 'ceph':
                    ceph = ClusterCephStorage.objects.get(cluster_id=cluster.id)
                    ceph_obj = ItemResource(resource_type=ItemResource.RESOURCE_TYPE_STORAGE, resource_id=ceph.id,
                                           item_id=item_resource['item_id'])
                    objs.append(ceph_obj)
                if cluster.deploy_type == Cluster.CLUSTER_DEPLOY_TYPE_AUTOMATIC:
                    plan_obj = ItemResource(resource_type=ItemResource.RESOURCE_TYPE_PLAN, resource_id=cluster.plan_id,
                                           item_id=item_resource['item_id'])
                    objs.append(plan_obj)
                cluster_backup = ClusterBackup.objects.filter(project_id=cluster.id)
                if len(cluster_backup) > 0:
                    backup_obj = ItemResource(resource_type=ItemResource.RESOURCE_TYPE_BACKUP_STORAGE, resource_id=cluster_backup[0].backup_storage_id,
                                           item_id=item_resource['item_id'])
                    objs.append(backup_obj)

        result = ItemResource.objects.bulk_create(objs)
        response = HttpResponse(content_type='application/json')
        response.write(json.dumps({'msg': '授权成功'}))
        return response


class ResourceView(APIView):

    def get(self, request, *args, **kwargs):
        item_name = kwargs['item_name']
        resource_type = kwargs['resource_type']
        item = Item.objects.get(name=item_name)
        data = []
        resource_ids = ItemResource.objects.filter(item_id=item.id).values_list('resource_id', flat=True)
        if resource_type == ItemResource.RESOURCE_TYPE_CLUSTER:
            result = Cluster.objects.exclude(id__in=resource_ids)
            for re in result:
                item_resource_dto = Resource(resource_id=re.id, resource_type=resource_type, data=re, checked=False)
                data.append(item_resource_dto.__dict__)
        response = HttpResponse(content_type='application/json')
        response.write(json.dumps(data, cls=JsonResourceEncoder))
        return response
