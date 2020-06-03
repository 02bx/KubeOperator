import {BaseModel} from '../../shared/class/BaseModel';

export class Cluster extends BaseModel {
    name: string;
    spec: Spec;
    nodeSize: string;
    status: string;
}

export class Spec {
    version: string;
}

export class ClusterStatusResponse {
    status: Status;
}

export class Status {
    phase: string;
    conditions: Condition[] = [];
}

export class Condition {
    status: string;
    message: string;
    name: string;
}

export class CreateNodeRequest {
    role: string;
    hostName: string;
}


export class ClusterCreateRequest extends BaseModel {
    name: string;
    version: string;
    provider: string;
    networkType: string;
    runtimeType: string;
    dockerStorageDir: string;
    containerdStorageDir: string;
    clusterCIDR: string;
    serviceCIDR: string;
    nodes: CreateNodeRequest[] = [];
}

export class InitClusterResponse {
    message: string;
}
