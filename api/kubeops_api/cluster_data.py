import json
from uuid import UUID


class ClusterData():

    def __init__(self, cluster, token, pods, nodes, namespaces, deployments):
        self.id = str(cluster.id)
        self.name = cluster.name
        self.pods = pods
        self.nodes = nodes
        self.token = token
        self.namespaces = namespaces
        self.deployments = deployments


class Pod():

    def __init__(self, name, cluster_name, restart_count, status, namespace, host_ip, pod_ip, host_name, containers):
        self.name = name
        self.cluster_name = cluster_name
        self.restart_count = restart_count
        self.status = status
        self.namespace = namespace
        self.host_ip = host_ip
        self.pod_ip = pod_ip
        self.host_name = host_name
        self.containers = containers


class NameSpace():

    def __init__(self, name, status):
        self.name = name
        self.status = status


class Node():

    def __init__(self, name, status):
        self.name = name
        self.status = status


class Container():

    def __init__(self, name, ready, restart_count):
        self.name = name
        self.ready = ready
        self.restart_count = restart_count


class Deployment():

    def __init__(self, name, ready_replicas, replicas, namespace):
        self.name = name
        self.ready_replicas = ready_replicas
        self.replicas = replicas
        self.namespace = namespace
