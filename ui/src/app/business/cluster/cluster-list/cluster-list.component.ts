import {Component, EventEmitter, OnDestroy, OnInit, Output} from '@angular/core';
import {ClusterService} from '../cluster.service';
import {BaseModelComponent} from '../../../shared/class/BaseModelComponent';
import {Cluster} from '../cluster';
import {CommonAlertService} from '../../../layout/common-alert/common-alert.service';
import {AlertLevels} from '../../../layout/common-alert/alert';
import {ActivatedRoute, Router} from '@angular/router';
import {Project} from '../../project/project';

@Component({
    selector: 'app-cluster-list',
    templateUrl: './cluster-list.component.html',
    styleUrls: ['./cluster-list.component.css']
})
export class ClusterListComponent extends BaseModelComponent<Cluster> implements OnInit, OnDestroy {

    constructor(private clusterService: ClusterService,
                private commonAlert: CommonAlertService,
                private router: Router,
                private route: ActivatedRoute) {
        super(clusterService);
    }

    @Output() statusDetailEvent = new EventEmitter<string>();
    @Output() importEvent = new EventEmitter();
    timer;
    currentProject: Project;
    loading = false;


    ngOnInit(): void {
        this.route.parent.data.subscribe(data => {
            this.currentProject = data.project;
            this.polling();
            this.pageBy();
        });
    }

    ngOnDestroy(): void {
        clearInterval(this.timer);
    }

    onDetail(item: Cluster) {
        if (item.status !== 'Running') {
            this.commonAlert.showAlert('cluster is not ready', AlertLevels.ERROR);
        } else {
        this.router.navigate(['projects/' + this.currentProject.name + '/clusters', item.name]).then();
        }
    }

    onImport() {
        this.importEvent.emit();
    }

    onNodeDetail(item: Cluster) {
        if (item.status !== 'Running') {
            this.commonAlert.showAlert('cluster is not ready', AlertLevels.ERROR);
        } else {
            this.router.navigate(['clusters', item.name, 'nodes']).then();
        }
    }


    onStatusDetail(name: string) {
        this.statusDetailEvent.emit(name);
    }

    polling() {
        this.timer = setInterval(() => {
            let flag = false;
            const needPolling = ['Waiting', 'Initializing', 'Terminating', 'Creating'];
            for (const item of this.items) {
                if (needPolling.indexOf(item.status) !== -1) {
                    flag = true;
                    break;
                }
            }
            if (flag) {
                this.clusterService.pageBy(this.page, this.size, this.currentProject.name).subscribe(data => {
                    data.items.forEach(n => {
                        this.items.forEach(item => {
                            if (item.name === n.name) {
                                if (item.status !== n.status) {
                                    item.status = n.status;
                                }
                            }
                        });
                    });
                });
            }
        }, 1000);
    }

    pageBy() {
        this.loading = true;
        this.clusterService.pageBy(this.page, this.size, this.currentProject.name).subscribe(data => {
            this.items = data.items;
            this.total = data.total;
            this.loading = false;
        });
    }
}
