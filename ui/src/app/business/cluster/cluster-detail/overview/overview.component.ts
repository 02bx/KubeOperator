import {Component, OnInit} from '@angular/core';
import {ClusterService} from '../../cluster.service';
import {ActivatedRoute} from '@angular/router';
import {Cluster} from '../../cluster';
import {KubernetesService} from '../../kubernetes.service';

@Component({
    selector: 'app-overview',
    templateUrl: './overview.component.html',
    styleUrls: ['./overview.component.css']
})
export class OverviewComponent implements OnInit {

    constructor(private service: ClusterService, private route: ActivatedRoute, private kubernetesService: KubernetesService) {
    }

    currentCluster: Cluster;

    ngOnInit(): void {
        this.route.parent.data.subscribe(data => {
            this.currentCluster = data.cluster.item;
            this.listNamespaces();
        });
    }

    listNamespaces() {
        this.kubernetesService.listNamespaces(this.currentCluster.name).subscribe(data => {
            console.log(data);
        });
    }


}
