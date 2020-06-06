import {Component, OnInit} from '@angular/core';
import {KubernetesService} from '../../../kubernetes.service';
import {ActivatedRoute} from '@angular/router';
import {Cluster} from '../../../cluster';
import {V1Namespace} from '@kubernetes/client-node';

@Component({
    selector: 'app-namespace-list',
    templateUrl: './namespace-list.component.html',
    styleUrls: ['./namespace-list.component.css']
})
export class NamespaceListComponent implements OnInit {

    loading = true;
    selected = [];
    items: V1Namespace[] = [];
    nextToken = '';
    previousToken = '';
    continueToken = '';
    page = 1;
    currentCluster: Cluster;

    constructor(private service: KubernetesService, private route: ActivatedRoute) {
    }


    ngOnInit(): void {
        this.route.parent.data.subscribe(data => {
            this.currentCluster = data.cluster.item;
            this.list();
        });
    }

    list() {
        this.loading = true;
        this.service.listNamespaces(this.currentCluster.name, this.continueToken).subscribe(data => {
            this.loading = false;
            this.items = data.items;
            this.nextToken = data.metadata[this.service.continueTokenKey] ? data.metadata[this.service.continueTokenKey] : '';
        });
    }
}
