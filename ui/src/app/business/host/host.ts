import {BaseModel, BaseRequest} from '../../shared/class/BaseModel';
import {Credential} from '../setting/credential/credential';
import {Cluster} from '../cluster/cluster';

export class Host extends BaseModel {

    id: string;
    name: string;
    ip: string;
    port: string;
    credentialId: string;
    os: string;
    osVersion: string;
    memory: string;
    cpuCore: number;
    gpuNum: number;
    gpuInfo: string;
    status: string;
    volumes: Volume[];
    credential: Credential;
    cluster: Cluster;
}

export class Volume extends BaseModel {
    ID: string;
    size: string;
    name: string;
    hostId: string;
}

export class HostCreateRequest extends BaseRequest {
    name: string;
    ip: string;
    port: string;
    credentialId: string;
}
