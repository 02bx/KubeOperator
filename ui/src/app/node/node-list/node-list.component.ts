import {Component, EventEmitter, Input, OnInit, Output, ViewChild} from '@angular/core';
import {NodeService} from '../node.service';
import {Node} from '../node';
import {Cluster} from '../../cluster/cluster';
import {AlertLevels} from '../../base/header/components/common-alert/alert';
import {CommonAlertService} from '../../base/header/common-alert.service';

@Component({
  selector: 'app-node-list',
  templateUrl: './node-list.component.html',
  styleUrls: ['./node-list.component.css']
})
export class NodeListComponent implements OnInit {

  loading = true;
  nodes: Node[] = [];
  @Input() currentCluster: Cluster;
  timeResult;
  openView = false;
  loadingTime = false;

  constructor(private nodeService: NodeService, private alertService: CommonAlertService) {
  }

  ngOnInit() {
    this.listNodes();
  }

  listNodes() {
    this.nodeService.listNodes(this.currentCluster.name).subscribe(data => {
      this.nodes = data.filter((node) => {
        return node.name !== 'localhost' && node.name !== '127.0.0.1' && node.name !== '::1';
      });
      this.loading = false;
    }, error => {
      this.loading = false;
    });
  }

  refresh() {
    this.listNodes();
  }

  toGrafana() {
    const url = 'http://grafana.apps.' + this.currentCluster.name + '.' + this.currentCluster.cluster_doamin_suffix + '/explore';
    window.open(url, '_blank');
  }

  syncTime() {
    this.loadingTime = true;
    this.openView = true;
    this.nodeService.syncHostTime(this.currentCluster.name).subscribe(data => {
      this.timeResult = data;
      this.loadingTime = false;
    }, error1 => {
      this.loadingTime = false;
    });
  }

  checkNodes() {
    this.nodeService.checkNodes(this.currentCluster.name).subscribe(data => {
      this.alertService.showAlert('同步成功', AlertLevels.SUCCESS);
      this.refresh();
    }, error1 => {
      this.alertService.showAlert('同步失败', AlertLevels.ERROR);
    });
  }
}
