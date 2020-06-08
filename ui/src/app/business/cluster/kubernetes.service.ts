import {Injectable} from '@angular/core';
import {HttpClient} from '@angular/common/http';
import {V1NamespaceList} from '@kubernetes/client-node/dist/gen/model/v1NamespaceList';
import {Observable} from 'rxjs';
import {
    V1beta1CronJobList, V1DaemonSetList,
    V1DeploymentList, V1JobList,
    V1NodeList,
    V1PersistentVolumeClaimList,
    V1PersistentVolumeList, V1StatefulSet, V1StatefulSetList
} from '@kubernetes/client-node';

@Injectable({
    providedIn: 'root'
})

export class KubernetesService {

    proxyUrl = '/api/v1/proxy/{cluster_name}/{resource_url}';
    limit = 10;
    continueTokenKey = 'continue';

    constructor(private client: HttpClient) {
    }

    namespaceUrl = '/api/v1/namespaces';
    persistentVolumesUrl = '/api/v1/persistentvolumes';
    persistentVolumeClaimsUrl = '/api/v1/persistentvolumeclaims';
    namespacePersistentVolumeClaimsUrl = '/api/v1/namespaces/{namespace}/deployments';
    deploymentUrl = 'apis/apps/v1/deployments';
    namespaceDeploymentUrl = 'apis/apps/v1/namespaces/{namespace}/deployments';
    daemonSetUrl = '/apis/apps/v1/daemonsets/';
    statefulSetUrl = '/apis/apps/v1/statefulsets/';
    namespaceStatefulSet = '/apis/apps/v1/namespaces/{namespace}/statefulsets/';
    namespaceDaemonSetUrl = '/apis/apps/v1/namespaces/{namespace}/daemonsets/';
    cornJobUrl = 'apis/batch/v1beta1/cronjobs';
    namespaceCornJobUrl = 'apis/batch/v1beta1/namespaces/{namespace}/cronjobs';
    jobUrl = '/apis/batch/v1/jobs';
    namespaceJobUrl = '/apis/batch/v1/namespaces/{namespace}/jobs';
    nodesUrl = '/api/v1/nodes';

    listNodes(clusterName: string, continueToken?: string): Observable<V1NodeList> {
        let url = this.proxyUrl.replace('{cluster_name}', clusterName).replace('{resource_url}', this.nodesUrl);
        url += '?limit=' + this.limit;
        if (continueToken) {
            url += '&continue=' + continueToken;
        }
        return this.client.get<V1NodeList>(url);
    }

    listNamespaces(clusterName: string, continueToken?: string): Observable<V1NamespaceList> {
        let url = this.proxyUrl.replace('{cluster_name}', clusterName).replace('{resource_url}', this.namespaceUrl);
        url += '?limit=' + this.limit;
        if (continueToken) {
            url += '&continue=' + continueToken;
        }
        return this.client.get<V1NamespaceList>(url);
    }


    listPersistentVolumes(clusterName: string, continueToken?: string): Observable<V1PersistentVolumeList> {
        let url = this.proxyUrl.replace('{cluster_name}', clusterName).replace('{resource_url}', this.persistentVolumesUrl);
        url += '?limit=' + this.limit;
        if (continueToken) {
            url += '&continue=' + continueToken;
        }
        return this.client.get<V1PersistentVolumeList>(url);
    }

    listPersistentVolumeClaims(clusterName: string, continueToken?: string, namespace?: string): Observable<V1PersistentVolumeClaimList> {
        let url = this.proxyUrl.replace('{cluster_name}', clusterName);
        url += '?limit=' + this.limit;
        if (continueToken) {
            url += '&continue=' + continueToken;
        }
        if (namespace) {
            url = url.replace('{resource_url}', this.namespacePersistentVolumeClaimsUrl).replace('{namespace}', namespace);
        } else {
            url = url.replace('{resource_url}', this.persistentVolumeClaimsUrl);
        }
        return this.client.get<V1PersistentVolumeClaimList>(url);
    }

    listDeployment(clusterName: string, continueToken?: string, namespace?: string): Observable<V1DeploymentList> {
        let url = this.proxyUrl.replace('{cluster_name}', clusterName);
        url += '?limit=' + this.limit;
        if (continueToken) {
            url += '&continue=' + continueToken;
        }
        if (namespace) {
            url = url.replace('{resource_url}', this.namespaceDeploymentUrl).replace('{namespace}', namespace);
        } else {
            url = url.replace('{resource_url}', this.deploymentUrl);
        }

        return this.client.get<V1DeploymentList>(url);
    }

    listDaemonSet(clusterName: string, continueToken?: string, namespace?: string): Observable<V1DaemonSetList> {
        let url = this.proxyUrl.replace('{cluster_name}', clusterName);
        url += '?limit=' + this.limit;
        if (continueToken) {
            url += '&continue=' + continueToken;
        }
        if (namespace) {
            url = url.replace('{resource_url}', this.namespaceDaemonSetUrl).replace('{namespace}', namespace);
        } else {
            url = url.replace('{resource_url}', this.daemonSetUrl);
        }

        return this.client.get<V1DaemonSetList>(url);
    }

    listStatefulSet(clusterName: string, continueToken?: string, namespace?: string): Observable<V1StatefulSetList> {
        let url = this.proxyUrl.replace('{cluster_name}', clusterName);
        url += '?limit=' + this.limit;
        if (continueToken) {
            url += '&continue=' + continueToken;
        }
        if (namespace) {
            url = url.replace('{resource_url}', this.namespaceStatefulSet).replace('{namespace}', namespace);
        } else {
            url = url.replace('{resource_url}', this.statefulSetUrl);
        }

        return this.client.get<V1StatefulSetList>(url);
    }

    listJob(clusterName: string, continueToken?: string, namespace?: string): Observable<V1JobList> {
        let url = this.proxyUrl.replace('{cluster_name}', clusterName);
        url += '?limit=' + this.limit;
        if (continueToken) {
            url += '&continue=' + continueToken;
        }
        if (namespace) {
            url = url.replace('{resource_url}', this.namespaceJobUrl).replace('{namespace}', namespace);
        } else {
            url = url.replace('{resource_url}', this.jobUrl);
        }
        return this.client.get<V1JobList>(url);
    }


    listCornJob(clusterName: string, continueToken?: string, namespace?: string): Observable<V1beta1CronJobList> {
        let url = this.proxyUrl.replace('{cluster_name}', clusterName);
        url += '?limit=' + this.limit;
        if (continueToken) {
            url += '&continue=' + continueToken;
        }
        if (namespace) {
            url = url.replace('{resource_url}', this.namespaceCornJobUrl).replace('{namespace}', namespace);
        } else {
            url = url.replace('{resource_url}', this.cornJobUrl);
        }
        return this.client.get<V1beta1CronJobList>(url);
    }
}
